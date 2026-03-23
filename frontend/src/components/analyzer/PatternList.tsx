import { useState } from 'react'
import { useAnalyzerStore } from '../../stores/analyzerStore'
import type { Pattern } from '../../types'
import { Card, CardHeader, CardTitle, CardAction, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'

const placeholderRe = /(\{UUID\}|\{IP\}|\{HEX\}|\{NUM\})/g

const placeholderColor: Record<string, string> = {
  '{UUID}': 'bg-purple-500/15 text-purple-400',
  '{IP}':   'bg-blue-500/15 text-blue-400',
  '{HEX}':  'bg-orange-500/15 text-orange-400',
  '{NUM}':  'bg-green-500/15 text-green-400',
}

function TemplateText({ template }: { template: string }) {
  const parts = template.split(placeholderRe)
  return (
    <span className="font-mono text-xs break-all">
      {parts.map((part, i) =>
        placeholderRe.test(part) ? (
          <span
            key={i}
            className={`inline-block rounded px-1 py-0 mx-0.5 text-[10px] font-semibold ${placeholderColor[part] ?? 'bg-muted text-muted-foreground'}`}
          >
            {part}
          </span>
        ) : (
          <span key={i}>{part}</span>
        ),
      )}
    </span>
  )
}

function PatternRow({ pattern, rank }: { pattern: Pattern; rank: number }) {
  const [expanded, setExpanded] = useState(false)

  return (
    <div className="border-b border-border/60 last:border-0 px-4 py-3">
      <div className="flex items-start gap-3">
        <span className="mt-0.5 min-w-[2rem] text-right text-xs text-muted-foreground tabular-nums">
          #{rank}
        </span>
        <div className="flex-1 min-w-0">
          <TemplateText template={pattern.template} />
        </div>
        <div className="flex items-center gap-2 shrink-0">
          <Badge variant="outline" className="tabular-nums">
            {pattern.count.toLocaleString()}
          </Badge>
          <span className="text-xs text-muted-foreground">hits</span>
          <Button
            variant="ghost"
            size="xs"
            onClick={() => setExpanded((v) => !v)}
            aria-expanded={expanded}
          >
            {expanded ? 'hide' : 'sample'}
          </Button>
        </div>
      </div>
      {expanded && (
        <div className="mt-2 ml-[3.25rem] rounded-md border-l-2 border-primary/30 bg-muted/30 px-3 py-2 font-mono text-xs text-muted-foreground break-all">
          {pattern.sample_line}
        </div>
      )}
    </div>
  )
}

export function PatternList() {
  const result = useAnalyzerStore((s) => s.result)

  if (!result) return null

  const { patterns } = result

  if (!patterns?.length) {
    return (
      <div className="flex items-center justify-center h-24 rounded-xl border border-border text-sm text-muted-foreground">
        No patterns detected
      </div>
    )
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-sm">Detected Patterns</CardTitle>
        <CardAction>
          <span className="text-xs text-muted-foreground">{patterns.length.toLocaleString()} templates</span>
        </CardAction>
      </CardHeader>
      <CardContent className="p-0">
        <div className="overflow-y-auto max-h-96">
          {patterns.map((p, i) => (
            <PatternRow key={p.template} pattern={p} rank={i + 1} />
          ))}
        </div>
      </CardContent>
    </Card>
  )
}
