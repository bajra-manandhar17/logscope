# 22 — Action Bar

**Complexity:** Moderate
**Phase:** 5 — Frontend Generator
**Blocked by:** 10-analyze-handler, 20-config-form-store-api
**Blocks:** 23-integration-build

## Objective

Download, copy, and send-to-analyzer actions for generated logs.

## Scope

- `src/components/generator/ActionBar.tsx`
- **Download:** join `generatorStore.lines` into text, trigger browser download as `.log` or `.json` file
- **Copy:** copy all lines to clipboard via `navigator.clipboard.writeText()`
- **Send to Analyzer:**
  1. Join lines into string
  2. Create `new Blob([text], { type: 'text/plain' })`
  3. Call `analyzerStore.upload(blob)`
  4. Navigate to `/analyze` (React Router)
  5. Analyzer page displays results automatically
- Buttons disabled during generation (only enabled when status = done)
- Loading state on send-to-analyzer (upload in progress)

## Acceptance Criteria

- [x] Download triggers browser file save with correct content
- [x] Copy places all lines in clipboard
- [x] Send-to-analyzer uploads blob and navigates to analyzer page
- [x] Analyzer page shows results from generated logs
- [x] Buttons disabled while generating
- [x] Loading state during send-to-analyzer upload
