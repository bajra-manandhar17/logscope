package analyzer

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

var jsonLogs = `{"timestamp":"2024-01-01T00:00:00Z","level":"info","message":"server started","source":"main"}
{"timestamp":"2024-01-01T00:01:00Z","level":"error","message":"connection refused 127.0.0.1","source":"db"}
{"timestamp":"2024-01-01T00:02:00Z","level":"warn","message":"slow query 450ms","source":"db"}
{"timestamp":"2024-01-01T00:03:00Z","level":"info","message":"request ok","source":"api"}
{"timestamp":"2024-01-01T00:04:00Z","level":"error","message":"connection refused 10.0.0.1","source":"db"}
`

var plaintextLogs = `2024-01-01T00:00:00Z INFO [main] server started
2024-01-01T00:01:00Z ERROR [db] connection refused 127.0.0.1
2024-01-01T00:02:00Z WARN [db] slow query 450ms
2024-01-01T00:03:00Z INFO [api] request ok
2024-01-01T00:04:00Z ERROR [db] connection refused 10.0.0.1
`

func TestAnalyze_JSONInput(t *testing.T) {
	result, err := Analyze(context.Background(), strings.NewReader(jsonLogs), "")
	if err != nil {
		t.Fatal(err)
	}

	if result.FormatDetected != "json" {
		t.Errorf("format: got %q", result.FormatDetected)
	}
	if result.Summary.TotalLines != 5 {
		t.Errorf("totalLines: got %d", result.Summary.TotalLines)
	}
	if result.Summary.ErrorCount != 2 {
		t.Errorf("errorCount: got %d", result.Summary.ErrorCount)
	}
	if len(result.Entries) != 5 {
		t.Errorf("entries: got %d", len(result.Entries))
	}
	if len(result.Patterns) == 0 {
		t.Error("patterns empty")
	}
	if len(result.TimeSeries) == 0 {
		t.Error("timeSeries empty")
	}
	if result.BucketInterval == "" {
		t.Error("bucketInterval empty")
	}
}

func TestAnalyze_PlaintextInput(t *testing.T) {
	result, err := Analyze(context.Background(), strings.NewReader(plaintextLogs), "")
	if err != nil {
		t.Fatal(err)
	}

	if result.FormatDetected != "plaintext" {
		t.Errorf("format: got %q", result.FormatDetected)
	}
	if result.Summary.TotalLines != 5 {
		t.Errorf("totalLines: got %d", result.Summary.TotalLines)
	}
	if result.Summary.ErrorCount != 2 {
		t.Errorf("errorCount: got %d", result.Summary.ErrorCount)
	}
}

func TestAnalyze_FormatHint(t *testing.T) {
	// Provide plaintext data but hint "plaintext" — should skip detection.
	result, err := Analyze(context.Background(), strings.NewReader(plaintextLogs), "plaintext")
	if err != nil {
		t.Fatal(err)
	}
	if result.FormatDetected != "plaintext" {
		t.Errorf("format: got %q", result.FormatDetected)
	}
}

func TestAnalyze_EntryCap(t *testing.T) {
	var sb strings.Builder
	ts := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < MaxEntries+100; i++ {
		sb.WriteString(fmt.Sprintf(
			"{\"timestamp\":%q,\"level\":\"info\",\"message\":\"line %d\"}\n",
			ts.Add(time.Duration(i)*time.Second).Format(time.RFC3339), i,
		))
	}

	result, err := Analyze(context.Background(), strings.NewReader(sb.String()), "json")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Entries) > MaxEntries {
		t.Errorf("entries not capped: got %d", len(result.Entries))
	}
	// Summary should reflect all lines, not just capped entries.
	if result.Summary.TotalLines != MaxEntries+100 {
		t.Errorf("summary totalLines: got %d, want %d", result.Summary.TotalLines, MaxEntries+100)
	}
}

func TestAnalyze_ContextCancellation(t *testing.T) {
	var sb strings.Builder
	for i := 0; i < 10000; i++ {
		sb.WriteString("{\"level\":\"info\",\"message\":\"line\"}\n")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	_, err := Analyze(ctx, strings.NewReader(sb.String()), "json")
	if err == nil {
		t.Error("expected error on cancelled context")
	}
}

func TestAnalyze_EmptyInput(t *testing.T) {
	result, err := Analyze(context.Background(), strings.NewReader(""), "")
	if err != nil {
		t.Fatal(err)
	}
	if result.Summary.TotalLines != 0 {
		t.Errorf("totalLines: got %d", result.Summary.TotalLines)
	}
	if len(result.Entries) != 0 {
		t.Errorf("entries: got %d", len(result.Entries))
	}
}

func TestAnalyze_PatternsGrouped(t *testing.T) {
	// Two messages that differ only in IP — should collapse to one pattern.
	input := `{"level":"error","message":"connection refused 127.0.0.1"}
{"level":"error","message":"connection refused 10.0.0.2"}
`
	result, err := Analyze(context.Background(), strings.NewReader(input), "json")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Patterns) != 1 {
		t.Errorf("want 1 pattern, got %d", len(result.Patterns))
	}
	if result.Patterns[0].Count != 2 {
		t.Errorf("pattern count: got %d", result.Patterns[0].Count)
	}
}
