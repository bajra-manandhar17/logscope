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
	Template   string `json:"template"`
	Count      int    `json:"count"`
	SampleLine string `json:"sample_line"`
}

type TimeBucket struct {
	Timestamp  time.Time `json:"timestamp"`
	Count      int       `json:"count"`
	ErrorCount int       `json:"error_count"`
}

type AnalysisResult struct {
	FormatDetected string       `json:"format_detected"`
	Summary        Summary      `json:"summary"`
	Entries        []LogEntry   `json:"entries"`
	Patterns       []Pattern    `json:"patterns"`
	TimeSeries     []TimeBucket `json:"time_series"`
	BucketInterval string       `json:"bucket_interval"`
}
