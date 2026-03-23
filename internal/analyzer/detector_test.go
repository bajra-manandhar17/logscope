package analyzer

import (
	"io"
	"strings"
	"testing"
)

func detect(input string) (string, string, error) {
	r := strings.NewReader(input)
	format, replay, err := DetectFormat(r)
	if err != nil {
		return "", "", err
	}
	b, _ := io.ReadAll(replay)
	return format, string(b), nil
}

func TestDetectFormat_PureJSON(t *testing.T) {
	input := `{"level":"info","msg":"started"}
{"level":"error","msg":"failed"}
{"level":"warn","msg":"slow"}
`
	format, replay, err := detect(input)
	if err != nil {
		t.Fatal(err)
	}
	if format != "json" {
		t.Errorf("got %q, want \"json\"", format)
	}
	if replay != input {
		t.Errorf("replay mismatch: got %q", replay)
	}
}

func TestDetectFormat_Plaintext(t *testing.T) {
	input := "2024-01-01 ERROR something went wrong\n2024-01-01 INFO all good\n"
	format, replay, err := detect(input)
	if err != nil {
		t.Fatal(err)
	}
	if format != "plaintext" {
		t.Errorf("got %q, want \"plaintext\"", format)
	}
	if replay != input {
		t.Errorf("replay mismatch")
	}
}

func TestDetectFormat_Empty(t *testing.T) {
	format, replay, err := detect("")
	if err != nil {
		t.Fatal(err)
	}
	if format != "plaintext" {
		t.Errorf("got %q, want \"plaintext\"", format)
	}
	if replay != "" {
		t.Errorf("replay should be empty")
	}
}

func TestDetectFormat_SingleLineJSON(t *testing.T) {
	input := `{"level":"info","msg":"only line"}`
	format, _, err := detect(input)
	if err != nil {
		t.Fatal(err)
	}
	if format != "json" {
		t.Errorf("got %q, want \"json\"", format)
	}
}

func TestDetectFormat_JSONWithLeadingWhitespace(t *testing.T) {
	input := "  {\"level\":\"info\"}\n  {\"level\":\"warn\"}\n"
	format, _, err := detect(input)
	if err != nil {
		t.Fatal(err)
	}
	if format != "json" {
		t.Errorf("got %q, want \"json\"", format)
	}
}

func TestDetectFormat_MixedFallsToPlaintext(t *testing.T) {
	input := `{"level":"info","msg":"ok"}
not json at all
{"level":"warn"}
another plain line
`
	format, _, err := detect(input)
	if err != nil {
		t.Fatal(err)
	}
	if format != "plaintext" {
		t.Errorf("got %q, want \"plaintext\"", format)
	}
}

func TestDetectFormat_ReplayFullStream(t *testing.T) {
	// More than 10 lines — replay must include lines beyond the peek window
	var sb strings.Builder
	for i := 0; i < 20; i++ {
		sb.WriteString("{\"n\":1}\n")
	}
	input := sb.String()
	_, replay, err := detect(input)
	if err != nil {
		t.Fatal(err)
	}
	if replay != input {
		t.Errorf("replay truncated: got %d bytes, want %d", len(replay), len(input))
	}
}
