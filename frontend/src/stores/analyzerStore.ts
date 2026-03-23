import { create } from 'zustand'
import { analyzeFile } from '../api/analyze'
import type { AnalysisResult } from '../types'

export interface LogFilters {
  levels: string[]  // empty = all
  search: string
  timeStart: string // ISO string or ''
  timeEnd: string   // ISO string or ''
}

type Status = 'idle' | 'uploading' | 'done' | 'error'

interface AnalyzerState {
  status: Status
  error: string | null
  result: AnalysisResult | null
  filters: LogFilters
  // actions
  upload: (file: File, format?: string) => Promise<void>
  cancel: () => void
  setFilters: (filters: Partial<LogFilters>) => void
  reset: () => void
}

let abortController: AbortController | null = null

export const useAnalyzerStore = create<AnalyzerState>((set) => ({
  status: 'idle',
  error: null,
  result: null,
  filters: { levels: [], search: '', timeStart: '', timeEnd: '' },

  upload: async (file, format = 'auto') => {
    abortController?.abort()
    abortController = new AbortController()

    set({ status: 'uploading', error: null, result: null })
    try {
      const result = await analyzeFile(file, format, abortController.signal)
      set({ status: 'done', result })
    } catch (err) {
      if (err instanceof DOMException && err.name === 'AbortError') {
        set({ status: 'idle' })
      } else {
        set({ status: 'error', error: err instanceof Error ? err.message : 'Upload failed' })
      }
    } finally {
      abortController = null
    }
  },

  cancel: () => {
    abortController?.abort()
    abortController = null
  },

  setFilters: (partial) =>
    set((s) => ({ filters: { ...s.filters, ...partial } })),

  reset: () => {
    abortController?.abort()
    abortController = null
    set({ status: 'idle', error: null, result: null, filters: { levels: [], search: '', timeStart: '', timeEnd: '' } })
  },
}))
