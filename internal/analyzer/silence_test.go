package analyzer

import (
	"testing"
	"time"
)

func TestDetectSilenceGaps_Empty(t *testing.T) {
	if gaps := DetectSilenceGaps(nil, nil, nil, time.Minute); gaps != nil {
		t.Errorf("want nil, got %v", gaps)
	}
}

func TestDetectSilenceGaps_SingleSource(t *testing.T) {
	// Need ≥2 sources to detect "others active" — single source returns nil.
	base := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	buckets := []TimeBucket{
		{Timestamp: base, Count: 5},
		{Timestamp: base.Add(time.Minute), Count: 5},
		{Timestamp: base.Add(2 * time.Minute), Count: 5},
	}
	entries := []LogEntry{
		{Timestamp: base, Source: "api", Message: "ok"},
	}
	if gaps := DetectSilenceGaps(entries, buckets, []string{"api"}, time.Minute); gaps != nil {
		t.Errorf("want nil for single source, got %v", gaps)
	}
}

func TestDetectSilenceGaps_AllActive(t *testing.T) {
	base := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	buckets := []TimeBucket{
		{Timestamp: base, Count: 10},
		{Timestamp: base.Add(time.Minute), Count: 10},
		{Timestamp: base.Add(2 * time.Minute), Count: 10},
	}
	entries := []LogEntry{
		{Timestamp: base, Source: "api", Message: "a"},
		{Timestamp: base.Add(time.Minute), Source: "api", Message: "b"},
		{Timestamp: base.Add(2 * time.Minute), Source: "api", Message: "c"},
		{Timestamp: base, Source: "db", Message: "d"},
		{Timestamp: base.Add(time.Minute), Source: "db", Message: "e"},
		{Timestamp: base.Add(2 * time.Minute), Source: "db", Message: "f"},
	}
	gaps := DetectSilenceGaps(entries, buckets, []string{"api", "db"}, time.Minute)
	if len(gaps) != 0 {
		t.Errorf("want 0 gaps when all sources active, got %d", len(gaps))
	}
}

func TestDetectSilenceGaps_GapDetected(t *testing.T) {
	base := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	buckets := make([]TimeBucket, 6)
	for i := range buckets {
		buckets[i] = TimeBucket{
			Timestamp: base.Add(time.Duration(i) * time.Minute),
			Count:     5,
		}
	}

	entries := []LogEntry{
		// "api" active in all buckets.
		{Timestamp: base, Source: "api", Message: "a"},
		{Timestamp: base.Add(time.Minute), Source: "api", Message: "b"},
		{Timestamp: base.Add(2 * time.Minute), Source: "api", Message: "c"},
		{Timestamp: base.Add(3 * time.Minute), Source: "api", Message: "d"},
		{Timestamp: base.Add(4 * time.Minute), Source: "api", Message: "e"},
		{Timestamp: base.Add(5 * time.Minute), Source: "api", Message: "f"},
		// "db" active only in first and last buckets — silent for buckets 1-4.
		{Timestamp: base, Source: "db", Message: "g"},
		{Timestamp: base.Add(5 * time.Minute), Source: "db", Message: "h"},
	}

	gaps := DetectSilenceGaps(entries, buckets, []string{"api", "db"}, time.Minute)
	if len(gaps) == 0 {
		t.Fatal("expected at least one gap for 'db'")
	}

	found := false
	for _, g := range gaps {
		if g.Source == "db" {
			found = true
			if g.ActiveSourcesDuringGap == 0 {
				t.Error("expected active sources > 0 during gap")
			}
		}
	}
	if !found {
		t.Error("expected gap for source 'db'")
	}
}

func TestDetectSilenceGaps_GapAtEnd(t *testing.T) {
	base := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	buckets := make([]TimeBucket, 5)
	for i := range buckets {
		buckets[i] = TimeBucket{
			Timestamp: base.Add(time.Duration(i) * time.Minute),
			Count:     5,
		}
	}

	entries := []LogEntry{
		// "api" active in all buckets.
		{Timestamp: base, Source: "api", Message: "a"},
		{Timestamp: base.Add(time.Minute), Source: "api", Message: "b"},
		{Timestamp: base.Add(2 * time.Minute), Source: "api", Message: "c"},
		{Timestamp: base.Add(3 * time.Minute), Source: "api", Message: "d"},
		{Timestamp: base.Add(4 * time.Minute), Source: "api", Message: "e"},
		// "db" only in first bucket — silent for rest.
		{Timestamp: base, Source: "db", Message: "f"},
	}

	gaps := DetectSilenceGaps(entries, buckets, []string{"api", "db"}, time.Minute)
	if len(gaps) == 0 {
		t.Fatal("expected gap at end for 'db'")
	}
	found := false
	for _, g := range gaps {
		if g.Source == "db" {
			found = true
		}
	}
	if !found {
		t.Error("expected gap for source 'db'")
	}
}
