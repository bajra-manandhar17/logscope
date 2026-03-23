package generator

import (
	"errors"
	"math"
	"time"
)

const levelSumTolerance = 0.01

// GenerateConfig holds parameters for log generation.
type GenerateConfig struct {
	Format     string             // "json" or "plaintext"
	TotalLines int                // 1..1,000,000
	Levels     map[string]float64 // keys: error/warn/info/debug; values must sum to ~1.0
	Start      time.Time
	End        time.Time
}

// Validate checks all constraints on the config.
func (c GenerateConfig) Validate() error {
	if c.Format != "json" && c.Format != "plaintext" {
		return errors.New("format must be \"json\" or \"plaintext\"")
	}
	if c.TotalLines < 1 || c.TotalLines > 1_000_000 {
		return errors.New("total_lines must be between 1 and 1,000,000")
	}
	var sum float64
	for _, v := range c.Levels {
		sum += v
	}
	if math.Abs(sum-1.0) > levelSumTolerance {
		return errors.New("level weights must sum to 1.0 (±0.01)")
	}
	if !c.End.After(c.Start) {
		return errors.New("end must be after start")
	}
	return nil
}
