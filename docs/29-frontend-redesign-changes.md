## Objective

Retheme frontend from generic shadcn to Datadog/Kibana observability aesthetic. Keep Radix primitives, replace design tokens.

## Decisions
- [x] Keep Radix primitives, retheme design tokens
- [x] Dark + light mode, dark default, toggle in nav (localStorage)
- [x] Geist (UI) + JetBrains Mono (data), drop Syne
- [x] Neutral gray chrome, color only for data signals
- [x] Full-width dense layout
- [x] Restyle charts to match
- [x] All pages/components in scope, visual-only
- [x] `@fontsource-variable/jetbrains-mono` dependency approved

## Acceptance Criteria
- [x] Dark mode loads by default
- [x] Sun/moon toggle works + persists across reload
- [x] Light mode renders correctly
- [x] Full-width layout, no max-w-7xl, dense spacing (gap-2/3)
- [x] No Syne font loading (Network tab)
- [x] JetBrains Mono on: log messages, timestamps, stats, patterns
- [x] Cards: rounded-md corners, reduced padding
- [x] Nav: compact (h-10), solid bg, no blur
- [x] Charts: flat bar fills, dotted grid, mono axes
- [x] Color only on data signals (levels, charts, errors), not chrome
- [x] LogTable fills viewport height
- [x] Radix a11y intact (focus, keyboard, ARIA)
- [x] `npm run build` passes
- [x] `npx vitest run` passes

## Steps

### 0. Dependency
Install `@fontsource-variable/jetbrains-mono`

### 1. Theme — `frontend/src/index.css`
- Add `@import "@fontsource-variable/jetbrains-mono"`
- Replace `--font-heading` with `--font-mono: 'JetBrains Mono Variable', monospace`
- `--radius: 0.5rem` → `0.25rem`
- Remove all `--sidebar-*` tokens and `--color-sidebar-*` mappings
- Rewrite `:root` light palette — neutral gray chrome:
  - `--background: oklch(0.97 0.003 250)`, `--foreground: oklch(0.13 0.005 250)`
  - `--card: oklch(1.0 0 0)`, `--primary: oklch(0.55 0.05 250)` (subtle blue, interactive only)
  - `--muted: oklch(0.94 0.003 250)`, `--muted-foreground: oklch(0.45 0.01 250)`
  - `--destructive: oklch(0.55 0.22 27)`, `--border: oklch(0.88 0.005 250)`
- Rewrite `.dark` palette:
  - `--background: oklch(0.11 0.005 250)`, `--foreground: oklch(0.90 0.003 250)`
  - `--card: oklch(0.15 0.006 250)`, `--primary: oklch(0.60 0.08 220)` (muted teal-blue)
  - `--muted: oklch(0.18 0.006 250)`, `--muted-foreground: oklch(0.55 0.01 250)`
  - `--border: oklch(0.22 0.008 250)`
- Chart colors (both modes): blue/red/green/amber/purple

### 2. Font cleanup — `frontend/index.html`
- Remove 3 Google Fonts `<link>` tags for Syne

### 3. Theme toggle — `frontend/src/components/shared/Navigation.tsx`
- Add Sun/Moon toggle (lucide-react, already available)
- localStorage persistence, default dark
- Remove heading font from logo
- Nav: `h-14` → `h-10`, solid bg, no blur

### 4. UI primitives — `frontend/src/components/ui/card.tsx`
- `rounded-xl` → `rounded-md`, reduce padding, `ring-1` → `border border-border`

### 5. UI primitives — other
- **button.tsx**: `rounded-lg` → `rounded-md`
- **badge.tsx**: `rounded-4xl` → `rounded`
- **input.tsx**: `rounded-lg` → `rounded-md`
- **tabs.tsx**: `rounded-lg` → `rounded-md`
- **checkbox.tsx**: `rounded-[4px]` → `rounded-sm`

### 6. Layout density — pages
- AnalyzerPage: full-width `px-4 py-3`, `gap-3`, remove heading font, smaller headings
- GeneratorPage: same treatment
- Layout.tsx: remove `selection:bg-primary/20`

### 7. Analyzer components
- SummaryCards: `text-2xl font-mono`, `gap-2`, no hover color
- FilterBar: `rounded-md`, no blur, `font-mono` time inputs
- LogTable: `rounded-md`, viewport height, `font-mono` on data cols
- TimeSeriesChart: flat fills, subtler grid, mono axes
- IntelligencePanel: `rounded-md`, `font-mono` stats
- FileUpload: `rounded-md`, `p-6`
- ReplayControlBar: `rounded-md`
- PatternList: `rounded-md`, neutral border
- SpikeList/SilenceGapList/CausalSequenceList: `font-mono`, `rounded-md`

### 8. Generator components
- ConfigForm: `font-mono` on number inputs

## Files Modified
1. `frontend/src/index.css`
2. `frontend/index.html`
3. `frontend/src/components/shared/Navigation.tsx`
4. `frontend/src/components/shared/Layout.tsx`
5. `frontend/src/components/ui/card.tsx`
6. `frontend/src/components/ui/button.tsx`
7. `frontend/src/components/ui/badge.tsx`
8. `frontend/src/components/ui/input.tsx`
9. `frontend/src/components/ui/tabs.tsx`
10. `frontend/src/components/ui/checkbox.tsx`
11. `frontend/src/pages/AnalyzerPage.tsx`
12. `frontend/src/pages/GeneratorPage.tsx`
13. `frontend/src/components/analyzer/SummaryCards.tsx`
14. `frontend/src/components/analyzer/FilterBar.tsx`
15. `frontend/src/components/analyzer/LogTable.tsx`
16. `frontend/src/components/analyzer/TimeSeriesChart.tsx`
17. `frontend/src/components/analyzer/IntelligencePanel.tsx`
18. `frontend/src/components/analyzer/FileUpload.tsx`
19. `frontend/src/components/analyzer/ReplayControlBar.tsx`
20. `frontend/src/components/analyzer/PatternList.tsx`
21. `frontend/src/components/analyzer/SpikeList.tsx`
22. `frontend/src/components/analyzer/SilenceGapList.tsx`
23. `frontend/src/components/analyzer/CausalSequenceList.tsx`
24. `frontend/src/components/generator/ConfigForm.tsx`
