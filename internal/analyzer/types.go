package analyzer

import "time"

const (
	MaxEntries  = 10_000
	MaxPatterns = 10_000
)

type LogEntry struct {
	Timestamp  time.Time `json:"timestamp"`
	Level      string    `json:"level"`
	Message    string    `json:"message"`
	Source     string    `json:"source"`
	Raw        string    `json:"raw"`
	LineNumber int       `json:"line_number"`
	Entropy    float64   `json:"entropy"`
}

type Summary struct {
	TotalLines int          `json:"total_lines"`
	ErrorCount int          `json:"error_count"`
	WarnCount  int          `json:"warn_count"`
	InfoCount  int          `json:"info_count"`
	DebugCount int          `json:"debug_count"`
	TimeRange  [2]time.Time `json:"time_range"`
	TopSources []string     `json:"top_sources"`
}

type Pattern struct {
	Template   string    `json:"template"`
	Count      int       `json:"count"`
	SampleLine string    `json:"sample_line"`
	FirstSeen  time.Time `json:"first_seen"`
	LastSeen   time.Time `json:"last_seen"`
}

type TimeBucket struct {
	Timestamp  time.Time `json:"timestamp"`
	Count      int       `json:"count"`
	ErrorCount int       `json:"error_count"`
}

type Spike struct {
	BucketTimestamp time.Time `json:"bucket_timestamp"`
	Count           int       `json:"count"`
	Threshold       float64   `json:"threshold"`
	Severity        string    `json:"severity"` // "high" (>3σ) or "medium" (>2σ)
}

type SilenceGap struct {
	Source                 string    `json:"source"`
	GapStart               time.Time `json:"gap_start"`
	GapEnd                 time.Time `json:"gap_end"`
	Duration               string    `json:"duration"`
	ActiveSourcesDuringGap int       `json:"active_sources_during_gap"`
}

type CausalSequence struct {
	PatternA      string  `json:"pattern_a"`
	PatternB      string  `json:"pattern_b"`
	Count         int     `json:"count"`
	AvgLagSeconds float64 `json:"avg_lag_seconds"`
}

type Intelligence struct {
	Spikes           []Spike          `json:"spikes"`
	SilenceGaps      []SilenceGap     `json:"silence_gaps"`
	CausalSequences  []CausalSequence `json:"causal_sequences"`
	AvgEntropy       float64          `json:"avg_entropy"`
	HighEntropyCount int              `json:"high_entropy_count"`
}

type AnalysisResult struct {
	FormatDetected string       `json:"format_detected"`
	Summary        Summary      `json:"summary"`
	Entries        []LogEntry   `json:"entries"`
	Patterns       []Pattern    `json:"patterns"`
	TimeSeries     []TimeBucket `json:"time_series"`
	BucketInterval string       `json:"bucket_interval"`
	Intelligence   Intelligence `json:"intelligence"`
}
