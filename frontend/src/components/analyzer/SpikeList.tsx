import type { Spike } from '../../types'
import { Badge } from '@/components/ui/badge'

function formatTs(iso: string): string {
  const d = new Date(iso)
  return isNaN(d.getTime()) ? iso : d.toLocaleString()
}

export function SpikeList({ spikes }: { spikes: Spike[] }) {
  if (!spikes?.length) return null

  return (
    <div className="overflow-y-auto max-h-64">
      <table className="w-full text-xs">
        <thead>
          <tr className="border-b border-border text-muted-foreground">
            <th className="text-left py-2 px-3 font-medium">Time</th>
            <th className="text-right py-2 px-3 font-medium">Count</th>
            <th className="text-right py-2 px-3 font-medium">Threshold</th>
            <th className="text-left py-2 px-3 font-medium">Severity</th>
          </tr>
        </thead>
        <tbody>
          {spikes.map((s, i) => (
            <tr key={i} className="border-b border-border/50 last:border-0">
              <td className="py-2 px-3 tabular-nums text-muted-foreground">{formatTs(s.bucket_timestamp)}</td>
              <td className="py-2 px-3 text-right tabular-nums font-medium">{s.count.toLocaleString()}</td>
              <td className="py-2 px-3 text-right tabular-nums text-muted-foreground">{Math.round(s.threshold).toLocaleString()}</td>
              <td className="py-2 px-3">
                <Badge variant={s.severity === 'high' ? 'destructive' : 'outline'} className="text-[10px] uppercase">
                  {s.severity}
                </Badge>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
