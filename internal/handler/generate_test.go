package handler

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func makeGenerateBody(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return bytes.NewBuffer(b)
}

func validGeneratePayload() map[string]any {
	return map[string]any{
		"Format":     "json",
		"TotalLines": 250,
		"Levels": map[string]float64{
			"error": 0.05,
			"warn":  0.15,
			"info":  0.70,
			"debug": 0.10,
		},
		"Start": "2024-01-01T00:00:00Z",
		"End":   "2024-01-02T00:00:00Z",
	}
}

// sseEvent holds a parsed SSE event.
type sseEvent struct {
	event string
	data  string
}

func parseSSE(body string) []sseEvent {
	var events []sseEvent
	var cur sseEvent
	scanner := bufio.NewScanner(strings.NewReader(body))
	for scanner.Scan() {
		line := scanner.Text()
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

func TestGenerateHandlerSSEHeaders(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/generate", makeGenerateBody(t, validGeneratePayload()))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	GenerateHandler(rr, req)

	if got := rr.Header().Get("Content-Type"); got != "text/event-stream" {
		t.Errorf("Content-Type = %q, want text/event-stream", got)
	}
	if got := rr.Header().Get("Cache-Control"); got != "no-cache" {
		t.Errorf("Cache-Control = %q, want no-cache", got)
	}
	if got := rr.Header().Get("Connection"); got != "keep-alive" {
		t.Errorf("Connection = %q, want keep-alive", got)
	}
}

func TestGenerateHandlerBatchEvents(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/generate", makeGenerateBody(t, validGeneratePayload()))
	rr := httptest.NewRecorder()

	GenerateHandler(rr, req)

	events := parseSSE(rr.Body.String())
	if len(events) < 2 {
		t.Fatalf("expected at least 2 events (batch + done), got %d", len(events))
	}

	for _, ev := range events[:len(events)-1] {
		if ev.event != "batch" {
			t.Errorf("expected batch event, got %q", ev.event)
		}
		var payload map[string][]string
		if err := json.Unmarshal([]byte(ev.data), &payload); err != nil {
			t.Fatalf("batch data not valid JSON: %v", err)
		}
		if len(payload["lines"]) == 0 {
			t.Error("batch event has empty lines array")
		}
	}
}

func TestGenerateHandlerDoneEvent(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/generate", makeGenerateBody(t, validGeneratePayload()))
	rr := httptest.NewRecorder()

	GenerateHandler(rr, req)

	events := parseSSE(rr.Body.String())
	last := events[len(events)-1]
	if last.event != "done" {
		t.Fatalf("last event = %q, want done", last.event)
	}

	var payload map[string]int
	if err := json.Unmarshal([]byte(last.data), &payload); err != nil {
		t.Fatalf("done data not valid JSON: %v", err)
	}
	if payload["totalLines"] != 250 {
		t.Errorf("totalLines = %d, want 250", payload["totalLines"])
	}
}

func TestGenerateHandlerFlushed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/generate", makeGenerateBody(t, validGeneratePayload()))
	rr := httptest.NewRecorder()

	GenerateHandler(rr, req)

	if !rr.Flushed {
		t.Error("expected Flush to be called")
	}
}

func TestGenerateHandlerInvalidConfig(t *testing.T) {
	cases := []struct {
		name    string
		payload any
	}{
		{"bad format", map[string]any{
			"Format": "xml", "TotalLines": 10,
			"Levels": map[string]float64{"info": 1.0},
			"Start":  "2024-01-01T00:00:00Z", "End": "2024-01-02T00:00:00Z",
		}},
		{"zero lines", map[string]any{
			"Format": "json", "TotalLines": 0,
			"Levels": map[string]float64{"info": 1.0},
			"Start":  "2024-01-01T00:00:00Z", "End": "2024-01-02T00:00:00Z",
		}},
		{"bad level sum", map[string]any{
			"Format": "json", "TotalLines": 10,
			"Levels": map[string]float64{"info": 0.5},
			"Start":  "2024-01-01T00:00:00Z", "End": "2024-01-02T00:00:00Z",
		}},
		{"end before start", map[string]any{
			"Format": "json", "TotalLines": 10,
			"Levels": map[string]float64{"info": 1.0},
			"Start":  "2024-01-02T00:00:00Z", "End": "2024-01-01T00:00:00Z",
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/generate", makeGenerateBody(t, tc.payload))
			rr := httptest.NewRecorder()
			GenerateHandler(rr, req)
			if rr.Code != http.StatusBadRequest {
				t.Errorf("got %d, want 400", rr.Code)
			}
			var resp map[string]string
			json.NewDecoder(rr.Body).Decode(&resp) //nolint:errcheck
			if resp["code"] != "invalid_config" {
				t.Errorf("code = %q, want invalid_config", resp["code"])
			}
		})
	}
}

func TestGenerateHandlerBadJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/generate", strings.NewReader("not-json"))
	rr := httptest.NewRecorder()
	GenerateHandler(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("got %d, want 400", rr.Code)
	}
}
