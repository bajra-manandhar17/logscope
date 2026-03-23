package analyzer

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

const peekLines = 10

// DetectFormat peeks at the first peekLines non-empty lines of r to determine
// whether the stream is "json" or "plaintext". It returns the detected format
// and a new reader that replays all bytes (including those already peeked).
func DetectFormat(r io.Reader) (format string, replay io.Reader, err error) {
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)

	scanner := bufio.NewScanner(tee)
	var sampled, jsonCount int
	for scanner.Scan() && sampled < peekLines {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		sampled++
		if json.Valid([]byte(line)) {
			jsonCount++
		}
	}
	if err = scanner.Err(); err != nil {
		return "", nil, err
	}

	// Drain remainder into buf so replay is complete.
	if _, err = io.Copy(&buf, r); err != nil {
		return "", nil, err
	}

	if sampled > 0 && jsonCount == sampled {
		format = "json"
	} else {
		format = "plaintext"
	}
	return format, &buf, nil
}
