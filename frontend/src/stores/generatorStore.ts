import { create } from 'zustand'
import { streamGenerate } from '../api/generate'
import { useAnalyzerStore } from './analyzerStore'
import type { GenerateConfig } from '../types'

type Status = 'idle' | 'generating' | 'done' | 'error'

const DEFAULT_CONFIG: GenerateConfig = {
  format: 'json',
  total_lines: 1000,
  levels: { error: 0.05, warn: 0.15, info: 0.70, debug: 0.10 },
  start: new Date(Date.now() - 3600 * 1000).toISOString(),
  end: new Date().toISOString(),
}

interface GeneratorState {
  config: GenerateConfig
  status: Status
  lines: string[]
  error: string | null
  // actions
  setConfig: (patch: Partial<GenerateConfig>) => void
  generate: () => Promise<void>
  abort: () => void
  sendToAnalyzer: (navigate: (path: string) => void) => void
  reset: () => void
}

let abortController: AbortController | null = null

export const useGeneratorStore = create<GeneratorState>((set, get) => ({
  config: DEFAULT_CONFIG,
  status: 'idle',
  lines: [],
  error: null,

  setConfig: (patch) =>
    set((s) => ({ config: { ...s.config, ...patch } })),

  generate: async () => {
    abortController?.abort()
    abortController = new AbortController()

    set({ status: 'generating', lines: [], error: null })

    // Map frontend snake_case config to what the Go handler expects (Go JSON uses exported field names).
    const cfg = get().config
    const body = {
      Format: cfg.format,
      TotalLines: cfg.total_lines,
      Levels: Object.fromEntries(
        Object.entries(cfg.levels).map(([k, v]) => [k, v]),
      ),
      Start: cfg.start,
      End: cfg.end,
    }

    try {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      for await (const event of streamGenerate(body as any, abortController.signal)) {
        if (event.type === 'batch') {
          set((s) => ({ lines: [...s.lines, ...event.lines] }))
        } else if (event.type === 'done') {
          set({ status: 'done' })
        }
      }
    } catch (err) {
      if (err instanceof DOMException && err.name === 'AbortError') {
        set({ status: 'idle' })
      } else {
        set({ status: 'error', error: err instanceof Error ? err.message : 'Generation failed' })
      }
    } finally {
      abortController = null
    }
  },

  abort: () => {
    abortController?.abort()
    abortController = null
  },

  sendToAnalyzer: (navigate) => {
    const { lines, config } = get()
    if (!lines.length) return

    const content = lines.join('\n')
    const file = new File([content], 'generated.log', { type: 'text/plain' })
    useAnalyzerStore.getState().upload(file, config.format)
    navigate('/analyze')
  },

  reset: () => {
    abortController?.abort()
    abortController = null
    set({ status: 'idle', lines: [], error: null, config: DEFAULT_CONFIG })
  },
}))
