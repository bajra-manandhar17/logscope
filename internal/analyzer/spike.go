package analyzer

import "math"

// DetectSpikes flags time buckets with anomalously high log counts using
// mean + standard deviation thresholds. Requires ≥3 buckets for meaningful stats.
func DetectSpikes(buckets []TimeBucket) []Spike {
	if len(buckets) < 3 {
		return nil
	}

	// Compute mean.
	sum := 0.0
	for _, b := range buckets {
		sum += float64(b.Count)
	}
	mean := sum / float64(len(buckets))

	// Compute stddev.
	variance := 0.0
	for _, b := range buckets {
		d := float64(b.Count) - mean
		variance += d * d
	}
	stddev := math.Sqrt(variance / float64(len(buckets)))

	if stddev == 0 {
		return nil // uniform distribution — no spikes
	}

	threshold2 := mean + 2*stddev
	threshold3 := mean + 3*stddev

	var spikes []Spike
	for _, b := range buckets {
		c := float64(b.Count)
		if c > threshold3 {
			spikes = append(spikes, Spike{
				BucketTimestamp: b.Timestamp,
				Count:           b.Count,
				Threshold:       threshold2,
				Severity:        "high",
			})
		} else if c > threshold2 {
			spikes = append(spikes, Spike{
				BucketTimestamp: b.Timestamp,
				Count:           b.Count,
				Threshold:       threshold2,
				Severity:        "medium",
			})
		}
	}
	return spikes
}
