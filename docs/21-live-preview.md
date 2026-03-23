# 21 — Live Preview

**Complexity:** Moderate
**Phase:** 5 — Frontend Generator
**Blocked by:** 20-config-form-store-api
**Blocks:** 23-integration-build

## Objective

Scrollable log preview that updates in real-time during SSE generation.

## Scope

- `src/components/generator/LivePreview.tsx`
- Read from `generatorStore.lines`
- Auto-scroll to bottom as new lines arrive (with option to pause auto-scroll)
- Virtualized if line count is large (reuse `@tanstack/react-virtual`)
- Show generation progress (lines generated / total requested)
- Syntax-highlighted or formatted log lines (basic: level color coding)
- Empty state before generation starts

## Performance

- Must handle up to 1M accumulated lines without freezing
- Virtualization critical at high line counts
- Batch DOM updates (React batching should handle this)

## Acceptance Criteria

- [x] Lines appear in real-time as batches arrive
- [x] Auto-scrolls to newest lines
- [x] User can pause/resume auto-scroll
- [x] Progress indicator shows lines generated vs total
- [x] Handles large line counts without UI lag
- [x] Empty state before generation
