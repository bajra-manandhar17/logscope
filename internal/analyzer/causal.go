package analyzer

import (
	"sort"
	"time"
)

const (
	maxCandidatePatterns = 100
	causalWindowSeconds  = 60
	minCooccurrence      = 3
	maxCausalResults     = 20
	maxForwardScan       = 500
)

type causalEvent struct {
	idx int
	ts  time.Time
}

type pairKey struct {
	a, b int
}

type pairStats struct {
	count    int
	totalLag float64 // seconds
}

// DetectCausalSequences finds temporal A→B patterns among error/warn log entries.
// Only considers the top maxCandidatePatterns by frequency. A must precede B within
// causalWindowSeconds. Returns top maxCausalResults pairs with ≥minCooccurrence hits.
func DetectCausalSequences(entries []LogEntry, patterns []Pattern) []CausalSequence {
	// Select candidate patterns (error/warn only).
	// We identify error/warn entries by level, then match their masked template
	// to the known pattern list.
	candidates := selectCandidates(patterns, entries)
	if len(candidates) < 2 {
		return nil
	}

	templateToIdx := make(map[string]int, len(candidates))
	for i, tmpl := range candidates {
		templateToIdx[tmpl] = i
	}

	// Collect timestamped events for candidate patterns.
	var events []causalEvent
	for _, e := range entries {
		if e.Timestamp.IsZero() {
			continue
		}
		if e.Level != "error" && e.Level != "warn" {
			continue
		}
		tmpl := maskTokens(e.Message)
		if idx, ok := templateToIdx[tmpl]; ok {
			events = append(events, causalEvent{idx: idx, ts: e.Timestamp})
		}
	}

	if len(events) < 2 {
		return nil
	}

	// Sort by timestamp.
	sort.Slice(events, func(i, j int) bool {
		return events[i].ts.Before(events[j].ts)
	})

	window := time.Duration(causalWindowSeconds) * time.Second
	pairs := make(map[pairKey]*pairStats)

	for i, ev := range events {
		scanned := 0
		for j := i + 1; j < len(events) && scanned < maxForwardScan; j++ {
			lag := events[j].ts.Sub(ev.ts)
			if lag > window {
				break
			}
			if lag <= 0 || events[j].idx == ev.idx {
				continue
			}
			scanned++
			key := pairKey{a: ev.idx, b: events[j].idx}
			ps := pairs[key]
			if ps == nil {
				ps = &pairStats{}
				pairs[key] = ps
			}
			ps.count++
			ps.totalLag += lag.Seconds()
		}
	}

	// Filter and build results.
	var results []CausalSequence
	for key, ps := range pairs {
		if ps.count < minCooccurrence {
			continue
		}
		results = append(results, CausalSequence{
			PatternA:      candidates[key.a],
			PatternB:      candidates[key.b],
			Count:         ps.count,
			AvgLagSeconds: ps.totalLag / float64(ps.count),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Count > results[j].Count
	})

	if len(results) > maxCausalResults {
		results = results[:maxCausalResults]
	}
	return results
}

// selectCandidates picks the top N pattern templates that appear in error/warn entries.
func selectCandidates(patterns []Pattern, entries []LogEntry) []string {
	// Build set of templates that have error/warn entries.
	errorWarnTemplates := make(map[string]int)
	for _, e := range entries {
		if e.Level != "error" && e.Level != "warn" {
			continue
		}
		if e.Message == "" {
			continue
		}
		tmpl := maskTokens(e.Message)
		errorWarnTemplates[tmpl]++
	}

	// Intersect with known patterns (which are already sorted by count desc).
	// Take top maxCandidatePatterns.
	var candidates []string
	for _, p := range patterns {
		if _, ok := errorWarnTemplates[p.Template]; ok {
			candidates = append(candidates, p.Template)
			if len(candidates) >= maxCandidatePatterns {
				break
			}
		}
	}
	return candidates
}
