import { useAnalyzerStore } from '../stores/analyzerStore'
import { useReplayStore } from '../stores/replayStore'
import { useReplayLoop } from '../hooks/useReplayLoop'
import { FileUpload } from '../components/analyzer/FileUpload'
import { SummaryCards } from '../components/analyzer/SummaryCards'
import { FilterBar } from '../components/analyzer/FilterBar'
import { LogTable } from '../components/analyzer/LogTable'
import { TimeSeriesChart } from '../components/analyzer/TimeSeriesChart'
import { PatternList } from '../components/analyzer/PatternList'
import { IntelligencePanel } from '../components/analyzer/IntelligencePanel'
import { ReplayControlBar } from '../components/analyzer/ReplayControlBar'
import { Separator } from '@/components/ui/separator'
import { Play } from 'lucide-react'

export function AnalyzerPage() {
  const result = useAnalyzerStore((s) => s.result)
  const status = useAnalyzerStore((s) => s.status)
  const replayMode = useReplayStore((s) => s.mode)
  const { start } = useReplayStore.getState()

  useReplayLoop()

  const hasResult = result !== null
  const isUploading = status === 'uploading'
  const canReplay = hasResult && !!result.summary.time_range?.[0] && replayMode === 'idle'

  return (
    <div className="px-4 py-3 flex flex-col gap-3">
      {/* Upload zone — hero when no result, compact when results visible */}
      <section className={hasResult ? 'flex flex-col items-center gap-2' : 'flex flex-col items-center justify-center py-16 gap-8'}>
        {!hasResult && (
          <div className="text-center">
            <h1 className="text-xl font-semibold tracking-tight text-foreground">
              Log Analyzer
            </h1>
            <p className="text-muted-foreground mt-2 text-sm">
              Upload a log file to extract patterns, timelines, and entries
            </p>
          </div>
        )}
        <FileUpload />
        {hasResult && (
          <p className="text-xs text-muted-foreground">Drop a new file to re-analyze</p>
        )}
      </section>

      {/* Skeleton cards while uploading */}
      {isUploading && !hasResult && <SummaryCards />}

      {/* Results */}
      {hasResult && (
        <>
          <Separator />
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-base font-semibold tracking-tight">
                Analysis Results
              </h2>
              <p className="text-xs text-muted-foreground mt-0.5">
                {result.format_detected} format · {result.summary.total_lines.toLocaleString()} lines
              </p>
            </div>
            {canReplay && (
              <button
                onClick={() => start(result.summary.time_range[0], result.summary.time_range[1])}
                className="flex items-center gap-1.5 text-xs px-3 py-1.5 rounded-md border border-border bg-card hover:border-primary/50 hover:text-primary transition-colors"
              >
                <Play className="size-3" />
                Replay
              </button>
            )}
          </div>
          <SummaryCards />
          <ReplayControlBar />
          <Separator />
          <TimeSeriesChart />
          <Separator />
          <IntelligencePanel />
          <Separator />
          <div className="flex flex-col gap-3">
            <div className="sticky top-10 z-[5] bg-background py-2 -mx-4 px-4">
              <FilterBar />
            </div>
            <LogTable />
          </div>
          <Separator />
          <PatternList />
        </>
      )}
    </div>
  )
}
