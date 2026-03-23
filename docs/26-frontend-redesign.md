# 26 — Frontend Redesign

**Complexity:** Moderate
**Phase:** 6 — Frontend Polish
**Blocked by:** 14-router-layout-nav, 15-file-upload-store-api, 16-summary-cards, 17-log-table-filter, 18-time-series-chart, 19-pattern-list, 20-config-form-store-api, 21-live-preview, 22-action-bar
**Blocks:** None

## Objective

Redesign LogScope v2 frontend for premium dev-tool aesthetic (Linear/Vercel caliber). Cool blue palette, shadcn/ui component polish, refined spacing and hierarchy. Purely visual — no store or logic changes.

**Decisions:** Cool blue primary, sticky filter bar, polish both light + dark modes.

## Scope

- `frontend/src/index.css` — full palette rewrite
- `frontend/src/components/shared/Navigation.tsx` — restyle nav
- `frontend/src/components/shared/Layout.tsx` — minor additions
- `frontend/src/pages/AnalyzerPage.tsx` — layout + sticky filter
- `frontend/src/pages/GeneratorPage.tsx` — Card wrapping
- `frontend/src/components/analyzer/FileUpload.tsx` — visual polish
- `frontend/src/components/analyzer/SummaryCards.tsx` — shadcn Card/Badge/Skeleton
- `frontend/src/components/analyzer/TimeSeriesChart.tsx` — shadcn Card, CSS vars
- `frontend/src/components/analyzer/FilterBar.tsx` — shadcn Input/Checkbox
- `frontend/src/components/analyzer/LogTable.tsx` — Badge, alt rows
- `frontend/src/components/analyzer/PatternList.tsx` — shadcn Card
- `frontend/src/components/generator/ConfigForm.tsx` — shadcn Input/Tabs
- `frontend/src/components/generator/LivePreview.tsx` — shadcn Card, gutter
- `frontend/src/components/generator/ActionBar.tsx` — Tooltip, separator
- `frontend/src/components/ui/*` — new shadcn scaffolds

## Implementation

### Phase 1: Foundation

#### 1a. Color palette (`frontend/src/index.css`)

Shift all tokens from golden-amber to cool blue. Blue-tinted dark surfaces.

**Dark mode (`.dark`):**

| Token | New Value |
|-------|-----------|
| `--primary` | `oklch(0.65 0.15 250)` |
| `--primary-foreground` | `oklch(0.98 0 0)` |
| `--background` | `oklch(0.09 0.005 260)` |
| `--card` | `oklch(0.13 0.006 260)` |
| `--border` | `oklch(0.22 0.01 260)` |
| `--input` | `oklch(0.18 0.008 260)` |
| `--muted` | `oklch(0.16 0.006 260)` |
| `--muted-foreground` | `oklch(0.52 0.01 260)` |
| `--secondary` | `oklch(0.16 0.006 260)` |
| `--accent` | `oklch(0.65 0.15 250 / 0.15)` |
| `--accent-foreground` | `oklch(0.72 0.12 250)` |
| `--ring` | `oklch(0.65 0.15 250 / 0.4)` |
| `--chart-1` | `oklch(0.65 0.15 250)` |
| `--chart-3` | `oklch(0.72 0.13 180)` |
| Sidebar tokens | Mirror new primaries |

**Light mode (`:root`):**

| Token | New Value |
|-------|-----------|
| `--primary` | `oklch(0.50 0.18 250)` |
| `--primary-foreground` | `oklch(0.98 0 0)` |
| `--background` | `oklch(0.98 0.002 260)` |
| `--card` | `oklch(0.995 0.001 260)` |
| `--border` | `oklch(0.90 0.005 260)` |
| `--accent` | `oklch(0.50 0.18 250 / 0.10)` |
| `--accent-foreground` | `oklch(0.50 0.18 250)` |
| `--ring` | `oklch(0.50 0.18 250 / 0.4)` |
| Sidebar tokens | Mirror new primaries |

#### 1b. Scaffold shadcn/ui components

`npx shadcn add input card badge separator skeleton checkbox tabs tooltip`

Scaffolds into `frontend/src/components/ui/`. No new npm packages — uses existing `radix-ui` + `cva`.

### Phase 2: Shell

#### 2a. Navigation (`Navigation.tsx`)
- Active link: underline indicator (`border-b-2 border-primary`) instead of filled bg pill
- Stronger glass: `backdrop-blur-md`, `shadow-[0_1px_0_0_rgba(0,0,0,0.1)]`
- Vertical `Separator` between logo and nav links
- Nav height `h-14`, transition `duration-100`

#### 2b. Layout (`Layout.tsx`)
- Add `antialiased` to root div
- Add `selection:bg-primary/20`

### Phase 3: Analyzer Page

#### 3a. AnalyzerPage
- `Separator` between major sections
- Compact upload state when results visible: file name + "re-analyze" button
- Sticky filter bar: `sticky top-14 z-[5] bg-background/95 backdrop-blur-sm py-2`

#### 3b. FileUpload
- Solid thin border `border-border/50` + subtle fill `bg-muted/20` (no dashes)
- Hover: `border-primary/40 bg-primary/5`
- Drag-over: `ring-2 ring-primary/20`
- Icon in circular container: `bg-muted/50 p-3 rounded-full`
- Uploading: thin CSS shimmer bar at top

#### 3c. SummaryCards
- Handwritten `Card` → shadcn `Card`/`CardHeader`/`CardContent`
- shadcn `Skeleton` for loading state
- `Badge` for level counts and top sources
- `hover:border-primary/20` transition

#### 3d. TimeSeriesChart
- Wrap in shadcn `Card` with `CardHeader`/`CardTitle`/`CardDescription`
- Replace hardcoded oklch with CSS variable references
- Tooltip styled to `--popover`/`--popover-foreground`
- Gradient fill under bar chart via Recharts `<defs>` + `<linearGradient>`

#### 3e. FilterBar
- shadcn `Input` with lucide `Search` icon
- shadcn `Checkbox` replacing raw checkboxes
- shadcn `Input` for datetime-local
- Vertical `Separator` between sections
- "Clear filters": `Button variant="ghost" size="xs"` + `X` icon
- Wrapper: `bg-card/50 backdrop-blur-sm`

#### 3f. LogTable
- Header: `uppercase tracking-wider`
- Level column: `Badge` (destructive=error, outline+yellow=warn, default=info, secondary=debug)
- Alternating row tint: `bg-muted/10` on even rows
- Hover: `bg-muted/40`
- Empty state: lucide `SearchX` icon
- `Separator` above footer

#### 3g. PatternList
- Wrap in shadcn `Card`
- Header: `CardHeader` pattern
- Hit count in `Badge`
- Toggle: `Button variant="ghost" size="xs"`
- Expanded sample: `border-l-2 border-primary/30`

### Phase 4: Generator Page

#### 4a. GeneratorPage
- Config panel: shadcn `Card` with `CardHeader`/`CardContent`

#### 4b. ConfigForm
- All `<input>` → shadcn `Input`
- Format toggle → shadcn `Tabs` (`TabsList` + `TabsTrigger`)
- Level distribution: `Badge` labels, `Input type="number"` in 2x2 grid
- Validation errors: `bg-destructive/10 border-destructive/20 rounded-lg p-3` + `TriangleAlert` icon
- Submit: `Loader2` spinner when generating

#### 4c. LivePreview
- Wrap in shadcn `Card`
- Status dot: `bg-green-500 size-2 rounded-full animate-pulse` when generating
- Progress bar: `h-1 rounded-full`
- Auto-scroll toggle: `Button variant="ghost" size="xs"` + `ArrowDown` icon
- Line gutter: muted line numbers
- Empty state: lucide `Terminal` icon

#### 4d. ActionBar
- `Separator` above bar
- `Tooltip` on each button
- "Send to Analyzer": `shadow-sm shadow-primary/20`
- Gap `gap-3`

### Phase 5: Polish

- Consistent `transition-colors duration-150` on all interactive elements
- Verify light mode renders correctly
- `npm run build` — catch type/import errors
- Visual QA: spacing consistency (`gap-3` tight, `gap-6` sections)

## Acceptance Criteria

- [ ] Cool blue palette applied in both light and dark modes
- [ ] 8 shadcn/ui components scaffolded (input, card, badge, separator, skeleton, checkbox, tabs, tooltip)
- [ ] Navigation uses underline active state instead of filled pill
- [ ] All analyzer components use shadcn Card/Badge/Input/Checkbox
- [ ] All generator components use shadcn Card/Input/Tabs
- [ ] Filter bar is sticky below nav
- [ ] LogTable has alternating row tints and Badge for levels
- [ ] FileUpload uses solid border style with circular icon
- [ ] TimeSeriesChart uses CSS variable colors and gradient fill
- [ ] `npm run build` passes
- [ ] `npm run test` passes
- [ ] Responsive layout works at md and sm breakpoints
