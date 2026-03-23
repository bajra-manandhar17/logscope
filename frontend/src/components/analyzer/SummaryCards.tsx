import { useAnalyzerStore } from '../../stores/analyzerStore'
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

  return (
    <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
      {/* Total lines */}
      <Card className="transition-colors hover:border-primary/20">
        <CardHeader>
          <CardTitle className="text-xs font-medium uppercase tracking-wide text-muted-foreground">
            Total Lines
          </CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col gap-1">
          <span className="text-3xl font-semibold tabular-nums tracking-tight">
            {summary.total_lines.toLocaleString()}
          </span>
          <span className="text-xs text-muted-foreground">{result.format_detected} format</span>
        </CardContent>
      </Card>

      {/* Level breakdown */}
      <Card className="transition-colors hover:border-primary/20">
        <CardHeader>
          <CardTitle className="text-xs font-medium uppercase tracking-wide text-muted-foreground">
            Level Breakdown
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col gap-1.5">
            {(['error', 'warn', 'info', 'debug'] as const).map((level) => {
              const count = summary[`${level}_count` as keyof typeof summary] as number
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
      <Card className="transition-colors hover:border-primary/20">
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
      <Card className="transition-colors hover:border-primary/20">
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
