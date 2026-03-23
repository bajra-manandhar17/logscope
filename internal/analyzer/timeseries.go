package analyzer

import (
	"sort"
	"time"
)

const (
	intervalMinute = "1m"
	intervalQuarter = "15m"
	intervalHour   = "1h"
	intervalDay    = "1d"
)

// BucketTimeSeries groups entries into time buckets, auto-selecting the bucket
// size from the overall time range. Returns buckets sorted chronologically and
// the chosen interval string. Zero-timestamp entries are skipped.
func BucketTimeSeries(entries []LogEntry) ([]TimeBucket, string) {
	// Determine time range from valid entries.
	var minTS, maxTS time.Time
	for _, e := range entries {
		if e.Timestamp.IsZero() {
			continue
		}
		if minTS.IsZero() || e.Timestamp.Before(minTS) {
			minTS = e.Timestamp
		}
		if e.Timestamp.After(maxTS) {
			maxTS = e.Timestamp
		}
	}

	interval, bucketSize := selectInterval(minTS, maxTS)

	buckets := make(map[time.Time]*TimeBucket)
	for _, e := range entries {
		if e.Timestamp.IsZero() {
			continue
		}
		start := bucketStart(e.Timestamp, bucketSize)
		b := buckets[start]
		if b == nil {
			b = &TimeBucket{Timestamp: start}
			buckets[start] = b
		}
		b.Count++
		if e.Level == "error" {
			b.ErrorCount++
		}
	}

	result := make([]TimeBucket, 0, len(buckets))
	for _, b := range buckets {
		result = append(result, *b)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})
	return result, interval
}

func selectInterval(minTS, maxTS time.Time) (string, time.Duration) {
	if minTS.IsZero() {
		return intervalMinute, time.Minute
	}
	span := maxTS.Sub(minTS)
	switch {
	case span < time.Hour:
		return intervalMinute, time.Minute
	case span < 24*time.Hour:
		return intervalQuarter, 15 * time.Minute
	case span < 7*24*time.Hour:
		return intervalHour, time.Hour
	default:
		return intervalDay, 24 * time.Hour
	}
}

// bucketStart truncates ts to the nearest bucketSize boundary in UTC.
func bucketStart(ts time.Time, bucketSize time.Duration) time.Time {
	return ts.UTC().Truncate(bucketSize)
}
