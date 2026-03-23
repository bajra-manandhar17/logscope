package analyzer

import (
	"context"
	"strings"
	"testing"
	"time"
)

// collect drains the Parse channel into a slice.
func collect(ctx context.Context, r *strings.Reader, format string) ([]LogEntry, error) {
	entries, errc := Parse(ctx, r, format)
	var out []LogEntry
	for e := range entries {
		out = append(out, e)
	}
	return out, <-errc
}

// --- JSON ---

func TestParseJSON_BasicFields(t *testing.T) {
	input := `{"timestamp":"2024-01-02T03:04:05Z","level":"error","message":"boom","source":"app.go"}
{"timestamp":"2024-01-02T03:04:06Z","level":"info","message":"ok","service":"svc"}
`
	entries, err := collect(context.Background(), strings.NewReader(input), "json")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("want 2 entries, got %d", len(entries))
	}

	e := entries[0]
	if e.Level != "error" {
		t.Errorf("level: got %q", e.Level)
	}
	if e.Message != "boom" {
		t.Errorf("message: got %q", e.Message)
	}
	if e.Source != "app.go" {
		t.Errorf("source: got %q", e.Source)
	}
	if e.LineNumber != 1 {
		t.Errorf("line_number: got %d", e.LineNumber)
	}

	// second entry uses "service" as source field
	if entries[1].Source != "svc" {
		t.Errorf("source from 'service': got %q", entries[1].Source)
	}
}

func TestParseJSON_SourceFieldAliases(t *testing.T) {
	lines := []struct {
		field string
		value string
	}{
		{"source", "s1"},
		{"service", "s2"},
		{"module", "s3"},
		{"logger", "s4"},
		{"component", "s5"},
	}
	for _, tc := range lines {
		input := `{"level":"info","message":"x","` + tc.field + `":"` + tc.value + `"}`
		entries, err := collect(context.Background(), strings.NewReader(input), "json")
		if err != nil {
			t.Fatal(err)
		}
		if len(entries) != 1 {
			t.Fatalf("%s: want 1 entry", tc.field)
		}
		if entries[0].Source != tc.value {
			t.Errorf("%s: got source %q, want %q", tc.field, entries[0].Source, tc.value)
		}
	}
}

func TestParseJSON_MalformedLineSkipped(t *testing.T) {
	input := `{"level":"info","message":"good"}
not json at all
{"level":"warn","message":"also good"}
`
	entries, err := collect(context.Background(), strings.NewReader(input), "json")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Errorf("want 2, got %d", len(entries))
	}
}

func TestParseJSON_EmptyLinesSkipped(t *testing.T) {
	input := "\n{\"level\":\"info\",\"message\":\"hi\"}\n\n"
	entries, err := collect(context.Background(), strings.NewReader(input), "json")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Errorf("want 1, got %d", len(entries))
	}
}

func TestParseJSON_TimestampParsed(t *testing.T) {
	input := `{"timestamp":"2024-06-15T10:20:30Z","level":"info","message":"ts test"}`
	entries, err := collect(context.Background(), strings.NewReader(input), "json")
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2024, 6, 15, 10, 20, 30, 0, time.UTC)
	if !entries[0].Timestamp.Equal(want) {
		t.Errorf("timestamp: got %v, want %v", entries[0].Timestamp, want)
	}
}

// --- Plaintext ---

func TestParsePlaintext_BasicFields(t *testing.T) {
	input := "2024-01-02T03:04:05Z ERROR [myservice] something failed\n2024-01-02T03:04:06Z INFO [api] request ok\n"
	entries, err := collect(context.Background(), strings.NewReader(input), "plaintext")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("want 2, got %d", len(entries))
	}
	e := entries[0]
	if e.Level != "error" {
		t.Errorf("level: got %q", e.Level)
	}
	if e.Source != "myservice" {
		t.Errorf("source: got %q", e.Source)
	}
	if e.Message != "something failed" {
		t.Errorf("message: got %q", e.Message)
	}
}

func TestParsePlaintext_NoSourceOrTimestamp(t *testing.T) {
	input := "ERROR something broke\nINFO all good\n"
	entries, err := collect(context.Background(), strings.NewReader(input), "plaintext")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("want 2, got %d", len(entries))
	}
	if entries[0].Level != "error" {
		t.Errorf("level: got %q", entries[0].Level)
	}
}

func TestParsePlaintext_RawAndLineNumber(t *testing.T) {
	input := "ERROR first\nINFO second\n"
	entries, err := collect(context.Background(), strings.NewReader(input), "plaintext")
	if err != nil {
		t.Fatal(err)
	}
	if entries[0].Raw != "ERROR first" {
		t.Errorf("raw: got %q", entries[0].Raw)
	}
	if entries[0].LineNumber != 1 {
		t.Errorf("line_number: got %d", entries[0].LineNumber)
	}
	if entries[1].LineNumber != 2 {
		t.Errorf("line_number: got %d", entries[1].LineNumber)
	}
}

// --- Context cancellation ---

func TestParse_ContextCancellation(t *testing.T) {
	var sb strings.Builder
	for i := 0; i < 1000; i++ {
		sb.WriteString("{\"level\":\"info\",\"message\":\"line\"}\n")
	}
	ctx, cancel := context.WithCancel(context.Background())
	entries, errc := Parse(ctx, strings.NewReader(sb.String()), "json")

	// Cancel after first entry.
	<-entries
	cancel()

	// Drain remaining.
	for range entries {
	}
	if err := <-errc; err != nil && err != context.Canceled {
		t.Errorf("unexpected error: %v", err)
	}
}
