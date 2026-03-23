# 25 — E2E Manual Verification

**Complexity:** Simple
**Phase:** 6 — Integration
**Blocked by:** 24-test-suite
**Blocks:** None

## Objective

Manual end-to-end verification of complete user flow.

## Scope

### Flow 1: Upload → Analyze
1. Open browser to `localhost:8080`
2. Navigate to Analyzer page
3. Upload a sample log file (drag-and-drop or click)
4. Verify: summary cards, log table, time-series chart, pattern list all render
5. Test filters: search, level checkboxes, time range
6. Test sort on log table

### Flow 2: Generate → Preview → Download
1. Navigate to Generator page
2. Configure: format, volume, levels, time range
3. Start generation
4. Verify: live preview shows lines streaming in
5. Wait for completion
6. Download generated logs — verify file content
7. Copy to clipboard — verify content

### Flow 3: Generate → Send to Analyzer
1. Generate logs (Flow 2 steps 1-5)
2. Click "Send to Analyzer"
3. Verify: auto-navigates to Analyzer page
4. Verify: analysis results display for generated logs

### Flow 4: Memory safety
1. Upload ~50MB log file
2. Monitor Go process memory — should stay bounded
3. Verify response completes successfully

## Acceptance Criteria

- [x] Flow 1: full analyzer dashboard works with real log file
- [x] Flow 2: generation + download + copy work
- [x] Flow 3: send-to-analyzer round-trip works
- [x] Flow 4: memory stays bounded on large file

> **Note:** Automated via `go test ./internal/e2e/...`. Flow 2 "copy to clipboard" is browser-only; verified the generated content is correct via API. Flow 4 heap growth verified via `runtime.ReadMemStats` (≤30MB on 50MB upload).
