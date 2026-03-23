import {
  ComposedChart,
  Bar,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts'
import { useAnalyzerStore } from '../../stores/analyzerStore'
import type { TimeBucket } from '../../types'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'

function formatTick(iso: string, interval: string): string {
  const d = new Date(iso)
  if (isNaN(d.getTime())) return iso
  if (interval.includes('h') || interval.includes('d')) {
    return d.toLocaleDateString(undefined, { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
  }
  return d.toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit' })
}

function formatTooltipLabel(iso: string): string {
  const d = new Date(iso)
  return isNaN(d.getTime()) ? iso : d.toLocaleString()
}

interface ChartRow {
  ts: string
  total: number
  errors: number
}

function toChartData(buckets: TimeBucket[]): ChartRow[] {
  return buckets.map((b) => ({
    ts: b.timestamp,
    total: b.count,
    errors: b.error_count,
  }))
}

function getCSSVar(name: string): string {
  return getComputedStyle(document.documentElement).getPropertyValue(name).trim()
}

export function TimeSeriesChart() {
  const result = useAnalyzerStore((s) => s.result)

  if (!result) return null

  const { time_series, bucket_interval } = result

  if (!time_series?.length) {
    return (
      <div className="flex items-center justify-center h-40 rounded-xl border border-border text-sm text-muted-foreground">
        No time-series data available
      </div>
    )
  }

  const data = toChartData(time_series)

  const chartColor = getCSSVar('--chart-1') || 'oklch(0.65 0.15 250)'
  const errorColor = getCSSVar('--chart-2') || 'oklch(0.65 0.22 25)'
  const gridColor = getCSSVar('--border') || 'oklch(0.22 0.01 260)'
  const textColor = getCSSVar('--muted-foreground') || 'oklch(0.52 0.01 260)'
  const popoverBg = getCSSVar('--popover') || 'oklch(0.13 0.006 260)'
  const popoverFg = getCSSVar('--popover-foreground') || 'oklch(0.9 0.003 260)'

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-sm">Log Volume Over Time</CardTitle>
        {bucket_interval && (
          <CardDescription>Bucket interval: {bucket_interval}</CardDescription>
        )}
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={240}>
          <ComposedChart data={data} margin={{ top: 4, right: 16, left: 0, bottom: 4 }}>
            <defs>
              <linearGradient id="barGradient" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stopColor={chartColor} stopOpacity={0.6} />
                <stop offset="100%" stopColor={chartColor} stopOpacity={0.15} />
              </linearGradient>
            </defs>
            <CartesianGrid strokeDasharray="3 3" stroke={gridColor} />
            <XAxis
              dataKey="ts"
              tickFormatter={(v) => formatTick(v, bucket_interval)}
              tick={{ fontSize: 11, fill: textColor }}
              tickLine={false}
              axisLine={false}
            />
            <YAxis
              tick={{ fontSize: 11, fill: textColor }}
              tickLine={false}
              axisLine={false}
              width={40}
            />
            <Tooltip
              labelFormatter={(label) => formatTooltipLabel(String(label))}
              formatter={(value, name) => [
                typeof value === 'number' ? value.toLocaleString() : value,
                name === 'total' ? 'Total' : 'Errors',
              ]}
              contentStyle={{
                fontSize: 12,
                backgroundColor: popoverBg,
                border: `1px solid ${gridColor}`,
                borderRadius: '8px',
                color: popoverFg,
              }}
            />
            <Legend
              formatter={(v) => (v === 'total' ? 'Total logs' : 'Errors')}
              wrapperStyle={{ fontSize: 12 }}
            />
            <Bar dataKey="total" fill="url(#barGradient)" radius={[2, 2, 0, 0]} name="total" />
            <Line
              type="monotone"
              dataKey="errors"
              stroke={errorColor}
              strokeWidth={2}
              dot={false}
              name="errors"
            />
          </ComposedChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  )
}
