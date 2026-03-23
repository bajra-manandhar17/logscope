package analyzer

import (
	"testing"
	"time"
)

var t0 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func makeEntries(base time.Time, offsets []time.Duration, levels []string) []LogEntry {
	entries := make([]LogEntry, len(offsets))
	for i, d := range offsets {
		entries[i] = LogEntry{Timestamp: base.Add(d), Level: levels[i]}
	}
	return entries
}

// --- Interval selection ---

func TestBucketInterval_LessThan1Hour(t *testing.T) {
	entries := makeEntries(t0, []time.Duration{0, 30 * time.Minute}, []string{"info", "info"})
	_, interval := BucketTimeSeries(entries)
	if interval != "1m" {
		t.Errorf("got %q, want \"1m\"", interval)
	}
}

func TestBucketInterval_LessThan24Hours(t *testing.T) {
	entries := makeEntries(t0, []time.Duration{0, 2 * time.Hour}, []string{"info", "info"})
	_, interval := BucketTimeSeries(entries)
	if interval != "15m" {
		t.Errorf("got %q, want \"15m\"", interval)
	}
}

func TestBucketInterval_LessThan7Days(t *testing.T) {
	entries := makeEntries(t0, []time.Duration{0, 3 * 24 * time.Hour}, []string{"info", "info"})
	_, interval := BucketTimeSeries(entries)
	if interval != "1h" {
		t.Errorf("got %q, want \"1h\"", interval)
	}
}

func TestBucketInterval_7DaysOrMore(t *testing.T) {
	entries := makeEntries(t0, []time.Duration{0, 7 * 24 * time.Hour}, []string{"info", "info"})
	_, interval := BucketTimeSeries(entries)
	if interval != "1d" {
		t.Errorf("got %q, want \"1d\"", interval)
	}
}

// --- Bucketing logic ---

func TestBucketTimeSeries_CountsPerBucket(t *testing.T) {
	entries := makeEntries(t0, []time.Duration{
		0, 30 * time.Second, 1 * time.Minute, 2 * time.Minute,
	}, []string{"info", "info", "info", "info"})

	buckets, _ := BucketTimeSeries(entries)
	// 1m interval: entries at 0s and 30s → bucket 0; 1m → bucket 1; 2m → bucket 2
	if len(buckets) < 3 {
		t.Fatalf("want >=3 buckets, got %d", len(buckets))
	}
	if buckets[0].Count != 2 {
		t.Errorf("bucket[0] count: got %d, want 2", buckets[0].Count)
	}
}

func TestBucketTimeSeries_ErrorCount(t *testing.T) {
	entries := makeEntries(t0, []time.Duration{
		0, 10 * time.Second, 20 * time.Second,
	}, []string{"error", "info", "error"})

	buckets, _ := BucketTimeSeries(entries)
	if buckets[0].ErrorCount != 2 {
		t.Errorf("errorCount: got %d, want 2", buckets[0].ErrorCount)
	}
}

func TestBucketTimeSeries_SortedChronologically(t *testing.T) {
	entries := makeEntries(t0, []time.Duration{
		5 * time.Minute, 0, 2 * time.Minute,
	}, []string{"info", "info", "info"})

	buckets, _ := BucketTimeSeries(entries)
	for i := 1; i < len(buckets); i++ {
		if buckets[i].Timestamp.Before(buckets[i-1].Timestamp) {
			t.Errorf("buckets not sorted at index %d", i)
		}
	}
}

func TestBucketTimeSeries_SkipsZeroTimestamps(t *testing.T) {
	entries := []LogEntry{
		{Timestamp: time.Time{}, Level: "info"},
		{Timestamp: t0, Level: "error"},
		{Timestamp: time.Time{}, Level: "warn"},
	}
	buckets, interval := BucketTimeSeries(entries)
	if interval != "1m" {
		t.Errorf("interval: got %q", interval)
	}
	total := 0
	for _, b := range buckets {
		total += b.Count
	}
	if total != 1 {
		t.Errorf("total count: got %d, want 1 (zero-ts entries skipped)", total)
	}
}

func TestBucketTimeSeries_Empty(t *testing.T) {
	buckets, interval := BucketTimeSeries(nil)
	if len(buckets) != 0 {
		t.Errorf("want 0 buckets")
	}
	if interval != "1m" {
		t.Errorf("default interval: got %q", interval)
	}
}

func TestBucketTimeSeries_BucketTimestampIsStart(t *testing.T) {
	// Entry at 1m30s should land in the bucket starting at 1m.
	entries := makeEntries(t0, []time.Duration{90 * time.Second}, []string{"info"})
	buckets, _ := BucketTimeSeries(entries)
	want := t0.Add(1 * time.Minute)
	if !buckets[0].Timestamp.Equal(want) {
		t.Errorf("bucket start: got %v, want %v", buckets[0].Timestamp, want)
	}
}
