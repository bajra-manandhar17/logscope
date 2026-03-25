import { useState } from 'react'
import { useAnalyzerStore } from '../../stores/analyzerStore'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { SpikeList } from './SpikeList'
import { SilenceGapList } from './SilenceGapList'
import { CausalSequenceList } from './CausalSequenceList'

interface SectionProps {
  title: string
  count: number
  children: React.ReactNode
}

function CollapsibleSection({ title, count, children }: SectionProps) {
  const [open, setOpen] = useState(false)

  if (count === 0) return null

  return (
    <div className="border border-border/50 rounded-lg overflow-hidden">
      <button
        onClick={() => setOpen((v) => !v)}
        className="w-full flex items-center justify-between px-4 py-2.5 text-sm font-medium hover:bg-muted/30 transition-colors"
      >
        <span className="flex items-center gap-2">
          <span className="text-muted-foreground">{open ? '▾' : '▸'}</span>
          {title}
        </span>
        <Badge variant="secondary" className="tabular-nums text-[10px]">{count}</Badge>
      </button>
      {open && <div className="border-t border-border/50">{children}</div>}
    </div>
  )
}

export function IntelligencePanel() {
  const result = useAnalyzerStore((s) => s.result)

  if (!result) return null

  const intel = result.intelligence
  if (!intel) return null

  const spikeCount = intel.spikes?.length ?? 0
  const silenceCount = intel.silence_gaps?.length ?? 0
  const causalCount = intel.causal_sequences?.length ?? 0
  const hasData = spikeCount > 0 || silenceCount > 0 || causalCount > 0 || intel.high_entropy_count > 0

  if (!hasData) {
    return (
      <div className="flex items-center justify-center h-24 rounded-xl border border-border text-sm text-muted-foreground">
        No intelligence insights detected — timestamps may be missing or log volume too uniform
      </div>
    )
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-sm">Intelligence Insights</CardTitle>
      </CardHeader>
      <CardContent className="flex flex-col gap-3">
        {/* Summary stats */}
        <div className="flex flex-wrap gap-4 text-xs text-muted-foreground pb-2">
          {spikeCount > 0 && <span>{spikeCount} spike{spikeCount !== 1 ? 's' : ''}</span>}
          {silenceCount > 0 && <span>{silenceCount} silence gap{silenceCount !== 1 ? 's' : ''}</span>}
          {causalCount > 0 && <span>{causalCount} causal sequence{causalCount !== 1 ? 's' : ''}</span>}
          <span>avg entropy: {intel.avg_entropy.toFixed(2)}</span>
          {intel.high_entropy_count > 0 && (
            <span>{intel.high_entropy_count} high-entropy entr{intel.high_entropy_count !== 1 ? 'ies' : 'y'}</span>
          )}
        </div>

        <CollapsibleSection title="Spike Detection" count={spikeCount}>
          <SpikeList spikes={intel.spikes} />
        </CollapsibleSection>

        <CollapsibleSection title="Silence Gaps" count={silenceCount}>
          <SilenceGapList gaps={intel.silence_gaps} />
        </CollapsibleSection>

        <CollapsibleSection title="Causal Sequences" count={causalCount}>
          <CausalSequenceList sequences={intel.causal_sequences} />
        </CollapsibleSection>
      </CardContent>
    </Card>
  )
}
