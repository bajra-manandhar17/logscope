# 14 — Router, Layout & Navigation

**Complexity:** Simple
**Phase:** 4 — Frontend Analyzer
**Blocked by:** 01-project-scaffolding, 02-server-entrypoint
**Blocks:** 15-file-upload-store-api, 16-summary-cards, 17-log-table-filter, 18-time-series-chart, 19-pattern-list, 20-config-form-store-api, 21-live-preview, 22-action-bar

## Objective

Set up React Router, shared layout shell, and navigation between Analyzer and Generator pages.

## Scope

- `src/App.tsx` — React Router: `/` redirects to `/analyze`, routes for `/analyze` and `/generate`
- `src/components/shared/Layout.tsx` — page shell wrapping children
- `src/components/shared/Navigation.tsx` — nav bar/sidebar linking Analyzer ↔ Generator
- `src/pages/AnalyzerPage.tsx` — placeholder, orchestrates analyzer components
- `src/pages/GeneratorPage.tsx` — placeholder, orchestrates generator components
- `src/types/index.ts` — shared TS types mirroring Go structs

## Acceptance Criteria

- [x] `/` redirects to `/analyze`
- [x] Navigation switches between `/analyze` and `/generate`
- [x] Layout renders with nav + content area
- [x] Types file exports all shared interfaces
- [x] Pages render placeholder content
