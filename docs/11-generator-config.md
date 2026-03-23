# 11 — Generator Config & Validation

**Complexity:** Simple
**Phase:** 3 — Backend Generator
**Blocked by:** 01-project-scaffolding
**Blocks:** 12-streaming-generator, 13-generate-handler

## Objective

Define `GenerateConfig` struct with validation logic.

## Scope

- `internal/generator/config.go`
- `GenerateConfig` struct:
  - `Format`: "json" | "plaintext"
  - `TotalLines`: int (1..1,000,000)
  - `Levels`: map with error/warn/info/debug floats (must sum to ~1.0)
  - `TimeRange`: start + end (RFC3339)
- `Validate() error` method — check all constraints
- `internal/generator/config_test.go`

## Acceptance Criteria

- [x] Valid config passes validation
- [x] Rejects totalLines outside range
- [x] Rejects levels not summing to 1.0 (with tolerance)
- [x] Rejects invalid format values
- [x] Rejects end before start in time range
- [x] Tests cover all validation paths
