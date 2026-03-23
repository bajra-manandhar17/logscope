package analyzer

import (
	"encoding/json"
	"testing"
	"time"
)

func TestConstants(t *testing.T) {
	if MaxEntries != 10_000 {
		t.Errorf("MaxEntries = %d, want 10000", MaxEntries)
	}
	if MaxPatterns != 10_000 {
		t.Errorf("MaxPatterns = %d, want 10000", MaxPatterns)
	}
}

func TestLogEntryJSONTags(t *testing.T) {
	e := LogEntry{
		Timestamp:  time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC),
		Level:      "error",
		Message:    "something failed",
		Source:     "app.go",
		Raw:        "2024-01-02 03:04:05 ERROR something failed",
		LineNumber: 42,
	}
	b, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	json.Unmarshal(b, &m)
	for _, key := range []string{"timestamp", "level", "message", "source", "raw", "line_number"} {
		if _, ok := m[key]; !ok {
			t.Errorf("LogEntry missing JSON key %q", key)
		}
	}
}

func TestSummaryJSONTags(t *testing.T) {
	s := Summary{
		TotalLines:  100,
		ErrorCount:  5,
		WarnCount:   10,
		InfoCount:   80,
		DebugCount:  5,
		TimeRange:   [2]time.Time{time.Now(), time.Now()},
		TopSources:  []string{"app.go"},
	}
	b, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	json.Unmarshal(b, &m)
	for _, key := range []string{"total_lines", "error_count", "warn_count", "info_count", "debug_count", "time_range", "top_sources"} {
		if _, ok := m[key]; !ok {
			t.Errorf("Summary missing JSON key %q", key)
		}
	}
}

func TestPatternJSONTags(t *testing.T) {
	p := Pattern{
		Template:   "connection refused *",
		Count:      3,
		SampleLine: "connection refused 127.0.0.1",
	}
	b, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	json.Unmarshal(b, &m)
	for _, key := range []string{"template", "count", "sample_line"} {
		if _, ok := m[key]; !ok {
			t.Errorf("Pattern missing JSON key %q", key)
		}
	}
}

func TestTimeBucketJSONTags(t *testing.T) {
	tb := TimeBucket{
		Timestamp:  time.Now(),
		Count:      50,
		ErrorCount: 2,
	}
	b, err := json.Marshal(tb)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	json.Unmarshal(b, &m)
	for _, key := range []string{"timestamp", "count", "error_count"} {
		if _, ok := m[key]; !ok {
			t.Errorf("TimeBucket missing JSON key %q", key)
		}
	}
}

func TestAnalysisResultJSONTags(t *testing.T) {
	ar := AnalysisResult{
		FormatDetected: "json",
		Summary:        Summary{},
		Entries:        []LogEntry{},
		Patterns:       []Pattern{},
		TimeSeries:     []TimeBucket{},
		BucketInterval: "1m",
	}
	b, err := json.Marshal(ar)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	json.Unmarshal(b, &m)
	for _, key := range []string{"format_detected", "summary", "entries", "patterns", "time_series", "bucket_interval"} {
		if _, ok := m[key]; !ok {
			t.Errorf("AnalysisResult missing JSON key %q", key)
		}
	}
}
