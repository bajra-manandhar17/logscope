import type { SilenceGap } from '../../types'
import { Badge } from '@/components/ui/badge'

function formatTs(iso: string): string {
  const d = new Date(iso)
  return isNaN(d.getTime()) ? iso : d.toLocaleString()
}

export function SilenceGapList({ gaps }: { gaps: SilenceGap[] }) {
  if (!gaps?.length) return null

  return (
    <div className="overflow-y-auto max-h-64">
      <table className="w-full text-xs">
        <thead>
          <tr className="border-b border-border text-muted-foreground">
            <th className="text-left py-2 px-3 font-medium">Source</th>
            <th className="text-left py-2 px-3 font-medium">Gap Start</th>
            <th className="text-left py-2 px-3 font-medium">Gap End</th>
            <th className="text-right py-2 px-3 font-medium">Duration</th>
            <th className="text-right py-2 px-3 font-medium">Active Sources</th>
          </tr>
        </thead>
        <tbody>
          {gaps.map((g, i) => (
            <tr key={i} className="border-b border-border/50 last:border-0">
              <td className="py-2 px-3">
                <Badge variant="outline" className="text-[10px] font-mono">{g.source}</Badge>
              </td>
              <td className="py-2 px-3 tabular-nums text-muted-foreground">{formatTs(g.gap_start)}</td>
              <td className="py-2 px-3 tabular-nums text-muted-foreground">{formatTs(g.gap_end)}</td>
              <td className="py-2 px-3 text-right font-medium">{g.duration}</td>
              <td className="py-2 px-3 text-right tabular-nums text-muted-foreground">{g.active_sources_during_gap}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
