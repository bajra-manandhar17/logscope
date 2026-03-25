import type { CausalSequence } from '../../types'
import { Badge } from '@/components/ui/badge'

export function CausalSequenceList({ sequences }: { sequences: CausalSequence[] }) {
  if (!sequences?.length) return null

  return (
    <div className="overflow-y-auto max-h-64 flex flex-col gap-2 p-3">
      {sequences.map((seq, i) => (
        <div key={i} className="flex items-center gap-2 text-xs border border-border/50 rounded-lg px-3 py-2">
          <span className="text-muted-foreground min-w-[1.5rem] text-right">#{i + 1}</span>
          <span className="font-mono truncate flex-1" title={seq.pattern_a}>{seq.pattern_a}</span>
          <span className="text-muted-foreground shrink-0">→</span>
          <span className="font-mono truncate flex-1" title={seq.pattern_b}>{seq.pattern_b}</span>
          <Badge variant="outline" className="tabular-nums shrink-0">{seq.count}×</Badge>
          <span className="text-muted-foreground tabular-nums shrink-0">~{seq.avg_lag_seconds.toFixed(1)}s</span>
        </div>
      ))}
    </div>
  )
}
