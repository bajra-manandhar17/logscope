import { useEffect, useRef, useState } from 'react'
import { useVirtualizer } from '@tanstack/react-virtual'
import { ArrowDown, Terminal } from 'lucide-react'
import { useGeneratorStore } from '../../stores/generatorStore'
import { cn } from '@/lib/utils'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'

const ROW_HEIGHT = 20

function detectLevel(line: string): string {
  if (line.trimStart().startsWith('{')) {
    try {
      const m = JSON.parse(line) as Record<string, unknown>
      return typeof m.level === 'string' ? m.level.toLowerCase() : ''
    } catch {
      return ''
    }
  }
  if (/\bERROR\b/.test(line)) return 'error'
  if (/\bWARN\b/.test(line)) return 'warn'
  if (/\bINFO\b/.test(line)) return 'info'
  if (/\bDEBUG\b/.test(line)) return 'debug'
  return ''
}

const levelColor: Record<string, string> = {
  error: 'text-red-400',
  warn:  'text-yellow-400',
  info:  'text-blue-400',
  debug: 'text-muted-foreground',
}

const SCROLL_THRESHOLD = 48

export function LivePreview() {
  const lines = useGeneratorStore((s) => s.lines)
  const status = useGeneratorStore((s) => s.status)
  const totalRequested = useGeneratorStore((s) => s.config.total_lines)

  const parentRef = useRef<HTMLDivElement>(null)
  const [pinned, setPinned] = useState(true)

  const virtualizer = useVirtualizer({
    count: lines.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => ROW_HEIGHT,
    overscan: 20,
  })

  useEffect(() => {
    if (pinned && lines.length > 0) {
      virtualizer.scrollToIndex(lines.length - 1, { align: 'end' })
    }
  }, [lines.length, pinned, virtualizer])

  function onScroll() {
    const el = parentRef.current
    if (!el) return
    const atBottom = el.scrollTop + el.clientHeight >= el.scrollHeight - SCROLL_THRESHOLD
    setPinned(atBottom)
  }

  const progress = totalRequested > 0 ? Math.min(lines.length / totalRequested, 1) : 0
  const generating = status === 'generating'

  return (
    <Card className="overflow-hidden">
      <CardHeader className="flex-row items-center justify-between gap-3 py-2">
        <div className="flex items-center gap-2 min-w-0">
          {generating && (
            <span className="size-2 rounded-full bg-green-500 animate-pulse shrink-0" />
          )}
          <CardTitle className="text-sm shrink-0">Live Preview</CardTitle>
          {lines.length > 0 && (
            <span className="text-xs text-muted-foreground tabular-nums">
              {lines.length.toLocaleString()}
              {status !== 'idle' && ` / ${totalRequested.toLocaleString()}`}
            </span>
          )}
        </div>

        <div className="flex items-center gap-2 shrink-0">
          {lines.length > 0 && (
            <Button
              variant="ghost"
              size="xs"
              onClick={() => {
                const next = !pinned
                setPinned(next)
                if (next) virtualizer.scrollToIndex(lines.length - 1, { align: 'end' })
              }}
              className={pinned ? 'text-primary' : ''}
            >
              <ArrowDown />
              {pinned ? 'Auto-scroll on' : 'Auto-scroll off'}
            </Button>
          )}
          {generating && (
            <span className="text-xs text-muted-foreground animate-pulse">generating…</span>
          )}
        </div>
      </CardHeader>

      {/* Progress bar */}
      {(generating || status === 'done') && lines.length > 0 && (
        <div className="h-1 bg-muted/60 mx-4 rounded-full overflow-hidden">
          <div
            className="h-full bg-primary rounded-full transition-all duration-300 ease-out"
            style={{ width: `${progress * 100}%` }}
          />
        </div>
      )}

      {/* Virtualized lines */}
      <CardContent className="p-0">
        <div
          ref={parentRef}
          onScroll={onScroll}
          className="overflow-auto font-mono text-xs"
          style={{ height: '400px' }}
        >
          {lines.length === 0 ? (
            <div className="flex flex-col items-center justify-center h-full gap-2 text-muted-foreground">
              <Terminal className="size-8 opacity-40" />
              <span className="text-sm">
                {status === 'idle' ? 'Configure and generate to see a preview' : 'Waiting for first batch…'}
              </span>
            </div>
          ) : (
            <div style={{ height: `${virtualizer.getTotalSize()}px`, position: 'relative' }}>
              {virtualizer.getVirtualItems().map((vItem) => {
                const line = lines[vItem.index]
                const level = detectLevel(line)
                return (
                  <div
                    key={vItem.key}
                    style={{ position: 'absolute', top: vItem.start, width: '100%', height: ROW_HEIGHT }}
                    className={cn(
                      'px-4 flex items-center gap-2 truncate leading-none',
                    )}
                  >
                    <span className="w-10 text-right text-muted-foreground/50 tabular-nums shrink-0 select-none">
                      {vItem.index + 1}
                    </span>
                    <span className={levelColor[level] ?? 'text-foreground'}>{line}</span>
                  </div>
                )
              })}
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  )
}
