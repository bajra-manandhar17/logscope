package analyzer

import "math"

const highEntropyThreshold = 4.0

// ShannonEntropy computes the byte-level Shannon entropy of a string.
// Returns 0.0 for empty strings. Max theoretical value ~8.0 for uniformly
// distributed bytes.
func ShannonEntropy(msg string) float64 {
	n := len(msg)
	if n == 0 {
		return 0.0
	}

	var freq [256]int
	for i := 0; i < n; i++ {
		freq[msg[i]]++
	}

	entropy := 0.0
	fn := float64(n)
	for _, f := range freq {
		if f == 0 {
			continue
		}
		p := float64(f) / fn
		entropy -= p * math.Log2(p)
	}
	return entropy
}

// EnrichEntropy sets Entropy on each entry in-place and returns aggregate stats.
func EnrichEntropy(entries []LogEntry) (avgEntropy float64, highCount int) {
	if len(entries) == 0 {
		return 0.0, 0
	}

	sum := 0.0
	for i := range entries {
		e := ShannonEntropy(entries[i].Message)
		entries[i].Entropy = e
		sum += e
		if e > highEntropyThreshold {
			highCount++
		}
	}
	avgEntropy = sum / float64(len(entries))
	return avgEntropy, highCount
}
