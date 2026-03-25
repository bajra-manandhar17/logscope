package analyzer

import (
	"regexp"
	"sort"
	"time"
)

// Masking regexes applied in order — UUID and IP before NUM to avoid partial matches.
var (
	reUUID = regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
	reIP   = regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)
	reHex  = regexp.MustCompile(`\b(?:0x[0-9a-fA-F]+|[0-9a-fA-F]{6,})\b`)
	reNum  = regexp.MustCompile(`\b\d+(?:\.\d+)?`)
)

// maskTokens replaces variable token types with placeholders.
func maskTokens(s string) string {
	s = reUUID.ReplaceAllString(s, "{UUID}")
	s = reIP.ReplaceAllString(s, "{IP}")
	s = reHex.ReplaceAllString(s, "{HEX}")
	s = reNum.ReplaceAllString(s, "{NUM}")
	return s
}

type templateEntry struct {
	sample    string
	count     int
	firstSeen time.Time
	lastSeen  time.Time
}

// PatternAggregator groups log messages by masked template and counts occurrences.
type PatternAggregator struct {
	templates map[string]*templateEntry
}

func NewPatternAggregator() *PatternAggregator {
	return &PatternAggregator{templates: make(map[string]*templateEntry)}
}

// Add masks the message and increments the matching template's count,
// tracking first-seen and last-seen timestamps.
func (a *PatternAggregator) Add(entry LogEntry) {
	tmpl := maskTokens(entry.Message)
	if e, ok := a.templates[tmpl]; ok {
		e.count++
		if !entry.Timestamp.IsZero() {
			if e.firstSeen.IsZero() || entry.Timestamp.Before(e.firstSeen) {
				e.firstSeen = entry.Timestamp
			}
			if entry.Timestamp.After(e.lastSeen) {
				e.lastSeen = entry.Timestamp
			}
		}
		return
	}
	if len(a.templates) >= MaxPatterns {
		return // cap reached — drop new unique templates
	}
	te := &templateEntry{sample: entry.Message, count: 1}
	if !entry.Timestamp.IsZero() {
		te.firstSeen = entry.Timestamp
		te.lastSeen = entry.Timestamp
	}
	a.templates[tmpl] = te
}

// Result returns up to topN patterns sorted by count descending.
func (a *PatternAggregator) Result(topN int) []Pattern {
	if len(a.templates) == 0 {
		return nil
	}

	patterns := make([]Pattern, 0, len(a.templates))
	for tmpl, e := range a.templates {
		patterns = append(patterns, Pattern{
			Template:   tmpl,
			Count:      e.count,
			SampleLine: e.sample,
			FirstSeen:  e.firstSeen,
			LastSeen:   e.lastSeen,
		})
	}
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Count > patterns[j].Count
	})

	if topN < len(patterns) {
		patterns = patterns[:topN]
	}
	return patterns
}
