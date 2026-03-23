package analyzer

import (
	"regexp"
	"sort"
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
	sample string
	count  int
}

// PatternAggregator groups log messages by masked template and counts occurrences.
type PatternAggregator struct {
	templates map[string]*templateEntry
}

func NewPatternAggregator() *PatternAggregator {
	return &PatternAggregator{templates: make(map[string]*templateEntry)}
}

// Add masks the message and increments the matching template's count.
func (a *PatternAggregator) Add(message string) {
	tmpl := maskTokens(message)
	if e, ok := a.templates[tmpl]; ok {
		e.count++
		return
	}
	if len(a.templates) >= MaxPatterns {
		return // cap reached — drop new unique templates
	}
	a.templates[tmpl] = &templateEntry{sample: message, count: 1}
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
