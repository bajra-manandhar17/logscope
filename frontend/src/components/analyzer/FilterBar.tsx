import { Search, X } from 'lucide-react'
import { useAnalyzerStore } from '../../stores/analyzerStore'
import { Input } from '@/components/ui/input'
import { Checkbox } from '@/components/ui/checkbox'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'

const LEVELS = ['error', 'warn', 'info', 'debug'] as const

const levelColor: Record<string, string> = {
  error: 'text-red-500',
  warn:  'text-yellow-500',
  info:  'text-blue-500',
  debug: 'text-muted-foreground',
}

export function FilterBar() {
  const filters = useAnalyzerStore((s) => s.filters)
  const setFilters = useAnalyzerStore((s) => s.setFilters)

  function toggleLevel(level: string) {
    const next = filters.levels.includes(level)
      ? filters.levels.filter((l) => l !== level)
      : [...filters.levels, level]
    setFilters({ levels: next })
  }

  const hasFilters = filters.search || filters.levels.length > 0 || filters.timeStart || filters.timeEnd

  return (
    <div className="flex flex-wrap items-center gap-3 rounded-lg border border-border bg-card/50 backdrop-blur-sm px-3 py-2 text-sm">
      {/* Search */}
      <div className="relative">
        <Search className="absolute left-2 top-1/2 -translate-y-1/2 size-3.5 text-muted-foreground" />
        <Input
          type="search"
          placeholder="Search messages…"
          value={filters.search}
          onChange={(e) => setFilters({ search: e.target.value })}
          className="h-7 min-w-48 pl-7 text-sm"
        />
      </div>

      <Separator orientation="vertical" className="h-5" />

      {/* Level checkboxes */}
      <div className="flex items-center gap-3">
        {LEVELS.map((level) => (
          <label key={level} className="flex items-center gap-1.5 cursor-pointer select-none">
            <Checkbox
              checked={filters.levels.includes(level)}
              onCheckedChange={() => toggleLevel(level)}
            />
            <span className={`capitalize text-sm font-medium ${levelColor[level]}`}>{level}</span>
          </label>
        ))}
      </div>

      <Separator orientation="vertical" className="h-5" />

      {/* Time range */}
      <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
        <span>From</span>
        <Input
          type="datetime-local"
          value={filters.timeStart}
          onChange={(e) => setFilters({ timeStart: e.target.value })}
          className="h-7 text-xs w-auto"
        />
        <span>To</span>
        <Input
          type="datetime-local"
          value={filters.timeEnd}
          onChange={(e) => setFilters({ timeEnd: e.target.value })}
          className="h-7 text-xs w-auto"
        />
      </div>

      {/* Clear */}
      {hasFilters && (
        <Button
          variant="ghost"
          size="xs"
          onClick={() => setFilters({ levels: [], search: '', timeStart: '', timeEnd: '' })}
          className="ml-auto"
        >
          <X />
          Clear
        </Button>
      )}
    </div>
  )
}
