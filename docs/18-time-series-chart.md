# 18 — Time-Series Chart

**Complexity:** Moderate
**Phase:** 4 — Frontend Analyzer
**Blocked by:** 15-file-upload-store-api
**Blocks:** 23-integration-build

## Objective

Recharts line/bar chart showing log volume + error rate over time buckets.

## Scope

- `src/components/analyzer/TimeSeriesChart.tsx`
- Read from `analyzerStore.result.timeSeries` + `analyzerStore.result.bucketInterval`
- Dual visualization:
  - Bar chart: total log count per bucket
  - Line overlay: error count per bucket
- X-axis: formatted timestamps (adapt label format to bucket interval)
- Y-axis: count
- Tooltip showing bucket time, total count, error count
- Responsive container
- Empty state when no time-series data

## Acceptance Criteria

- [x] Renders bar chart for log volume
- [x] Renders line overlay for error count
- [x] Tooltip displays correct data
- [x] X-axis labels adapt to bucket interval granularity
- [x] Responsive sizing
- [x] Empty state when no data
