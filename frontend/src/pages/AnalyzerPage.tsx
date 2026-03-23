import { useAnalyzerStore } from '../stores/analyzerStore'
import { FileUpload } from '../components/analyzer/FileUpload'
import { SummaryCards } from '../components/analyzer/SummaryCards'
import { FilterBar } from '../components/analyzer/FilterBar'
import { LogTable } from '../components/analyzer/LogTable'
import { TimeSeriesChart } from '../components/analyzer/TimeSeriesChart'
import { PatternList } from '../components/analyzer/PatternList'
import { Separator } from '@/components/ui/separator'

export function AnalyzerPage() {
  const result = useAnalyzerStore((s) => s.result)
  const status = useAnalyzerStore((s) => s.status)

  const hasResult = result !== null
  const isUploading = status === 'uploading'

  return (
    <div className="max-w-7xl mx-auto px-6 py-8 flex flex-col gap-8">
      {/* Upload zone — hero when no result, compact when results visible */}
      <section className={hasResult ? 'flex flex-col items-center gap-2' : 'flex flex-col items-center justify-center py-16 gap-8'}>
        {!hasResult && (
          <div className="text-center">
            <h1
              className="text-3xl font-bold tracking-tight text-foreground"
              style={{ fontFamily: 'var(--font-heading)' }}
            >
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
          <div>
            <h2
              className="text-lg font-bold tracking-tight"
              style={{ fontFamily: 'var(--font-heading)' }}
            >
              Analysis Results
            </h2>
            <p className="text-xs text-muted-foreground mt-0.5">
              {result.format_detected} format · {result.summary.total_lines.toLocaleString()} lines
            </p>
          </div>
          <SummaryCards />
          <Separator />
          <TimeSeriesChart />
          <Separator />
          <div className="flex flex-col gap-3">
            <div className="sticky top-14 z-[5] bg-background/95 backdrop-blur-sm py-2 -mx-6 px-6">
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
