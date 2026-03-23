package generator_test

import (
	"testing"
	"time"

	"github.com/bajra-manandhar17/logscope-v2/internal/generator"
)

func validConfig() generator.GenerateConfig {
	return generator.GenerateConfig{
		Format:     "json",
		TotalLines: 100,
		Levels: map[string]float64{
			"error": 0.1,
			"warn":  0.2,
			"info":  0.6,
			"debug": 0.1,
		},
		Start: time.Now().Add(-time.Hour),
		End:   time.Now(),
	}
}

func TestValidConfig(t *testing.T) {
	if err := validConfig().Validate(); err != nil {
		t.Fatalf("expected valid config, got: %v", err)
	}
}

func TestRejectsTotalLinesZero(t *testing.T) {
	c := validConfig()
	c.TotalLines = 0
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for TotalLines=0")
	}
}

func TestRejectsTotalLinesTooLarge(t *testing.T) {
	c := validConfig()
	c.TotalLines = 1_000_001
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for TotalLines > 1,000,000")
	}
}

func TestAcceptsTotalLinesBoundary(t *testing.T) {
	c := validConfig()
	c.TotalLines = 1
	if err := c.Validate(); err != nil {
		t.Fatalf("expected valid for TotalLines=1: %v", err)
	}
	c.TotalLines = 1_000_000
	if err := c.Validate(); err != nil {
		t.Fatalf("expected valid for TotalLines=1,000,000: %v", err)
	}
}

func TestRejectsLevelsNotSummingToOne(t *testing.T) {
	c := validConfig()
	c.Levels = map[string]float64{
		"error": 0.5,
		"warn":  0.5,
		"info":  0.5,
		"debug": 0.5,
	}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for levels sum != 1.0")
	}
}

func TestAcceptsLevelsWithinTolerance(t *testing.T) {
	c := validConfig()
	// sum = 1.001, within typical float tolerance
	c.Levels = map[string]float64{
		"error": 0.1001,
		"warn":  0.2,
		"info":  0.6,
		"debug": 0.1,
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("expected valid within tolerance: %v", err)
	}
}

func TestRejectsInvalidFormat(t *testing.T) {
	c := validConfig()
	c.Format = "csv"
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestRejectsEndBeforeStart(t *testing.T) {
	c := validConfig()
	c.Start, c.End = c.End, c.Start // swap
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for end before start")
	}
}

func TestRejectsEqualStartEnd(t *testing.T) {
	c := validConfig()
	c.End = c.Start
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for start == end")
	}
}
