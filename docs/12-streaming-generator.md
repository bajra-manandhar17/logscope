# 12 — Streaming Log Generator

**Complexity:** Complex
**Phase:** 3 — Backend Generator
**Blocked by:** 11-generator-config
**Blocks:** 13-generate-handler

## Objective

Generate synthetic log lines in batches, respecting config for format, volume, level distribution, and time range.

## Scope

- `internal/generator/generator.go`
- Accept `GenerateConfig` + `context.Context`
- Generate lines with:
  - Timestamps distributed across configured time range
  - Log levels matching configured distribution (weighted random)
  - Realistic message content (pool of templates with variable parts)
  - Source/service names from a pool
- Yield batches of ~100 lines via callback or channel
- Respect context cancellation (stop generating on client disconnect)
- `internal/generator/generator_test.go`

## Design Notes

- Use `math/rand` seeded per-request for reproducibility (optional)
- Message templates: mix of realistic patterns (HTTP requests, DB queries, auth events, etc.)
- JSON format: structured JSON lines; plaintext format: traditional syslog-style

## Acceptance Criteria

- [x] Generates correct number of total lines
- [x] Level distribution approximately matches config (within statistical tolerance)
- [x] Timestamps fall within configured time range
- [x] Yields batches of ~100 lines
- [x] Respects context cancellation mid-generation
- [x] Both JSON and plaintext formats produce valid output
- [x] Generated JSON logs are parseable by analyzer
- [x] Tests verify line count, format, level distribution
