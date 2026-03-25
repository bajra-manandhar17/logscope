package analyzer

import (
	"fmt"
	"testing"
	"time"
)

// --- Masking ---

func TestMaskTokens_Numbers(t *testing.T) {
	cases := []struct{ in, want string }{
		{"retry 3 times", "retry {NUM} times"},
		{"port 8080 open", "port {NUM} open"},
		{"took 123.456ms", "took {NUM}ms"},
	}
	for _, c := range cases {
		if got := maskTokens(c.in); got != c.want {
			t.Errorf("maskTokens(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestMaskTokens_UUID(t *testing.T) {
	in := "user 550e8400-e29b-41d4-a716-446655440000 logged in"
	want := "user {UUID} logged in"
	if got := maskTokens(in); got != want {
		t.Errorf("got %q", got)
	}
}

func TestMaskTokens_IP(t *testing.T) {
	in := "connection from 192.168.1.100 rejected"
	want := "connection from {IP} rejected"
	if got := maskTokens(in); got != want {
		t.Errorf("got %q", got)
	}
}

func TestMaskTokens_Hex(t *testing.T) {
	cases := []struct{ in, want string }{
		{"addr 0x1a2b3c4d fault", "addr {HEX} fault"},
		{"hash deadbeef stored", "hash {HEX} stored"},
	}
	for _, c := range cases {
		if got := maskTokens(c.in); got != c.want {
			t.Errorf("maskTokens(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestMaskTokens_NoVariableParts(t *testing.T) {
	in := "server started"
	if got := maskTokens(in); got != in {
		t.Errorf("got %q, want unchanged", got)
	}
}

func TestMaskTokens_AllVariableParts(t *testing.T) {
	in := "550e8400-e29b-41d4-a716-446655440000 192.168.0.1 42 0xdeadbeef"
	want := "{UUID} {IP} {NUM} {HEX}"
	if got := maskTokens(in); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// --- Aggregator ---

func entry(msg string) LogEntry {
	return LogEntry{Message: msg}
}

func entryAt(msg string, ts time.Time) LogEntry {
	return LogEntry{Message: msg, Timestamp: ts}
}

func TestPatternAggregator_GroupsAndCounts(t *testing.T) {
	a := NewPatternAggregator()
	a.Add(entry("retry 3 times"))
	a.Add(entry("retry 7 times"))
	a.Add(entry("retry 99 times"))

	patterns := a.Result(10)
	if len(patterns) != 1 {
		t.Fatalf("want 1 pattern, got %d", len(patterns))
	}
	if patterns[0].Template != "retry {NUM} times" {
		t.Errorf("template: got %q", patterns[0].Template)
	}
	if patterns[0].Count != 3 {
		t.Errorf("count: got %d", patterns[0].Count)
	}
	if patterns[0].SampleLine == "" {
		t.Error("sample line empty")
	}
}

func TestPatternAggregator_SampleLineRetained(t *testing.T) {
	a := NewPatternAggregator()
	a.Add(entry("retry 3 times"))
	a.Add(entry("retry 7 times"))

	p := a.Result(10)[0]
	if p.SampleLine != "retry 3 times" && p.SampleLine != "retry 7 times" {
		t.Errorf("unexpected sample: %q", p.SampleLine)
	}
}

func TestPatternAggregator_SortedByCountDesc(t *testing.T) {
	a := NewPatternAggregator()
	for i := 0; i < 5; i++ {
		a.Add(entry(fmt.Sprintf("connect from 10.0.0.%d", i)))
	}
	for i := 0; i < 2; i++ {
		a.Add(entry(fmt.Sprintf("disconnect user %d", i)))
	}

	patterns := a.Result(10)
	if len(patterns) < 2 {
		t.Fatalf("want >=2 patterns, got %d", len(patterns))
	}
	if patterns[0].Count < patterns[1].Count {
		t.Errorf("not sorted: %d < %d", patterns[0].Count, patterns[1].Count)
	}
}

func TestPatternAggregator_TopNCapped(t *testing.T) {
	a := NewPatternAggregator()
	for i := 0; i < 20; i++ {
		a.Add(entry(fmt.Sprintf("event type-%d happened", i)))
	}
	patterns := a.Result(5)
	if len(patterns) > 5 {
		t.Errorf("want <=5, got %d", len(patterns))
	}
}

func TestPatternAggregator_TemplateCap(t *testing.T) {
	a := NewPatternAggregator()
	// Insert MaxPatterns distinct templates.
	for i := 0; i < MaxPatterns; i++ {
		a.Add(entry(fmt.Sprintf("unique message alpha-%d", i)))
	}
	// One more distinct message — should be silently dropped.
	a.Add(entry("brand new unique message zzzz"))

	// Existing template should still increment.
	a.Add(entry("unique message alpha-0")) // matches template "unique message alpha-{NUM}"

	patterns := a.Result(MaxPatterns + 1)
	if len(patterns) > MaxPatterns {
		t.Errorf("template count exceeded cap: %d", len(patterns))
	}
}

func TestPatternAggregator_Empty(t *testing.T) {
	a := NewPatternAggregator()
	if patterns := a.Result(10); len(patterns) != 0 {
		t.Errorf("want 0 patterns, got %d", len(patterns))
	}
}

// --- First/Last Seen ---

func TestPatternAggregator_FirstLastSeen(t *testing.T) {
	a := NewPatternAggregator()
	t1 := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 15, 10, 5, 0, 0, time.UTC)
	t3 := time.Date(2024, 1, 15, 10, 10, 0, 0, time.UTC)

	a.Add(entryAt("retry 3 times", t2))
	a.Add(entryAt("retry 7 times", t1)) // earlier → updates firstSeen
	a.Add(entryAt("retry 99 times", t3)) // later → updates lastSeen

	p := a.Result(10)[0]
	if !p.FirstSeen.Equal(t1) {
		t.Errorf("firstSeen: want %v, got %v", t1, p.FirstSeen)
	}
	if !p.LastSeen.Equal(t3) {
		t.Errorf("lastSeen: want %v, got %v", t3, p.LastSeen)
	}
}

func TestPatternAggregator_FirstLastSeen_ZeroTimestamp(t *testing.T) {
	a := NewPatternAggregator()
	t1 := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	a.Add(entry("retry 3 times"))          // zero timestamp
	a.Add(entryAt("retry 7 times", t1))    // first non-zero

	p := a.Result(10)[0]
	if !p.FirstSeen.Equal(t1) {
		t.Errorf("firstSeen should be %v, got %v", t1, p.FirstSeen)
	}
	if !p.LastSeen.Equal(t1) {
		t.Errorf("lastSeen should be %v, got %v", t1, p.LastSeen)
	}
}
