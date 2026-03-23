// Package e2e runs full-stack integration tests against a real httptest.Server.
// No browser required — all flows are exercised via HTTP.
package e2e_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/bajra-manandhar17/logscope-v2/internal/server"
)

var stubFS = fstest.MapFS{
	"dist/index.html": &fstest.MapFile{Data: []byte("<html></html>")},
}

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	srv := server.New(stubFS, false)
	ts := httptest.NewServer(srv.Handler)
	t.Cleanup(ts.Close)
	return ts
}

// multipartUpload builds a multipart/form-data body with a "file" field.
func multipartUpload(t *testing.T, content string) (*bytes.Buffer, string) {
	t.Helper()
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, err := mw.CreateFormFile("file", "test.log")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.WriteString(fw, content); err != nil {
		t.Fatal(err)
	}
	mw.Close()
	return &buf, mw.FormDataContentType()
}

type sseEvent struct{ event, data string }

// parseSSE reads an SSE response body and returns all events.
func parseSSE(body string) []sseEvent {
	var events []sseEvent
	var cur sseEvent
	sc := bufio.NewScanner(strings.NewReader(body))
	for sc.Scan() {
		line := sc.Text()
		switch {
		case strings.HasPrefix(line, "event: "):
			cur.event = strings.TrimPrefix(line, "event: ")
		case strings.HasPrefix(line, "data: "):
			cur.data = strings.TrimPrefix(line, "data: ")
		case line == "":
			if cur.event != "" {
				events = append(events, cur)
			}
			cur = sseEvent{}
		}
	}
	return events
}

// generatePayload is the JSON body for POST /api/generate.
const generatePayload = `{
	"Format": "json",
	"TotalLines": 300,
	"Levels": {"error": 0.05, "warn": 0.15, "info": 0.70, "debug": 0.10},
	"Start": "2024-01-01T00:00:00Z",
	"End":   "2024-01-02T00:00:00Z"
}`

// ── Flow 1: Upload → Analyze ──────────────────────────────────────────────────

// TestE2E_Flow1_JSONLogs uploads a JSON log file and verifies all result fields.
func TestE2E_Flow1_JSONLogs(t *testing.T) {
	ts := newTestServer(t)

	logs := strings.Join([]string{
		`{"timestamp":"2024-01-01T00:00:00Z","level":"info","message":"server started","source":"main"}`,
		`{"timestamp":"2024-01-01T00:01:00Z","level":"error","message":"connection refused 127.0.0.1","source":"db"}`,
		`{"timestamp":"2024-01-01T00:02:00Z","level":"warn","message":"slow query 450ms","source":"db"}`,
		`{"timestamp":"2024-01-01T00:03:00Z","level":"info","message":"request completed","source":"api"}`,
		`{"timestamp":"2024-01-01T00:04:00Z","level":"error","message":"connection refused 10.0.0.1","source":"db"}`,
	}, "\n") + "\n"

	body, ct := multipartUpload(t, logs)
	resp, err := http.Post(ts.URL+"/api/analyze", ct, body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("status %d: %s", resp.StatusCode, b)
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}

	// Summary cards
	summary, ok := result["summary"].(map[string]any)
	if !ok {
		t.Fatal("missing summary")
	}
	if summary["total_lines"].(float64) != 5 {
		t.Errorf("total_lines = %v, want 5", summary["total_lines"])
	}
	if summary["error_count"].(float64) != 2 {
		t.Errorf("error_count = %v, want 2", summary["error_count"])
	}

	// Log table entries
	entries, ok := result["entries"].([]any)
	if !ok || len(entries) != 5 {
		t.Errorf("entries: got %v, want 5", len(entries))
	}

	// Pattern list
	patterns, ok := result["patterns"].([]any)
	if !ok || len(patterns) == 0 {
		t.Error("patterns empty or missing")
	}

	// Time-series chart
	ts2, ok := result["time_series"].([]any)
	if !ok || len(ts2) == 0 {
		t.Error("time_series empty or missing")
	}

	// Bucket interval
	if result["bucket_interval"] == "" {
		t.Error("bucket_interval missing")
	}

	// Format detection
	if result["format_detected"] != "json" {
		t.Errorf("format_detected = %v, want json", result["format_detected"])
	}
}

// TestE2E_Flow1_PlaintextLogs uploads a plaintext log file.
func TestE2E_Flow1_PlaintextLogs(t *testing.T) {
	ts := newTestServer(t)

	logs := `2024-01-01T00:00:00Z INFO [main] server started
2024-01-01T00:01:00Z ERROR [db] connection refused
2024-01-01T00:02:00Z WARN [db] slow query
`
	body, ct := multipartUpload(t, logs)
	resp, err := http.Post(ts.URL+"/api/analyze", ct, body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("status %d: %s", resp.StatusCode, b)
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if result["format_detected"] != "plaintext" {
		t.Errorf("format_detected = %v, want plaintext", result["format_detected"])
	}

	summary := result["summary"].(map[string]any)
	if summary["total_lines"].(float64) != 3 {
		t.Errorf("total_lines = %v, want 3", summary["total_lines"])
	}
}

// TestE2E_Flow1_FilterableFields verifies entries carry the fields needed for
// client-side filtering (level, source, timestamp, message, line_number).
func TestE2E_Flow1_FilterableFields(t *testing.T) {
	ts := newTestServer(t)

	logs := `{"timestamp":"2024-01-01T00:00:00Z","level":"info","message":"hello","source":"svc"}` + "\n"
	body, ct := multipartUpload(t, logs)
	resp, err := http.Post(ts.URL+"/api/analyze", ct, body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result) //nolint:errcheck

	entries := result["entries"].([]any)
	entry := entries[0].(map[string]any)

	for _, field := range []string{"level", "message", "line_number", "raw"} {
		if _, ok := entry[field]; !ok {
			t.Errorf("entry missing field %q", field)
		}
	}
}

// ── Flow 2: Generate → Preview → Verify content ───────────────────────────────

// TestE2E_Flow2_SSEStreamComplete verifies the generator streams batches and
// a done event whose totalLines matches the sum of all batch line counts.
func TestE2E_Flow2_SSEStreamComplete(t *testing.T) {
	ts := newTestServer(t)

	resp, err := http.Post(ts.URL+"/api/generate", "application/json",
		strings.NewReader(generatePayload))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("status %d: %s", resp.StatusCode, b)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "text/event-stream" {
		t.Errorf("Content-Type = %q, want text/event-stream", ct)
	}

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	events := parseSSE(string(rawBody))
	if len(events) < 2 {
		t.Fatalf("expected ≥2 events, got %d", len(events))
	}

	// Collect lines from batch events.
	var allLines []string
	for _, ev := range events[:len(events)-1] {
		if ev.event != "batch" {
			t.Errorf("non-terminal event = %q, want batch", ev.event)
			continue
		}
		var payload map[string][]string
		if err := json.Unmarshal([]byte(ev.data), &payload); err != nil {
			t.Fatalf("batch data invalid JSON: %v", err)
		}
		if len(payload["lines"]) == 0 {
			t.Error("batch event has empty lines")
		}
		allLines = append(allLines, payload["lines"]...)
	}

	// Verify done event.
	done := events[len(events)-1]
	if done.event != "done" {
		t.Fatalf("last event = %q, want done", done.event)
	}
	var donePayload map[string]int
	if err := json.Unmarshal([]byte(done.data), &donePayload); err != nil {
		t.Fatalf("done data invalid JSON: %v", err)
	}
	if donePayload["totalLines"] != 300 {
		t.Errorf("totalLines = %d, want 300", donePayload["totalLines"])
	}
	if len(allLines) != 300 {
		t.Errorf("collected %d lines, want 300", len(allLines))
	}
}

// TestE2E_Flow2_LinesAreValidJSON verifies each generated line is valid JSON.
func TestE2E_Flow2_LinesAreValidJSON(t *testing.T) {
	ts := newTestServer(t)

	resp, err := http.Post(ts.URL+"/api/generate", "application/json",
		strings.NewReader(generatePayload))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	rawBody, _ := io.ReadAll(resp.Body)
	events := parseSSE(string(rawBody))

	for _, ev := range events {
		if ev.event != "batch" {
			continue
		}
		var payload map[string][]string
		json.Unmarshal([]byte(ev.data), &payload) //nolint:errcheck
		for i, line := range payload["lines"] {
			if !json.Valid([]byte(line)) {
				t.Errorf("batch line %d not valid JSON: %s", i, line)
			}
		}
	}
}

// ── Flow 3: Generate → Send to Analyzer ──────────────────────────────────────

// TestE2E_Flow3_GenerateThenAnalyze pipes generated lines directly into the
// analyzer — the round-trip the "Send to Analyzer" button triggers.
func TestE2E_Flow3_GenerateThenAnalyze(t *testing.T) {
	ts := newTestServer(t)

	// Step 1: generate.
	resp, err := http.Post(ts.URL+"/api/generate", "application/json",
		strings.NewReader(generatePayload))
	if err != nil {
		t.Fatal(err)
	}
	rawBody, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	events := parseSSE(string(rawBody))
	var sb strings.Builder
	for _, ev := range events {
		if ev.event != "batch" {
			continue
		}
		var payload map[string][]string
		json.Unmarshal([]byte(ev.data), &payload) //nolint:errcheck
		for _, line := range payload["lines"] {
			sb.WriteString(line)
			sb.WriteByte('\n')
		}
	}

	generatedLogs := sb.String()
	if generatedLogs == "" {
		t.Fatal("no generated logs to analyze")
	}

	// Step 2: analyze the generated logs.
	body, ct := multipartUpload(t, generatedLogs)
	resp2, err := http.Post(ts.URL+"/api/analyze", ct, body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp2.Body)
		t.Fatalf("analyze status %d: %s", resp2.StatusCode, b)
	}

	var result map[string]any
	if err := json.NewDecoder(resp2.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}

	summary := result["summary"].(map[string]any)
	if summary["total_lines"].(float64) != 300 {
		t.Errorf("total_lines = %v, want 300", summary["total_lines"])
	}
	if result["format_detected"] != "json" {
		t.Errorf("format_detected = %v, want json", result["format_detected"])
	}
	if patterns, ok := result["patterns"].([]any); !ok || len(patterns) == 0 {
		t.Error("patterns empty or missing after generate→analyze round-trip")
	}
}

// ── Flow 4: Memory safety ─────────────────────────────────────────────────────

// TestE2E_Flow4_MemorySafety uploads ~50MB of log data and verifies that heap
// growth is bounded (streaming should prevent proportional allocation).
func TestE2E_Flow4_MemorySafety(t *testing.T) {
	ts := newTestServer(t)

	// Build ~50MB of JSON log data. Each line ~110 bytes → ~455K lines.
	const targetBytes = 50 << 20 // 50 MB
	line := `{"timestamp":"2024-01-01T00:00:00Z","level":"info","message":"memory safety stress test line","source":"e2e"}` + "\n"
	linesNeeded := targetBytes / len(line)

	var sb strings.Builder
	sb.Grow(targetBytes)
	for i := 0; i < linesNeeded; i++ {
		sb.WriteString(line)
	}
	logData := sb.String()

	// Baseline heap after GC.
	runtime.GC()
	var msBefore runtime.MemStats
	runtime.ReadMemStats(&msBefore)

	body, ct := multipartUpload(t, logData)
	resp, err := http.Post(ts.URL+"/api/analyze", ct, body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("status %d: %s", resp.StatusCode, b)
	}

	// Verify response is valid.
	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}
	summary := result["summary"].(map[string]any)
	if summary["total_lines"].(float64) == 0 {
		t.Error("expected non-zero total_lines for large file")
	}

	// Post-request heap.
	runtime.GC()
	var msAfter runtime.MemStats
	runtime.ReadMemStats(&msAfter)

	heapGrowthMB := int64(msAfter.HeapInuse) - int64(msBefore.HeapInuse)
	heapGrowthMB /= 1 << 20

	// Streaming should keep heap growth well below the file size (50MB).
	// We allow up to 30MB headroom for MaxEntries (10k entries) + overhead.
	const limitMB = 30
	t.Logf("heap growth: %+d MB (limit: %d MB)", heapGrowthMB, limitMB)
	if heapGrowthMB > limitMB {
		t.Errorf("heap grew by %d MB on a 50MB upload; want ≤ %d MB", heapGrowthMB, limitMB)
	}

	// Report throughput info.
	t.Logf("upload size: %.1f MB, total_lines processed: %v",
		float64(len(logData))/(1<<20), summary["total_lines"])
}

// TestE2E_Flow4_EntryCapEnforced verifies that even with a large file,
// the entries array in the response is capped (not proportional to input size).
func TestE2E_Flow4_EntryCapEnforced(t *testing.T) {
	ts := newTestServer(t)

	// 15k lines — above the 10k MaxEntries cap.
	var sb strings.Builder
	for i := 0; i < 15_000; i++ {
		sb.WriteString(fmt.Sprintf(
			`{"timestamp":"2024-01-01T%02d:%02d:%02d Z","level":"info","message":"line %d"}`,
			(i/3600)%24, (i/60)%60, i%60, i,
		) + "\n")
	}

	body, ct := multipartUpload(t, sb.String())
	resp, err := http.Post(ts.URL+"/api/analyze", ct, body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("status %d: %s", resp.StatusCode, b)
	}

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result) //nolint:errcheck

	entries := result["entries"].([]any)
	if len(entries) > 10_000 {
		t.Errorf("entries not capped: got %d", len(entries))
	}

	// Summary must reflect all 15k lines.
	summary := result["summary"].(map[string]any)
	if summary["total_lines"].(float64) != 15_000 {
		t.Errorf("total_lines = %v, want 15000", summary["total_lines"])
	}
}
