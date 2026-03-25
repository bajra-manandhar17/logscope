package analyzer

import (
	"testing"
	"time"
)

func TestDetectCausalSequences_Empty(t *testing.T) {
	if seqs := DetectCausalSequences(nil, nil); seqs != nil {
		t.Errorf("want nil, got %v", seqs)
	}
}

func TestDetectCausalSequences_NoErrorWarn(t *testing.T) {
	base := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	entries := []LogEntry{
		{Timestamp: base, Level: "info", Message: "all good"},
		{Timestamp: base.Add(10 * time.Second), Level: "debug", Message: "trace"},
	}
	patterns := []Pattern{
		{Template: "all good", Count: 1},
		{Template: "trace", Count: 1},
	}
	if seqs := DetectCausalSequences(entries, patterns); seqs != nil {
		t.Errorf("want nil for info/debug only, got %v", seqs)
	}
}

func TestDetectCausalSequences_DetectsAB(t *testing.T) {
	base := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	// Pattern A: "connection timeout 123ms" → "connection timeout {NUM}ms"
	// Pattern B: "circuit breaker open" → "circuit breaker open"
	// A→B happens 4 times within 60s.
	var entries []LogEntry
	for i := 0; i < 4; i++ {
		offset := time.Duration(i) * 2 * time.Minute
		entries = append(entries,
			LogEntry{
				Timestamp: base.Add(offset),
				Level:     "error",
				Message:   "connection timeout 100ms",
			},
			LogEntry{
				Timestamp: base.Add(offset + 30*time.Second),
				Level:     "error",
				Message:   "circuit breaker open",
			},
		)
	}

	patterns := []Pattern{
		{Template: "connection timeout {NUM}ms", Count: 4},
		{Template: "circuit breaker open", Count: 4},
	}

	seqs := DetectCausalSequences(entries, patterns)
	if len(seqs) == 0 {
		t.Fatal("expected at least one causal sequence")
	}

	found := false
	for _, s := range seqs {
		if s.PatternA == "connection timeout {NUM}ms" && s.PatternB == "circuit breaker open" {
			found = true
			if s.Count != 4 {
				t.Errorf("want count 4, got %d", s.Count)
			}
			if s.AvgLagSeconds < 29 || s.AvgLagSeconds > 31 {
				t.Errorf("want avg lag ~30s, got %.1f", s.AvgLagSeconds)
			}
		}
	}
	if !found {
		t.Error("expected A→B sequence not found")
	}
}

func TestDetectCausalSequences_BelowThreshold(t *testing.T) {
	base := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	// A→B only 2 times — below minCooccurrence of 3.
	entries := []LogEntry{
		{Timestamp: base, Level: "error", Message: "timeout 100ms"},
		{Timestamp: base.Add(10 * time.Second), Level: "error", Message: "breaker open"},
		{Timestamp: base.Add(2 * time.Minute), Level: "error", Message: "timeout 200ms"},
		{Timestamp: base.Add(2*time.Minute + 10*time.Second), Level: "error", Message: "breaker open"},
	}

	patterns := []Pattern{
		{Template: "timeout {NUM}ms", Count: 2},
		{Template: "breaker open", Count: 2},
	}

	seqs := DetectCausalSequences(entries, patterns)
	if len(seqs) != 0 {
		t.Errorf("want 0 sequences (below threshold), got %d", len(seqs))
	}
}

func TestDetectCausalSequences_OutsideWindow(t *testing.T) {
	base := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	// A→B with >60s gap — should not be detected.
	var entries []LogEntry
	for i := 0; i < 5; i++ {
		offset := time.Duration(i) * 5 * time.Minute
		entries = append(entries,
			LogEntry{Timestamp: base.Add(offset), Level: "error", Message: "timeout 100ms"},
			LogEntry{Timestamp: base.Add(offset + 90*time.Second), Level: "error", Message: "breaker open"},
		)
	}

	patterns := []Pattern{
		{Template: "timeout {NUM}ms", Count: 5},
		{Template: "breaker open", Count: 5},
	}

	seqs := DetectCausalSequences(entries, patterns)
	if len(seqs) != 0 {
		t.Errorf("want 0 sequences (outside window), got %d", len(seqs))
	}
}

func TestDetectCausalSequences_ReverseIsSeparate(t *testing.T) {
	base := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	// A→B and B→A both happen 3 times — should be 2 separate sequences.
	var entries []LogEntry
	for i := 0; i < 3; i++ {
		offset := time.Duration(i) * 3 * time.Minute
		entries = append(entries,
			LogEntry{Timestamp: base.Add(offset), Level: "warn", Message: "slow query 500ms"},
			LogEntry{Timestamp: base.Add(offset + 20*time.Second), Level: "warn", Message: "cache miss rate high"},
			LogEntry{Timestamp: base.Add(offset + 40*time.Second), Level: "warn", Message: "slow query 600ms"},
		)
	}

	patterns := []Pattern{
		{Template: "slow query {NUM}ms", Count: 6},
		{Template: "cache miss rate high", Count: 3},
	}

	seqs := DetectCausalSequences(entries, patterns)
	// Should find at least "slow query" → "cache miss" and "cache miss" → "slow query"
	if len(seqs) < 1 {
		t.Errorf("expected causal sequences, got %d", len(seqs))
	}
}
