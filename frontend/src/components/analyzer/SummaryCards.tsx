import { useMemo } from 'react'
import { useAnalyzerStore } from '../../stores/analyzerStore'
import { useReplayStore } from '../../stores/replayStore'
import { entryIndexAtTime } from '../../lib/binarySearch'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'

function formatTime(iso: string): string {
  if (!iso) return '—'
  const d = new Date(iso)
  return isNaN(d.getTime()) ? iso : d.toLocaleString()
}

const levelConfig: Record<string, { label: string; variant: 'destructive' | 'outline' | 'default' | 'secondary'; className?: string }> = {
  error: { label: 'Error', variant: 'destructive' },
  warn:  { label: 'Warn',  variant: 'outline', className: 'border-yellow-500/30 text-yellow-500' },
  info:  { label: 'Info',  variant: 'default' },
  debug: { label: 'Debug', variant: 'secondary' },
}

export function SummaryCards() {
  const result = useAnalyzerStore((s) => s.result)
  const status = useAnalyzerStore((s) => s.status)
  const replayMode = useReplayStore((s) => s.mode)
  const currentTime = useReplayStore((s) => s.currentTime)

  const liveCounts = useMemo(() => {
    if (!result || replayMode === 'idle' || !currentTime) return null
    const cutoffIdx = entryIndexAtTime(result.entries, currentTime)
    const counts = { total: 0, error: 0, warn: 0, info: 0, debug: 0 }
    const limit = cutoffIdx < 0 ? 0 : cutoffIdx + 1
    for (let i = 0; i < limit; i++) {
      const e = result.entries[i]
      // Entries without timestamps are always counted
      if (e.timestamp && new Date(e.timestamp).getTime() > new Date(currentTime).getTime()) continue
      counts.total++
      const lvl = e.level as keyof typeof counts
      if (lvl in counts && lvl !== 'total') counts[lvl]++
    }
    return counts
  }, [result, currentTime, replayMode])

  if (status === 'uploading') {
    return (
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
        {Array.from({ length: 4 }).map((_, i) => (
          <Card key={i} className="transition-colors">
            <CardHeader>
              <Skeleton className="h-3 w-20" />
            </CardHeader>
            <CardContent>
              <Skeleton className="h-8 w-24" />
            </CardContent>
          </Card>
        ))}
      </div>
    )
  }

  if (!result) return null

  const { summary } = result
  const [start, end] = summary.time_range ?? []

  const displayTotal = liveCounts ? liveCounts.total : summary.total_lines

  return (
    <div className="grid grid-cols-2 lg:grid-cols-4 gap-2">
      {/* Total lines */}
      <Card className="transition-colors">
        <CardHeader>
          <CardTitle className="text-xs font-medium uppercase tracking-wide text-muted-foreground">
            Total Lines
          </CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col gap-1">
          <span className="text-2xl font-mono font-semibold tabular-nums tracking-tight">
            {displayTotal.toLocaleString()}
          </span>
          <span className="text-xs text-muted-foreground">{result.format_detected} format</span>
        </CardContent>
      </Card>

      {/* Level breakdown */}
      <Card className="transition-colors">
        <CardHeader>
          <CardTitle className="text-xs font-medium uppercase tracking-wide text-muted-foreground">
            Level Breakdown
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col gap-1.5">
            {(['error', 'warn', 'info', 'debug'] as const).map((level) => {
              const count = liveCounts
                ? liveCounts[level]
                : (summary[`${level}_count` as keyof typeof summary] as number)
              const cfg = levelConfig[level]
              return (
                <div key={level} className="flex items-center justify-between">
                  <Badge variant={cfg.variant} className={cfg.className}>
                    {cfg.label}
                  </Badge>
                  <span className="tabular-nums text-sm text-foreground">{count.toLocaleString()}</span>
                </div>
              )
            })}
          </div>
        </CardContent>
      </Card>

      {/* Time range */}
      <Card className="transition-colors">
        <CardHeader>
          <CardTitle className="text-xs font-medium uppercase tracking-wide text-muted-foreground">
            Time Range
          </CardTitle>
        </CardHeader>
        <CardContent>
          {start ? (
            <div className="flex flex-col gap-1 text-sm">
              <div>
                <span className="text-muted-foreground text-xs">From </span>
                <span className="font-medium">{formatTime(start)}</span>
              </div>
              <div>
                <span className="text-muted-foreground text-xs">To </span>
                <span className="font-medium">{formatTime(end)}</span>
              </div>
            </div>
          ) : (
            <span className="text-sm text-muted-foreground">No timestamps detected</span>
          )}
        </CardContent>
      </Card>

      {/* Top sources */}
      <Card className="transition-colors">
        <CardHeader>
          <CardTitle className="text-xs font-medium uppercase tracking-wide text-muted-foreground">
            Top Sources
          </CardTitle>
        </CardHeader>
        <CardContent>
          {summary.top_sources?.length ? (
            <div className="flex flex-wrap gap-1.5">
              {summary.top_sources.slice(0, 5).map((src) => (
                <Badge key={src} variant="outline" className="font-mono text-xs">
                  {src}
                </Badge>
              ))}
            </div>
          ) : (
            <span className="text-sm text-muted-foreground">No sources detected</span>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
