package analyzer

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"regexp"
	"strings"
	"time"
)

// Parse streams LogEntry values from r, parsed according to format ("json" or "plaintext").
// Returns a channel of entries and a single-value error channel (nil on clean EOF).
// Closes both channels when done.
func Parse(ctx context.Context, r io.Reader, format string) (<-chan LogEntry, <-chan error) {
	entries := make(chan LogEntry)
	errc := make(chan error, 1)

	go func() {
		defer close(entries)
		defer close(errc)

		scanner := bufio.NewScanner(r)
		// Allow up to 1 MiB per line.
		scanner.Buffer(make([]byte, 64*1024), 1024*1024)

		lineNum := 0
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				errc <- ctx.Err()
				return
			default:
			}

			lineNum++
			raw := scanner.Text()
			if strings.TrimSpace(raw) == "" {
				continue
			}

			var entry LogEntry
			var ok bool
			if format == "json" {
				entry, ok = parseJSONLine(raw, lineNum)
			} else {
				entry, ok = parsePlaintextLine(raw, lineNum)
			}
			if !ok {
				continue
			}

			select {
			case <-ctx.Done():
				errc <- ctx.Err()
				return
			case entries <- entry:
			}
		}

		if err := scanner.Err(); err != nil {
			errc <- err
			return
		}
		errc <- nil
	}()

	return entries, errc
}

// sourceFields lists JSON keys tried in order when looking for a source value.
var sourceFields = []string{"source", "service", "module", "logger", "component"}

func parseJSONLine(raw string, lineNum int) (LogEntry, bool) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return LogEntry{}, false
	}

	entry := LogEntry{Raw: raw, LineNumber: lineNum}

	if v, ok := m["timestamp"]; ok {
		var s string
		if json.Unmarshal(v, &s) == nil {
			if t, err := time.Parse(time.RFC3339, s); err == nil {
				entry.Timestamp = t
			}
		}
	}
	if v, ok := m["level"]; ok {
		json.Unmarshal(v, &entry.Level) //nolint:errcheck
		entry.Level = strings.ToLower(entry.Level)
	}
	if v, ok := m["message"]; ok {
		json.Unmarshal(v, &entry.Message) //nolint:errcheck
	} else if v, ok := m["msg"]; ok {
		json.Unmarshal(v, &entry.Message) //nolint:errcheck
	}
	for _, f := range sourceFields {
		if v, ok := m[f]; ok {
			json.Unmarshal(v, &entry.Source) //nolint:errcheck
			break
		}
	}

	return entry, true
}

// plaintextRe matches optional timestamp, level, optional [source], and message.
// Examples:
//
//	2024-01-02T03:04:05Z ERROR [svc] msg
//	ERROR msg
//	2024-01-02 03:04:05 INFO msg
var plaintextRe = regexp.MustCompile(
	`^(?:(\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}(?:Z|[+-]\d{2}:\d{2})?)\s+)?` +
		`(ERROR|WARN|WARNING|INFO|DEBUG)\s+` +
		`(?:\[([^\]]+)\]\s+)?` +
		`(.*)$`,
)

var plaintextTimestampLayouts = []string{
	time.RFC3339,
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05",
}

func parsePlaintextLine(raw string, lineNum int) (LogEntry, bool) {
	m := plaintextRe.FindStringSubmatch(raw)
	if m == nil {
		// No recognisable level — emit with just raw + line number.
		return LogEntry{Raw: raw, LineNumber: lineNum, Message: raw}, true
	}

	entry := LogEntry{
		Raw:        raw,
		LineNumber: lineNum,
		Level:      normalizeLevel(m[2]),
		Source:     m[3],
		Message:    strings.TrimSpace(m[4]),
	}

	if m[1] != "" {
		for _, layout := range plaintextTimestampLayouts {
			if t, err := time.Parse(layout, m[1]); err == nil {
				entry.Timestamp = t
				break
			}
		}
	}

	return entry, true
}

func normalizeLevel(s string) string {
	l := strings.ToLower(s)
	if l == "warning" {
		return "warn"
	}
	return l
}
