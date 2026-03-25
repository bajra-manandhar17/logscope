# 27 — Intelligence & Pattern Detection System

**Complexity:** Complex
**Phase:** 7 — Intelligence Layer
**Blocked by:** 09-analyzer-orchestrator, 08-time-series, 07-pattern-detector
**Blocks:** none

## Objective

Add 5 intelligence features to the analyzer pipeline: spike detection, silence detection, first/last-seen tracking, causal sequence detection, and entropy scoring. All stateless, single-file session scope.

## Scope

### Backend (Go)

**New files:**
- `internal/analyzer/entropy.go` — Shannon entropy per message + aggregate stats
- `internal/analyzer/spike.go` — mean+2σ/3σ anomaly detection on time buckets
- `internal/analyzer/silence.go` — per-source gap detection across time buckets
- `internal/analyzer/causal.go` — temporal A→B sequence detection (60s window)
- Tests for each: `*_test.go`

**Modified files:**
- `internal/analyzer/types.go` — new types: `Spike`, `SilenceGap`, `CausalSequence`, `Intelligence`. Extend `Pattern` (FirstSeen/LastSeen), `LogEntry` (Entropy), `AnalysisResult` (Intelligence)
- `internal/analyzer/pattern.go` — `Add(string)` → `Add(LogEntry)`, track first/last seen timestamps
- `internal/analyzer/pattern_test.go` — update all Add() calls
- `internal/analyzer/analyzer.go` — wire 5 post-processing steps

### Frontend (React/TS)

**New files:**
- `frontend/src/components/analyzer/IntelligencePanel.tsx` — collapsible accordion container
- `frontend/src/components/analyzer/SpikeList.tsx` — spike table
- `frontend/src/components/analyzer/SilenceGapList.tsx` — silence gap table
- `frontend/src/components/analyzer/CausalSequenceList.tsx` — A→B chain display

**Modified files:**
- `frontend/src/types/index.ts` — mirror new Go types
- `frontend/src/components/analyzer/TimeSeriesChart.tsx` — spike markers
- `frontend/src/components/analyzer/PatternList.tsx` — first/last seen
- `frontend/src/components/analyzer/LogTable.tsx` — toggleable entropy column (hidden default)
- `frontend/src/pages/AnalyzerPage.tsx` — add IntelligencePanel

## Feature Details

### 1. Entropy Scoring
- Byte-level Shannon entropy per `LogEntry.Message`
- Threshold: >4.0 = high entropy
- Summary: avg entropy + high-entropy count
- Frontend: hidden-by-default sortable column in LogTable

### 2. Spike Detection
- Mean + σ over `TimeBucket.Count` values
- >2σ = medium severity, >3σ = high severity
- Min 3 buckets required
- Frontend: ReferenceDot markers on TimeSeriesChart + list in IntelligencePanel

### 3. First/Last Seen per Pattern
- Track earliest and latest `LogEntry.Timestamp` per masked template
- Zero-timestamp entries ignored
- Frontend: timestamps displayed in PatternList

### 4. Silence Detection (Per-Source)
- Track per-source presence across time buckets
- Gap = ≥2 consecutive empty buckets while other sources active
- Cap: 50 results, sorted by duration desc
- Frontend: table in IntelligencePanel

### 5. Causal Sequence Detection
- Scope: error + warn patterns only, top 100 by frequency
- Window: 60 seconds
- Min co-occurrence: ≥3 times
- Forward scan cap: 500 events per entry
- Result cap: top 20 by count
- Reuses `maskTokens` from pattern.go
- Frontend: A→B chain display in IntelligencePanel

## Algorithms

**Spike:** `mean + 2*stddev` / `mean + 3*stddev` over bucket counts. stddev=0 → no spikes.

**Silence:** `map[bucketStart]map[source]bool` presence matrix. Walk chronologically per source, emit gap on ≥2 consecutive absent buckets where activeCount > 0.

**Causal:** Filter error/warn entries → match top-100 templates → sort by timestamp → sliding window ≤60s → accumulate `(A,B)→{count, totalLag}` → filter count≥3 → top 20.

**Entropy:** `H = -Σ p(b) * log2(p(b))` over 256 byte values.

## Memory Budget

| Feature | Memory | Notes |
|---------|--------|-------|
| Entropy | 80KB | float64 per 10K entries |
| Spikes | ~1KB | operates on existing TimeBucket slice |
| Silence | ~40KB | 1K sources × ~1K buckets boolean matrix |
| Causal | ~80KB | 100 candidates → max 10K pairs |
| First/Last seen | ~160KB | 2 timestamps per 10K patterns |

## Edge Cases

- Log file with no timestamps → spike, silence, causal return nil; entropy still works
- Uniform log rate → no spikes detected
- Single source → no silence gaps (need ≥2 sources)
- No error/warn entries → no causal sequences
- All identical messages → entropy = 0.0

## Acceptance Criteria

- [x] `ShannonEntropy` returns 0 for empty/uniform, correct values for known inputs
- [x] `EnrichEntropy` sets Entropy on entries in-place, computes avg + high count
- [x] `DetectSpikes` returns nil for uniform/insufficient data, flags anomalous buckets
- [x] `Pattern.FirstSeen`/`LastSeen` populated correctly, zero-timestamps ignored
- [x] `DetectSilenceGaps` detects per-source gaps, requires ≥2 sources, handles gap-at-end
- [x] `DetectCausalSequences` finds A→B chains, respects threshold/window/caps
- [x] All existing analyzer tests pass (no regression)
- [x] Frontend displays intelligence panel with accordion sections
- [x] Spike markers on time-series chart
- [x] Entropy column toggleable in log table
- [x] First/last seen shown on pattern list
- [x] `go test ./internal/analyzer/...` passes
- [x] `cd frontend && npx vitest run` passes
