package generator

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

var baseConfig = GenerateConfig{
	Format:     "json",
	TotalLines: 500,
	Levels: map[string]float64{
		"error": 0.05,
		"warn":  0.15,
		"info":  0.70,
		"debug": 0.10,
	},
	Start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	End:   time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
}

func collectLines(t *testing.T, cfg GenerateConfig) []string {
	t.Helper()
	ctx := context.Background()
	batches, errc := Generate(ctx, cfg)
	var lines []string
	for b := range batches {
		lines = append(lines, b...)
	}
	if err := <-errc; err != nil {
		t.Fatalf("generator error: %v", err)
	}
	return lines
}

func TestGenerateLineCount(t *testing.T) {
	lines := collectLines(t, baseConfig)
	if len(lines) != baseConfig.TotalLines {
		t.Errorf("got %d lines, want %d", len(lines), baseConfig.TotalLines)
	}
}

func TestGenerateBatchSize(t *testing.T) {
	cfg := baseConfig
	cfg.TotalLines = 350

	ctx := context.Background()
	batches, errc := Generate(ctx, cfg)

	total := 0
	for b := range batches {
		if len(b) > batchSize {
			t.Errorf("batch size %d exceeds max %d", len(b), batchSize)
		}
		total += len(b)
	}
	if err := <-errc; err != nil {
		t.Fatalf("generator error: %v", err)
	}
	if total != cfg.TotalLines {
		t.Errorf("total lines %d != %d", total, cfg.TotalLines)
	}
}

func TestGenerateLevelDistribution(t *testing.T) {
	cfg := baseConfig
	cfg.TotalLines = 10_000

	lines := collectLines(t, cfg)

	counts := map[string]int{}
	for _, l := range lines {
		var m map[string]string
		if err := json.Unmarshal([]byte(l), &m); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
		counts[m["level"]]++
	}

	total := float64(len(lines))
	for level, want := range cfg.Levels {
		got := float64(counts[level]) / total
		if diff := got - want; diff < -0.05 || diff > 0.05 {
			t.Errorf("level %s: got %.3f, want %.3f (tolerance ±0.05)", level, got, want)
		}
	}
}

func TestGenerateTimestampsInRange(t *testing.T) {
	lines := collectLines(t, baseConfig)
	for _, l := range lines {
		var m map[string]string
		if err := json.Unmarshal([]byte(l), &m); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
		ts, err := time.Parse(time.RFC3339, m["timestamp"])
		if err != nil {
			t.Fatalf("bad timestamp %q: %v", m["timestamp"], err)
		}
		if ts.Before(baseConfig.Start) || !ts.Before(baseConfig.End) {
			t.Errorf("timestamp %v out of range [%v, %v)", ts, baseConfig.Start, baseConfig.End)
		}
	}
}

func TestGenerateJSONParseable(t *testing.T) {
	lines := collectLines(t, baseConfig)
	for i, l := range lines {
		var m map[string]string
		if err := json.Unmarshal([]byte(l), &m); err != nil {
			t.Errorf("line %d not valid JSON: %v", i, err)
		}
		for _, field := range []string{"timestamp", "level", "message", "source"} {
			if m[field] == "" {
				t.Errorf("line %d missing field %q", i, field)
			}
		}
	}
}

func TestGeneratePlaintextFormat(t *testing.T) {
	cfg := baseConfig
	cfg.Format = "plaintext"
	lines := collectLines(t, cfg)

	for i, l := range lines {
		// Must contain a recognisable level token.
		hasLevel := strings.Contains(l, "ERROR") || strings.Contains(l, "WARN") ||
			strings.Contains(l, "INFO") || strings.Contains(l, "DEBUG")
		if !hasLevel {
			t.Errorf("line %d has no level: %s", i, l)
		}
		// Must contain a timestamp-shaped prefix.
		if len(l) < 20 {
			t.Errorf("line %d too short: %s", i, l)
		}
	}
}

func TestGenerateContextCancellation(t *testing.T) {
	cfg := baseConfig
	cfg.TotalLines = 1_000_000

	ctx, cancel := context.WithCancel(context.Background())
	batches, errc := Generate(ctx, cfg)

	// Drain one batch then cancel.
	<-batches
	cancel()

	// Drain remaining batches (must not block).
	for range batches {
	}

	err := <-errc
	if err == nil {
		t.Error("expected cancellation error, got nil")
	}
}

func TestGenerateSmallCount(t *testing.T) {
	cfg := baseConfig
	cfg.TotalLines = 1
	lines := collectLines(t, cfg)
	if len(lines) != 1 {
		t.Errorf("got %d lines, want 1", len(lines))
	}
}
