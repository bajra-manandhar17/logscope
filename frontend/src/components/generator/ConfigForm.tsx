import { useState } from 'react'
import { Loader2, TriangleAlert } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useGeneratorStore } from '../../stores/generatorStore'
import { cn } from '@/lib/utils'

const LEVELS = ['error', 'warn', 'info', 'debug'] as const
const LEVEL_LABEL: Record<string, string> = { error: 'Error', warn: 'Warn', info: 'Info', debug: 'Debug' }
const LEVEL_BADGE: Record<string, { variant: 'destructive' | 'outline' | 'default' | 'secondary'; className?: string }> = {
  error: { variant: 'destructive' },
  warn:  { variant: 'outline', className: 'border-yellow-500/30 text-yellow-500' },
  info:  { variant: 'default' },
  debug: { variant: 'secondary' },
}

function toLocalInput(iso: string): string {
  if (!iso) return ''
  const d = new Date(iso)
  if (isNaN(d.getTime())) return ''
  return d.toISOString().slice(0, 16)
}

function fromLocalInput(val: string): string {
  if (!val) return ''
  return new Date(val).toISOString()
}

function validate(config: ReturnType<typeof useGeneratorStore.getState>['config']): string[] {
  const errors: string[] = []
  if (config.total_lines < 1 || config.total_lines > 1_000_000) {
    errors.push('Total lines must be between 1 and 1,000,000')
  }
  const sum = Object.values(config.levels).reduce((a, b) => a + b, 0)
  if (Math.abs(sum - 1.0) > 0.01) {
    errors.push(`Level weights must sum to 1.0 (currently ${sum.toFixed(2)})`)
  }
  if (!config.start || !config.end) {
    errors.push('Start and end time are required')
  } else if (new Date(config.end) <= new Date(config.start)) {
    errors.push('End time must be after start time')
  }
  return errors
}

export function ConfigForm() {
  const config = useGeneratorStore((s) => s.config)
  const setConfig = useGeneratorStore((s) => s.setConfig)
  const generate = useGeneratorStore((s) => s.generate)
  const status = useGeneratorStore((s) => s.status)
  const [submitted, setSubmitted] = useState(false)

  const errors = validate(config)
  const levelSum = Object.values(config.levels).reduce((a, b) => a + b, 0)
  const generating = status === 'generating'

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setSubmitted(true)
    if (errors.length > 0) return
    generate()
  }

  function setLevel(level: string, value: number) {
    setConfig({ levels: { ...config.levels, [level]: value } })
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-5">
      {/* Format */}
      <div className="flex flex-col gap-1.5">
        <label className="text-sm font-medium">Format</label>
        <Tabs value={config.format} onValueChange={(v) => setConfig({ format: v as 'json' | 'plaintext' })}>
          <TabsList>
            <TabsTrigger value="json">JSON</TabsTrigger>
            <TabsTrigger value="plaintext">Plaintext</TabsTrigger>
          </TabsList>
        </Tabs>
      </div>

      {/* Total lines */}
      <div className="flex flex-col gap-1.5">
        <label htmlFor="total-lines" className="text-sm font-medium">
          Total Lines
        </label>
        <Input
          id="total-lines"
          type="number"
          min={1}
          max={1_000_000}
          value={config.total_lines}
          onChange={(e) => setConfig({ total_lines: Math.max(1, parseInt(e.target.value, 10) || 1) })}
          className="w-48 font-mono"
        />
      </div>

      {/* Level distribution */}
      <div className="flex flex-col gap-2">
        <div className="flex items-center justify-between">
          <span className="text-sm font-medium">Level Distribution</span>
          <span className={cn('text-xs tabular-nums', Math.abs(levelSum - 1) > 0.01 ? 'text-destructive' : 'text-muted-foreground')}>
            sum: {levelSum.toFixed(2)}
          </span>
        </div>
        <div className="grid grid-cols-2 gap-2">
          {LEVELS.map((level) => {
            const cfg = LEVEL_BADGE[level]
            return (
              <div key={level} className="flex items-center gap-2">
                <Badge variant={cfg.variant} className={cn('w-14 justify-center', cfg.className)}>
                  {LEVEL_LABEL[level]}
                </Badge>
                <Input
                  type="number"
                  min={0}
                  max={1}
                  step={0.01}
                  value={config.levels[level] ?? 0}
                  onChange={(e) => setLevel(level, parseFloat(e.target.value) || 0)}
                  className="w-24 tabular-nums font-mono"
                />
              </div>
            )
          })}
        </div>
      </div>

      {/* Time range */}
      <div className="flex flex-col gap-2">
        <span className="text-sm font-medium">Time Range</span>
        <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:gap-3">
          <div className="flex flex-col gap-1">
            <label className="text-xs text-muted-foreground">Start</label>
            <Input
              type="datetime-local"
              value={toLocalInput(config.start)}
              onChange={(e) => setConfig({ start: fromLocalInput(e.target.value) })}
            />
          </div>
          <div className="flex flex-col gap-1">
            <label className="text-xs text-muted-foreground">End</label>
            <Input
              type="datetime-local"
              value={toLocalInput(config.end)}
              onChange={(e) => setConfig({ end: fromLocalInput(e.target.value) })}
            />
          </div>
        </div>
      </div>

      {/* Validation errors */}
      {submitted && errors.length > 0 && (
        <div className="flex items-start gap-2 rounded-lg border border-destructive/20 bg-destructive/10 p-3">
          <TriangleAlert className="size-4 text-destructive shrink-0 mt-0.5" />
          <ul role="alert" className="flex flex-col gap-1">
            {errors.map((e) => (
              <li key={e} className="text-sm text-destructive">{e}</li>
            ))}
          </ul>
        </div>
      )}

      <Button type="submit" disabled={generating} className="self-start">
        {generating ? (
          <>
            <Loader2 className="animate-spin" />
            Generating…
          </>
        ) : (
          'Generate'
        )}
      </Button>
    </form>
  )
}
