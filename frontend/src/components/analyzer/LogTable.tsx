import { useMemo, useRef, useState } from 'react'
import { useVirtualizer } from '@tanstack/react-virtual'
import { SearchX } from 'lucide-react'
import { useAnalyzerStore } from '../../stores/analyzerStore'
import type { LogEntry } from '../../types'
import { cn } from '@/lib/utils'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'

type SortKey = 'line_number' | 'timestamp' | 'level'
type SortDir = 'asc' | 'desc'

const levelBadge: Record<string, { variant: 'destructive' | 'outline' | 'default' | 'secondary'; className?: string }> = {
  error: { variant: 'destructive' },
  warn:  { variant: 'outline', className: 'border-yellow-500/30 text-yellow-500' },
  info:  { variant: 'default' },
  debug: { variant: 'secondary' },
}

function formatTs(iso: string): string {
  if (!iso) return '—'
  const d = new Date(iso)
  return isNaN(d.getTime()) ? iso : d.toLocaleTimeString()
}

const LEVEL_ORDER: Record<string, number> = { error: 0, warn: 1, info: 2, debug: 3 }

function sortEntries(entries: LogEntry[], key: SortKey, dir: SortDir): LogEntry[] {
  return [...entries].sort((a, b) => {
    let cmp = 0
    if (key === 'timestamp') {
      cmp = (a.timestamp ?? '').localeCompare(b.timestamp ?? '')
    } else if (key === 'level') {
      cmp = (LEVEL_ORDER[a.level] ?? 99) - (LEVEL_ORDER[b.level] ?? 99)
    } else {
      cmp = a.line_number - b.line_number
    }
    return dir === 'asc' ? cmp : -cmp
  })
}

const ROW_HEIGHT = 32

export function LogTable() {
  const result = useAnalyzerStore((s) => s.result)
  const filters = useAnalyzerStore((s) => s.filters)

  const [sortKey, setSortKey] = useState<SortKey>('line_number')
  const [sortDir, setSortDir] = useState<SortDir>('asc')

  const parentRef = useRef<HTMLDivElement>(null)

  const filtered = useMemo(() => {
    if (!result) return []
    let rows = result.entries

    if (filters.levels.length > 0) {
      rows = rows.filter((e) => filters.levels.includes(e.level))
    }
    if (filters.search) {
      const q = filters.search.toLowerCase()
      rows = rows.filter((e) => e.message.toLowerCase().includes(q))
    }
    if (filters.timeStart) {
      const start = new Date(filters.timeStart).getTime()
      rows = rows.filter((e) => e.timestamp && new Date(e.timestamp).getTime() >= start)
    }
    if (filters.timeEnd) {
      const end = new Date(filters.timeEnd).getTime()
      rows = rows.filter((e) => e.timestamp && new Date(e.timestamp).getTime() <= end)
    }

    return sortEntries(rows, sortKey, sortDir)
  }, [result, filters, sortKey, sortDir])

  const virtualizer = useVirtualizer({
    count: filtered.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => ROW_HEIGHT,
    overscan: 10,
  })

  if (!result) return null

  function handleSort(key: SortKey) {
    if (sortKey === key) {
      setSortDir((d) => (d === 'asc' ? 'desc' : 'asc'))
    } else {
      setSortKey(key)
      setSortDir('asc')
    }
  }

  function sortIndicator(key: SortKey) {
    if (sortKey !== key) return <span className="opacity-25">↕</span>
    return <span className="text-primary">{sortDir === 'asc' ? '↑' : '↓'}</span>
  }

  return (
    <div className="rounded-xl border border-border overflow-hidden flex flex-col bg-card">
      {/* Header */}
      <div className="grid grid-cols-[4rem_9rem_5rem_9rem_1fr] gap-x-3 bg-muted/40 px-3 py-2 text-xs font-medium uppercase tracking-wider text-muted-foreground select-none">
        <button className="text-left flex items-center gap-1" onClick={() => handleSort('line_number')}>
          # {sortIndicator('line_number')}
        </button>
        <button className="text-left flex items-center gap-1" onClick={() => handleSort('timestamp')}>
          Time {sortIndicator('timestamp')}
        </button>
        <button className="text-left flex items-center gap-1" onClick={() => handleSort('level')}>
          Level {sortIndicator('level')}
        </button>
        <span>Source</span>
        <span>Message</span>
      </div>

      {/* Body */}
      <div ref={parentRef} className="overflow-auto" style={{ height: '480px' }}>
        {filtered.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-full gap-2 text-muted-foreground">
            <SearchX className="size-8 opacity-40" />
            <span className="text-sm">No entries match the current filters</span>
          </div>
        ) : (
          <div style={{ height: `${virtualizer.getTotalSize()}px`, position: 'relative' }}>
            {virtualizer.getVirtualItems().map((vItem) => {
              const entry = filtered[vItem.index]
              const cfg = levelBadge[entry.level]
              return (
                <div
                  key={vItem.key}
                  style={{ position: 'absolute', top: vItem.start, width: '100%', height: ROW_HEIGHT }}
                  className={cn(
                    'grid grid-cols-[4rem_9rem_5rem_9rem_1fr] gap-x-3 px-3 items-center text-xs border-b border-border/50 transition-colors duration-100',
                    vItem.index % 2 === 0 ? 'bg-transparent' : 'bg-muted/10',
                    'hover:bg-muted/40',
                  )}
                >
                  <span className="tabular-nums text-muted-foreground">{entry.line_number}</span>
                  <span className="tabular-nums text-muted-foreground truncate">{formatTs(entry.timestamp)}</span>
                  <Badge variant={cfg?.variant ?? 'secondary'} className={cn('text-[10px] uppercase', cfg?.className)}>
                    {entry.level}
                  </Badge>
                  <span className="truncate font-mono text-muted-foreground">{entry.source || '—'}</span>
                  <span className="truncate">{entry.message || entry.raw}</span>
                </div>
              )
            })}
          </div>
        )}
      </div>

      <Separator />
      <div className="px-3 py-1.5 text-xs text-muted-foreground bg-muted/30">
        {filtered.length.toLocaleString()} of {result.entries.length.toLocaleString()} entries
      </div>
    </div>
  )
}
