// Mirrors Go structs in internal/analyzer/types.go and internal/generator/config.go

export interface LogEntry {
  timestamp: string
  level: string
  message: string
  source: string
  raw: string
  line_number: number
}

export interface Summary {
  total_lines: number
  error_count: number
  warn_count: number
  info_count: number
  debug_count: number
  time_range: [string, string]
  top_sources: string[]
}

export interface Pattern {
  template: string
  count: number
  sample_line: string
}

export interface TimeBucket {
  timestamp: string
  count: number
  error_count: number
}

export interface AnalysisResult {
  format_detected: string
  summary: Summary
  entries: LogEntry[]
  patterns: Pattern[]
  time_series: TimeBucket[]
  bucket_interval: string
}

export interface GenerateConfig {
  format: 'json' | 'plaintext'
  total_lines: number
  levels: Record<string, number>
  start: string
  end: string
}

export interface ApiError {
  error: string
  code: string
}
