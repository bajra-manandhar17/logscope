# 28 — Live Replay System

**Complexity:** Medium
**Phase:** 8 — Replay & Playback
**Blocked by:** 27-intelligence-system, 18-time-series-chart, 17-log-table-filter, 16-summary-cards
**Blocks:** none

## Objective

Add frontend-only time-cursor playback over existing `AnalysisResult`. User hits "Replay" and the UI progressively reveals log entries, sweeps a cursor across the time-series chart, and live-counts summary stats — simulating logs arriving in real time. No backend changes.

## Scope

### Pre-requisite

- Run `npx shadcn@latest add slider` to add shadcn Slider component

### New files

- `frontend/src/types/index.ts` — add `ReplayMode`, `ReplaySpeed` types
- `frontend/src/lib/binarySearch.ts` — binary search utilities for timestamp lookup
- `frontend/src/stores/replayStore.ts` — Zustand clock store
- `frontend/src/hooks/useReplayLoop.ts` — `requestAnimationFrame` animation loop
- `frontend/src/components/analyzer/ReplayControlBar.tsx` — playback controls UI

### Modified files

- `frontend/src/pages/AnalyzerPage.tsx` — Replay button + control bar placement + hook
- `frontend/src/components/analyzer/TimeSeriesChart.tsx` — ReferenceLine cursor + dimmed future bars
- `frontend/src/components/analyzer/SummaryCards.tsx` — live counting mode
- `frontend/src/components/analyzer/LogTable.tsx` — progressive reveal + auto-scroll

### Test files (new)

- `frontend/src/lib/__tests__/binarySearch.test.ts`
- `frontend/src/stores/__tests__/replayStore.test.ts`

## Design Decisions

- Frontend-only — no backend changes; all data exists in `AnalysisResult`
- 3 components participate: TimeSeriesChart, LogTable, SummaryCards
- Controls: Play/Pause, Speed (1x/5x/10x/50x), Scrub bar (shadcn Slider), Reset
- Normalized speed: 1x = full replay in 60 seconds
- Inline control bar between SummaryCards and TimeSeriesChart
- "Replay" button post-analysis auto-starts playback
- Separate `replayStore` — thin clock; components filter independently via `useMemo`
- ~15fps tick throttle + binary search on sorted entries for performance
- Scrub: live seek, resumes on release; if paused before drag, stays paused
- Chart: ReferenceLine sweep + dimmed future bars via Recharts `<Cell>`

## Feature Details

### 1. Types (`frontend/src/types/index.ts`)

Add at bottom:
```typescript
export type ReplayMode = 'idle' | 'playing' | 'paused'
export type ReplaySpeed = 1 | 5 | 10 | 50
```

### 2. Binary Search Utility (`frontend/src/lib/binarySearch.ts`)

```typescript
export function upperBoundByTime<T>(arr: T[], targetMs: number, getTimeMs: (item: T) => number): number
export function entryIndexAtTime(entries: LogEntry[], time: string): number
export function bucketIndexAtTime(buckets: TimeBucket[], time: string): number
```

Standard binary search. Pre-parse target to ms outside loop. Entries without timestamps return -1 from convenience wrappers (treated as always-included by consumers).

### 3. Replay Store (`frontend/src/stores/replayStore.ts`)

Follow `analyzerStore.ts` Zustand pattern.

**State:** `mode: ReplayMode`, `speed: ReplaySpeed`, `currentTime: string | null`, `startTime: string | null`, `endTime: string | null`, `progress: number` (0–1)

**Actions:** `start(startTime, endTime)`, `pause()`, `resume()`, `seek(time)`, `setSpeed(speed)`, `tick(time, progress)`, `stop()`

Auto-stop: subscribe to `useAnalyzerStore` — when `result` becomes null, call `stop()`.

### 4. Animation Hook (`frontend/src/hooks/useReplayLoop.ts`)

`useReplayLoop()` — called unconditionally in AnalyzerPage; no-ops when `mode === 'idle'`.

Algorithm:
1. `requestAnimationFrame` loop
2. `deltaMs = performance.now() - lastFrameTime`
3. `timeScaleMs = (endMs - startMs) / 60_000` (maps log span to 60s at 1x)
4. `currentMs += deltaMs * speed * timeScaleMs`
5. Clamp to `[startMs, endMs]`
6. Throttle `tick()` calls to ~66ms intervals (~15fps)
7. At end: auto-pause (not stop), progress = 1
8. Guard: `endMs <= startMs` → complete instantly

Cleanup: `cancelAnimationFrame` on unmount or mode → idle.

### 5. ReplayControlBar (`frontend/src/components/analyzer/ReplayControlBar.tsx`)

Renders only when `mode !== 'idle'`. Single-row flex layout:
```
[ Play/Pause ] [ Speed: Nx ] [ ──●────── scrub ] [ HH:MM:SS ] [ Reset ]
```

- Play/Pause: lucide `Play`/`Pause` icons
- Speed: cycle button `[1, 5, 10, 50]`, display "1x" etc
- Scrub: shadcn `<Slider>` — `value={[progress * 1000]}`, `max={1000}`, `step={1}`
  - `onValueChange`: interpolate ISO time → `seek()`
  - `onPointerDown`: if playing, store `wasPlaying`, pause
  - `onPointerUp`: if `wasPlaying`, resume
- Timestamp: formatted `currentTime`
- Reset: lucide `RotateCcw` → `stop()`

Style: `bg-card border border-border rounded-lg px-4 py-2`

### 6. AnalyzerPage Wiring (`frontend/src/pages/AnalyzerPage.tsx`)

- Call `useReplayLoop()` unconditionally
- "Replay" button next to "Analysis Results" heading (lucide `Play` icon)
  - Visible when `result.summary.time_range[0]` is truthy and `mode === 'idle'`
  - On click: `replayStore.start(time_range[0], time_range[1])`
- `<ReplayControlBar />` between `<SummaryCards />` and first `<Separator />`

Layout:
```
SummaryCards
ReplayControlBar  ← conditional
Separator
TimeSeriesChart
...rest unchanged
```

### 7. TimeSeriesChart Modifications

- Import `ReferenceLine`, `Cell` from recharts
- When replaying:
  - `cutoffIdx = bucketIndexAtTime(time_series, currentTime)`
  - Add `<Cell fillOpacity={i > cutoffIdx ? 0.15 : 1} />` inside `<Bar>`
  - Add `<ReferenceLine x={currentTime} stroke="oklch(0.65 0.2 150)" strokeDasharray="4 4" />`
- When idle: unchanged

### 8. SummaryCards Modifications

- When replaying: compute live counts via `useMemo` keyed on `[result, currentTime, mode]`
  - `cutoffIdx = entryIndexAtTime(result.entries, currentTime)`
  - Iterate `entries[0..cutoffIdx]` counting error/warn/info/debug
  - Total lines = `cutoffIdx + 1`
- When idle: unchanged (reads `summary` directly)
- Entries without timestamps always counted

### 9. LogTable Modifications

- Add replay filter AFTER existing filters in `filtered` useMemo:
  ```typescript
  if (mode !== 'idle' && currentTime) {
    const cutoff = new Date(currentTime).getTime()
    rows = rows.filter(e => !e.timestamp || new Date(e.timestamp).getTime() <= cutoff)
  }
  ```
  Uses `.filter()` not binary search (existing filters break sort order; 10K at 15fps is ~0.5ms).
- Auto-scroll: `useEffect` when `mode === 'playing'` → `virtualizer.scrollToIndex(filtered.length - 1, { align: 'end' })`
- No auto-scroll when paused

## Edge Cases

| Case | Handling |
|------|----------|
| No timestamps in data | Replay button hidden (`time_range[0]` falsy) |
| Entries without timestamps | Always visible during replay (pass filter) |
| Replay reaches end | Auto-pause, progress=1, user can scrub back or reset |
| New file upload during replay | Store subscription detects result change → `stop()` |
| `endTime === startTime` | Guard: complete instantly, auto-pause |
| FilterBar + replay | Compose naturally; both filters apply, no interference |

## Acceptance Criteria

- [x] `npx shadcn@latest add slider` installs without errors
- [x] `upperBoundByTime` returns correct indices for edge cases (empty, single, boundary, no-timestamp)
- [x] `replayStore` start/pause/resume/seek/stop/tick/setSpeed work correctly
- [x] Auto-stop triggers when analyzer result becomes null
- [x] "Replay" button appears only when timestamps exist; hidden during active replay
- [x] Control bar renders with play/pause, speed, scrub, timestamp, reset
- [x] Speed selector cycles through 1x/5x/10x/50x
- [x] Scrub bar seeks live; paused state preserved across drag
- [x] TimeSeriesChart shows sweeping ReferenceLine + dimmed future bars
- [x] SummaryCards show incrementing counts during replay
- [x] LogTable progressively reveals entries + auto-scrolls during playback
- [x] Filters compose with replay (level, search, time range all work during playback)
- [x] Replay auto-pauses at end
- [x] New upload during replay triggers clean stop
- [x] `cd frontend && npm run build` — no type errors (2 pre-existing errors in IntelligencePanel + analyzerStore.test excluded)
- [x] `cd frontend && npx vitest run` — all tests pass (45/45)
