package analyzer

import (
	"testing"
	"time"
)

func makeEntry(level, source string, ts time.Time) LogEntry {
	return LogEntry{Level: level, Source: source, Timestamp: ts}
}

func TestSummary_LevelCounts(t *testing.T) {
	a := NewSummaryAggregator()
	a.Add(makeEntry("error", "", time.Time{}))
	a.Add(makeEntry("error", "", time.Time{}))
	a.Add(makeEntry("warn", "", time.Time{}))
	a.Add(makeEntry("info", "", time.Time{}))
	a.Add(makeEntry("debug", "", time.Time{}))

	s := a.Result()
	if s.TotalLines != 5 {
		t.Errorf("TotalLines: got %d", s.TotalLines)
	}
	if s.ErrorCount != 2 {
		t.Errorf("ErrorCount: got %d", s.ErrorCount)
	}
	if s.WarnCount != 1 {
		t.Errorf("WarnCount: got %d", s.WarnCount)
	}
	if s.InfoCount != 1 {
		t.Errorf("InfoCount: got %d", s.InfoCount)
	}
	if s.DebugCount != 1 {
		t.Errorf("DebugCount: got %d", s.DebugCount)
	}
}

func TestSummary_TimeRange(t *testing.T) {
	t1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	t3 := time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)

	a := NewSummaryAggregator()
	a.Add(makeEntry("info", "", t2))
	a.Add(makeEntry("info", "", t1))
	a.Add(makeEntry("info", "", t3))

	s := a.Result()
	if !s.TimeRange[0].Equal(t1) {
		t.Errorf("TimeRange[0]: got %v, want %v", s.TimeRange[0], t1)
	}
	if !s.TimeRange[1].Equal(t3) {
		t.Errorf("TimeRange[1]: got %v, want %v", s.TimeRange[1], t3)
	}
}

func TestSummary_TopSources(t *testing.T) {
	a := NewSummaryAggregator()
	for i := 0; i < 5; i++ {
		a.Add(makeEntry("info", "alpha", time.Time{}))
	}
	for i := 0; i < 3; i++ {
		a.Add(makeEntry("info", "beta", time.Time{}))
	}
	a.Add(makeEntry("info", "gamma", time.Time{}))

	s := a.Result()
	if len(s.TopSources) == 0 {
		t.Fatal("TopSources empty")
	}
	if s.TopSources[0] != "alpha" {
		t.Errorf("top source: got %q, want \"alpha\"", s.TopSources[0])
	}
	if s.TopSources[1] != "beta" {
		t.Errorf("second source: got %q, want \"beta\"", s.TopSources[1])
	}
}

func TestSummary_TopSourcesCappedAt10(t *testing.T) {
	a := NewSummaryAggregator()
	for i := 0; i < 20; i++ {
		src := string(rune('a' + i))
		a.Add(makeEntry("info", src, time.Time{}))
	}
	s := a.Result()
	if len(s.TopSources) > 10 {
		t.Errorf("TopSources len: got %d, want <=10", len(s.TopSources))
	}
}

func TestSummary_NoSourcesOmitsTopSources(t *testing.T) {
	a := NewSummaryAggregator()
	a.Add(makeEntry("info", "", time.Time{}))
	a.Add(makeEntry("error", "", time.Time{}))

	s := a.Result()
	if s.TopSources != nil {
		t.Errorf("expected nil TopSources when no source found, got %v", s.TopSources)
	}
}

func TestSummary_MissingTimestampsIgnored(t *testing.T) {
	t1 := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	a := NewSummaryAggregator()
	a.Add(makeEntry("info", "", time.Time{})) // zero — should be ignored
	a.Add(makeEntry("info", "", t1))

	s := a.Result()
	if !s.TimeRange[0].Equal(t1) {
		t.Errorf("TimeRange[0]: got %v, want %v", s.TimeRange[0], t1)
	}
	if !s.TimeRange[1].Equal(t1) {
		t.Errorf("TimeRange[1]: got %v, want %v", s.TimeRange[1], t1)
	}
}

func TestSummary_UnknownLevelStillCountedInTotal(t *testing.T) {
	a := NewSummaryAggregator()
	a.Add(makeEntry("trace", "", time.Time{}))
	s := a.Result()
	if s.TotalLines != 1 {
		t.Errorf("TotalLines: got %d", s.TotalLines)
	}
}

func TestSummary_Empty(t *testing.T) {
	a := NewSummaryAggregator()
	s := a.Result()
	if s.TotalLines != 0 {
		t.Errorf("TotalLines: got %d", s.TotalLines)
	}
	if s.TopSources != nil {
		t.Errorf("TopSources should be nil")
	}
}

func TestSummary_SourceMapBoundedBeforeResult(t *testing.T) {
	// Feed more than maxSourceKeys distinct sources — map must stay bounded.
	a := NewSummaryAggregator()
	for i := 0; i < 2000; i++ {
		src := string(rune(i + 0x4E00)) // unique CJK characters
		a.Add(makeEntry("info", src, time.Time{}))
	}
	s := a.Result()
	if len(s.TopSources) > 10 {
		t.Errorf("TopSources len: got %d", len(s.TopSources))
	}
}
