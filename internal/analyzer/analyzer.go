package analyzer

import (
	"context"
	"fmt"
	"io"
)

// Analyze detects the log format, streams and parses entries, then aggregates
// summary, patterns, and time-series into a single AnalysisResult.
// formatHint skips auto-detection when non-empty ("json" or "plaintext").
func Analyze(ctx context.Context, r io.Reader, formatHint string) (*AnalysisResult, error) {
	// 1. Detect format (or use hint).
	format := formatHint
	if format == "" {
		var replay io.Reader
		var err error
		format, replay, err = DetectFormat(r)
		if err != nil {
			return nil, fmt.Errorf("detect format: %w", err)
		}
		r = replay
	}

	// 2. Stream-parse entries.
	entries, errc := Parse(ctx, r, format)

	summary := NewSummaryAggregator()
	patterns := NewPatternAggregator()
	var capped []LogEntry

	for e := range entries {
		summary.Add(e)
		if e.Message != "" {
			patterns.Add(e)
		}
		if len(capped) < MaxEntries {
			capped = append(capped, e)
		}
	}

	if err := <-errc; err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	// 3. Compute time-series from capped entries (good enough per design doc).
	timeSeries, bucketInterval := BucketTimeSeries(capped)
	summaryResult := summary.Result()

	// 4. Intelligence post-processing.
	avgEntropy, highEntropyCount := EnrichEntropy(capped)
	spikes := DetectSpikes(timeSeries)

	// Silence detection needs bucket duration.
	_, bucketSize := selectInterval(summaryResult.TimeRange[0], summaryResult.TimeRange[1])
	silenceGaps := DetectSilenceGaps(capped, timeSeries, summaryResult.TopSources, bucketSize)

	// Causal sequence detection.
	patternResult := patterns.Result(MaxPatterns)
	causalSeqs := DetectCausalSequences(capped, patternResult)

	return &AnalysisResult{
		FormatDetected: format,
		Summary:        summaryResult,
		Entries:        capped,
		Patterns:       patternResult,
		TimeSeries:     timeSeries,
		BucketInterval: bucketInterval,
		Intelligence: Intelligence{
			Spikes:           spikes,
			SilenceGaps:      silenceGaps,
			CausalSequences:  causalSeqs,
			AvgEntropy:       avgEntropy,
			HighEntropyCount: highEntropyCount,
		},
	}, nil
}
