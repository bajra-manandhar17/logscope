package analyzer

import (
	"fmt"
	"sort"
	"time"
)

const maxSilenceGaps = 50

// DetectSilenceGaps finds periods where a source has no log entries while other
// sources remain active. Requires ≥2 consecutive empty buckets to qualify as a gap.
func DetectSilenceGaps(entries []LogEntry, buckets []TimeBucket, topSources []string, bucketSize time.Duration) []SilenceGap {
	if len(buckets) < 2 || len(topSources) < 2 {
		return nil
	}

	// Build sorted bucket timestamps.
	bucketTimes := make([]time.Time, len(buckets))
	for i, b := range buckets {
		bucketTimes[i] = b.Timestamp
	}
	sort.Slice(bucketTimes, func(i, j int) bool {
		return bucketTimes[i].Before(bucketTimes[j])
	})

	// Build per-bucket-per-source presence and total active source count.
	sourceSet := make(map[string]bool, len(topSources))
	for _, s := range topSources {
		sourceSet[s] = true
	}

	type bucketKey = time.Time
	// presence[bucket][source] = true if source has entries in that bucket.
	presence := make(map[bucketKey]map[string]bool, len(bucketTimes))
	// activeCount[bucket] = number of distinct sources active in bucket.
	activeCount := make(map[bucketKey]int, len(bucketTimes))

	for _, e := range entries {
		if e.Timestamp.IsZero() || e.Source == "" || !sourceSet[e.Source] {
			continue
		}
		bk := bucketStart(e.Timestamp, bucketSize)
		if presence[bk] == nil {
			presence[bk] = make(map[string]bool)
		}
		if !presence[bk][e.Source] {
			presence[bk][e.Source] = true
			activeCount[bk]++
		}
	}

	var gaps []SilenceGap

	for _, src := range topSources {
		var gapStart time.Time
		gapLen := 0

		for _, bt := range bucketTimes {
			srcActive := presence[bt] != nil && presence[bt][src]

			if !srcActive && activeCount[bt] > 0 {
				// Source is silent but others are active.
				if gapLen == 0 {
					gapStart = bt
				}
				gapLen++
			} else {
				if gapLen >= 2 {
					gapEnd := bucketTimes[0] // placeholder
					// Find the last silent bucket.
					for _, bt2 := range bucketTimes {
						if bt2.After(gapStart) || bt2.Equal(gapStart) {
							if bt2.Before(bt) {
								gapEnd = bt2
							}
						}
					}
					gapEnd = gapEnd.Add(bucketSize) // end of last silent bucket
					dur := gapEnd.Sub(gapStart)
					avgActive := avgActiveSourcesDuring(bucketTimes, gapStart, bt, activeCount)
					gaps = append(gaps, SilenceGap{
						Source:                 src,
						GapStart:               gapStart,
						GapEnd:                 gapEnd,
						Duration:               formatDuration(dur),
						ActiveSourcesDuringGap: avgActive,
					})
				}
				gapLen = 0
			}
		}

		// Handle gap at end of time range.
		if gapLen >= 2 {
			lastBucket := bucketTimes[len(bucketTimes)-1]
			gapEnd := lastBucket.Add(bucketSize)
			dur := gapEnd.Sub(gapStart)
			avgActive := avgActiveSourcesDuring(bucketTimes, gapStart, gapEnd, activeCount)
			gaps = append(gaps, SilenceGap{
				Source:                 src,
				GapStart:               gapStart,
				GapEnd:                 gapEnd,
				Duration:               formatDuration(dur),
				ActiveSourcesDuringGap: avgActive,
			})
		}
	}

	sort.Slice(gaps, func(i, j int) bool {
		return gaps[i].GapEnd.Sub(gaps[i].GapStart) > gaps[j].GapEnd.Sub(gaps[j].GapStart)
	})

	if len(gaps) > maxSilenceGaps {
		gaps = gaps[:maxSilenceGaps]
	}
	return gaps
}

func avgActiveSourcesDuring(bucketTimes []time.Time, start, end time.Time, activeCount map[time.Time]int) int {
	total, count := 0, 0
	for _, bt := range bucketTimes {
		if (bt.Equal(start) || bt.After(start)) && bt.Before(end) {
			total += activeCount[bt]
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return total / count
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}
