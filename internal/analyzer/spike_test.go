package analyzer

import (
	"testing"
	"time"
)

func TestDetectSpikes_Empty(t *testing.T) {
	if spikes := DetectSpikes(nil); spikes != nil {
		t.Errorf("want nil, got %v", spikes)
	}
}

func TestDetectSpikes_TooFewBuckets(t *testing.T) {
	buckets := []TimeBucket{
		{Count: 100},
		{Count: 200},
	}
	if spikes := DetectSpikes(buckets); spikes != nil {
		t.Errorf("want nil for <3 buckets, got %v", spikes)
	}
}

func TestDetectSpikes_UniformData(t *testing.T) {
	buckets := make([]TimeBucket, 10)
	for i := range buckets {
		buckets[i] = TimeBucket{Count: 50}
	}
	if spikes := DetectSpikes(buckets); spikes != nil {
		t.Errorf("want nil for uniform data, got %v", spikes)
	}
}

func TestDetectSpikes_DetectsBurst(t *testing.T) {
	base := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	buckets := make([]TimeBucket, 20)
	for i := range buckets {
		buckets[i] = TimeBucket{
			Timestamp: base.Add(time.Duration(i) * time.Minute),
			Count:     10,
		}
	}
	// Inject a massive spike.
	buckets[10].Count = 500

	spikes := DetectSpikes(buckets)
	if len(spikes) == 0 {
		t.Fatal("expected at least one spike")
	}
	if !spikes[0].BucketTimestamp.Equal(buckets[10].Timestamp) {
		t.Errorf("spike at wrong bucket: %v", spikes[0].BucketTimestamp)
	}
	if spikes[0].Severity != "high" {
		t.Errorf("expected high severity for 50x spike, got %q", spikes[0].Severity)
	}
}

func TestDetectSpikes_MediumSeverity(t *testing.T) {
	base := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	// Create data where one bucket is between 2σ and 3σ.
	// 9 buckets at 10, 1 bucket at ~40 should be medium.
	buckets := make([]TimeBucket, 10)
	for i := range buckets {
		buckets[i] = TimeBucket{
			Timestamp: base.Add(time.Duration(i) * time.Minute),
			Count:     10,
		}
	}
	buckets[5].Count = 35

	spikes := DetectSpikes(buckets)
	if len(spikes) == 0 {
		t.Fatal("expected a spike")
	}
	found := false
	for _, s := range spikes {
		if s.Severity == "medium" {
			found = true
		}
	}
	if !found {
		t.Error("expected at least one medium severity spike")
	}
}

func TestDetectSpikes_AtThresholdNotFlagged(t *testing.T) {
	// All same counts → stddev=0 → no spikes.
	buckets := []TimeBucket{
		{Count: 10},
		{Count: 10},
		{Count: 10},
	}
	if spikes := DetectSpikes(buckets); spikes != nil {
		t.Errorf("equal counts should not produce spikes, got %v", spikes)
	}
}
