package analyzer

import (
	"sort"
	"time"
)

const (
	topSourcesN   = 10
	maxSourceKeys = 1000 // cap source map to stay O(1) memory
)

// SummaryAggregator accumulates stats incrementally from LogEntry values.
type SummaryAggregator struct {
	totalLines int
	errorCount int
	warnCount  int
	infoCount  int
	debugCount int
	minTS      time.Time
	maxTS      time.Time
	sources    map[string]int
	hasSource  bool
}

func NewSummaryAggregator() *SummaryAggregator {
	return &SummaryAggregator{sources: make(map[string]int)}
}

// Add incorporates one entry into the running totals.
func (a *SummaryAggregator) Add(e LogEntry) {
	a.totalLines++

	switch e.Level {
	case "error":
		a.errorCount++
	case "warn":
		a.warnCount++
	case "info":
		a.infoCount++
	case "debug":
		a.debugCount++
	}

	if !e.Timestamp.IsZero() {
		if a.minTS.IsZero() || e.Timestamp.Before(a.minTS) {
			a.minTS = e.Timestamp
		}
		if e.Timestamp.After(a.maxTS) {
			a.maxTS = e.Timestamp
		}
	}

	if e.Source != "" {
		a.hasSource = true
		if len(a.sources) < maxSourceKeys {
			a.sources[e.Source]++
		} else if _, exists := a.sources[e.Source]; exists {
			a.sources[e.Source]++
		}
		// If map is full and source is new, silently drop to keep memory bounded.
	}
}

// Result returns the aggregated Summary.
func (a *SummaryAggregator) Result() Summary {
	s := Summary{
		TotalLines: a.totalLines,
		ErrorCount: a.errorCount,
		WarnCount:  a.warnCount,
		InfoCount:  a.infoCount,
		DebugCount: a.debugCount,
		TimeRange:  [2]time.Time{a.minTS, a.maxTS},
	}

	if a.hasSource {
		type kv struct {
			key   string
			count int
		}
		ranked := make([]kv, 0, len(a.sources))
		for k, v := range a.sources {
			ranked = append(ranked, kv{k, v})
		}
		sort.Slice(ranked, func(i, j int) bool {
			return ranked[i].count > ranked[j].count
		})
		n := topSourcesN
		if len(ranked) < n {
			n = len(ranked)
		}
		s.TopSources = make([]string, n)
		for i := range n {
			s.TopSources[i] = ranked[i].key
		}
	}

	return s
}
