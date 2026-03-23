package generator

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

const batchSize = 100

// Batch is a slice of formatted log lines.
type Batch []string

var sources = []string{
	"api-gateway", "auth-service", "billing-service", "cache-layer",
	"db-proxy", "email-worker", "frontend", "job-scheduler",
	"metrics-collector", "notification-service", "order-service", "payment-gateway",
	"queue-consumer", "search-indexer", "session-manager", "user-service",
}

// messageTemplates maps level to a pool of message templates.
// %s placeholders are filled with random values from their respective pools.
var messageTemplates = map[string][]string{
	"error": {
		"database connection failed: timeout after %dms",
		"failed to process request: %s",
		"unhandled exception in %s: nil pointer dereference",
		"authentication failed for user %s: invalid credentials",
		"payment declined for order %s: insufficient funds",
		"circuit breaker open for upstream %s",
		"disk write error on volume %s: no space left",
		"TLS handshake failed with peer %s",
	},
	"warn": {
		"slow query detected: %dms (threshold: 500ms)",
		"retry attempt %d/3 for operation %s",
		"cache miss rate above threshold: %.1f%%",
		"memory usage at %d%% of limit",
		"deprecated endpoint %s called by client %s",
		"rate limit approaching for tenant %s: %d/1000 req/min",
		"session token expiring soon for user %s",
		"queue depth %d exceeds warning threshold",
	},
	"info": {
		"GET %s 200 %dms",
		"POST %s 201 %dms",
		"user %s logged in from %s",
		"job %s completed in %dms",
		"cache warmed: %d entries loaded",
		"started worker pool with %d goroutines",
		"order %s transitioned to state %s",
		"email sent to %s via %s",
	},
	"debug": {
		"entering function %s with args: %v",
		"SQL: SELECT * FROM %s WHERE id = %d",
		"cache lookup for key %s: hit=%v",
		"HTTP response headers: %v",
		"parsed %d records from upstream response",
		"lock acquired on resource %s by goroutine %d",
		"feature flag %s = %v for tenant %s",
		"deserialized payload: %d bytes",
	},
}

var (
	endpoints  = []string{"/api/v1/users", "/api/v1/orders", "/api/v1/products", "/health", "/metrics", "/api/v2/auth"}
	operations = []string{"flush-cache", "send-email", "reindex", "snapshot", "gc", "rotate-keys"}
	userIDs    = []string{"usr_4f2a", "usr_9c1b", "usr_7d3e", "usr_2a8f", "usr_6b5c"}
	orderIDs   = []string{"ord_001", "ord_042", "ord_199", "ord_378", "ord_512"}
	ipAddrs    = []string{"10.0.1.12", "192.168.3.44", "172.16.0.8", "10.2.0.99"}
)

// Generate yields batches of log lines to the returned channel until all lines
// are generated or ctx is cancelled. The error channel receives nil on success
// or ctx.Err() on cancellation.
func Generate(ctx context.Context, cfg GenerateConfig) (<-chan Batch, <-chan error) {
	batches := make(chan Batch)
	errc := make(chan error, 1)

	go func() {
		defer close(batches)
		defer close(errc)

		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		levels := buildLevelWheel(cfg.Levels)
		spanNs := cfg.End.Sub(cfg.Start).Nanoseconds()

		batch := make(Batch, 0, batchSize)
		for i := 0; i < cfg.TotalLines; i++ {
			select {
			case <-ctx.Done():
				errc <- ctx.Err()
				return
			default:
			}

			ts := cfg.Start.Add(time.Duration(rng.Int63n(spanNs)))
			level := pickLevel(rng, levels)
			source := sources[rng.Intn(len(sources))]
			msg := renderMessage(rng, level)

			var line string
			if cfg.Format == "json" {
				line = formatJSON(ts, level, source, msg)
			} else {
				line = formatPlaintext(ts, level, source, msg)
			}

			batch = append(batch, line)
			if len(batch) >= batchSize {
				select {
				case <-ctx.Done():
					errc <- ctx.Err()
					return
				case batches <- batch:
				}
				batch = make(Batch, 0, batchSize)
			}
		}

		if len(batch) > 0 {
			select {
			case <-ctx.Done():
				errc <- ctx.Err()
				return
			case batches <- batch:
			}
		}

		errc <- nil
	}()

	return batches, errc
}

// levelEntry is one entry in the cumulative wheel used for weighted random selection.
type levelEntry struct {
	level     string
	threshold float64
}

func buildLevelWheel(weights map[string]float64) []levelEntry {
	wheel := make([]levelEntry, 0, len(weights))
	var cum float64
	for _, l := range []string{"error", "warn", "info", "debug"} {
		w, ok := weights[l]
		if !ok || w == 0 {
			continue
		}
		cum += w
		wheel = append(wheel, levelEntry{level: l, threshold: cum})
	}
	return wheel
}

func pickLevel(rng *rand.Rand, wheel []levelEntry) string {
	r := rng.Float64()
	for _, e := range wheel {
		if r < e.threshold {
			return e.level
		}
	}
	return wheel[len(wheel)-1].level
}

func renderMessage(rng *rand.Rand, level string) string {
	tmpls := messageTemplates[level]
	tmpl := tmpls[rng.Intn(len(tmpls))]
	// Fill format verbs with plausible random values.
	return fillTemplate(rng, tmpl)
}

func fillTemplate(rng *rand.Rand, tmpl string) string {
	// Replace each verb type with a context-appropriate value.
	// We use fmt.Sprintf after converting verbs to concrete values sequentially.
	// To avoid mismatches, we render by scanning for verbs manually.
	out := make([]byte, 0, len(tmpl)*2)
	i := 0
	for i < len(tmpl) {
		if tmpl[i] != '%' || i+1 >= len(tmpl) {
			out = append(out, tmpl[i])
			i++
			continue
		}
		verb := tmpl[i+1]
		var replacement string
		switch verb {
		case 'd':
			replacement = fmt.Sprintf("%d", rng.Intn(10000)+1)
		case 'f':
			replacement = fmt.Sprintf("%.1f", rng.Float64()*100)
		case 's':
			replacement = randomStringValue(rng)
		case 'v':
			replacement = randomVValue(rng)
		default:
			out = append(out, tmpl[i], tmpl[i+1])
			i += 2
			continue
		}
		out = append(out, []byte(replacement)...)
		i += 2
	}
	return string(out)
}

func randomStringValue(rng *rand.Rand) string {
	pools := [][]string{endpoints, operations, userIDs, orderIDs, ipAddrs, sources}
	pool := pools[rng.Intn(len(pools))]
	return pool[rng.Intn(len(pool))]
}

func randomVValue(rng *rand.Rand) string {
	switch rng.Intn(3) {
	case 0:
		return fmt.Sprintf("%v", rng.Intn(2) == 1)
	case 1:
		return fmt.Sprintf("[%s %s]", randomStringValue(rng), randomStringValue(rng))
	default:
		return fmt.Sprintf("%d", rng.Intn(1000))
	}
}

func formatJSON(ts time.Time, level, source, msg string) string {
	m := map[string]string{
		"timestamp": ts.UTC().Format(time.RFC3339),
		"level":     level,
		"source":    source,
		"message":   msg,
	}
	b, _ := json.Marshal(m)
	return string(b)
}

func formatPlaintext(ts time.Time, level, source, msg string) string {
	return fmt.Sprintf("%s %s [%s] %s",
		ts.UTC().Format(time.RFC3339),
		levelUpper(level),
		source,
		msg,
	)
}

func levelUpper(l string) string {
	switch l {
	case "error":
		return "ERROR"
	case "warn":
		return "WARN"
	case "info":
		return "INFO"
	case "debug":
		return "DEBUG"
	}
	return l
}
