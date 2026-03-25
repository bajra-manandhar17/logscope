package analyzer

import (
	"math"
	"testing"
)

func TestShannonEntropy_Empty(t *testing.T) {
	if e := ShannonEntropy(""); e != 0.0 {
		t.Errorf("want 0.0, got %f", e)
	}
}

func TestShannonEntropy_SingleRepeated(t *testing.T) {
	// All same byte → zero entropy.
	if e := ShannonEntropy("aaaaaaa"); e != 0.0 {
		t.Errorf("want 0.0, got %f", e)
	}
}

func TestShannonEntropy_TwoEqualChars(t *testing.T) {
	// "ab" → each has p=0.5 → H = -2*(0.5*log2(0.5)) = 1.0
	e := ShannonEntropy("ab")
	if math.Abs(e-1.0) > 0.001 {
		t.Errorf("want ~1.0, got %f", e)
	}
}

func TestShannonEntropy_HighEntropy(t *testing.T) {
	// 256 distinct bytes → max entropy = 8.0
	var buf [256]byte
	for i := range buf {
		buf[i] = byte(i)
	}
	e := ShannonEntropy(string(buf[:]))
	if math.Abs(e-8.0) > 0.001 {
		t.Errorf("want ~8.0, got %f", e)
	}
}

func TestShannonEntropy_StructuredMessage(t *testing.T) {
	// Typical structured log message should have moderate entropy.
	e := ShannonEntropy("INFO 2024-01-15 user logged in from 192.168.1.1")
	if e <= 0 || e > 5.0 {
		t.Errorf("expected moderate entropy, got %f", e)
	}
}

func TestEnrichEntropy_Empty(t *testing.T) {
	avg, high := EnrichEntropy(nil)
	if avg != 0.0 || high != 0 {
		t.Errorf("want (0,0), got (%f,%d)", avg, high)
	}
}

func TestEnrichEntropy_SetsFieldsInPlace(t *testing.T) {
	entries := []LogEntry{
		{Message: "hello"},
		{Message: "world"},
	}
	avg, _ := EnrichEntropy(entries)
	if avg <= 0 {
		t.Errorf("avg should be positive, got %f", avg)
	}
	for i, e := range entries {
		if e.Entropy <= 0 {
			t.Errorf("entry %d entropy not set: %f", i, e.Entropy)
		}
	}
}

func TestEnrichEntropy_HighEntropyCount(t *testing.T) {
	// Build a high-entropy message (all 256 byte values repeated).
	var buf [256]byte
	for i := range buf {
		buf[i] = byte(i)
	}
	highMsg := string(buf[:])

	entries := []LogEntry{
		{Message: "aaaa"},    // low entropy
		{Message: highMsg},   // high entropy (8.0)
		{Message: "bbbbb"},   // low entropy
	}
	_, high := EnrichEntropy(entries)
	if high != 1 {
		t.Errorf("want 1 high-entropy entry, got %d", high)
	}
}
