import { beforeEach, describe, expect, it, vi } from 'vitest'
import { useAnalyzerStore } from './analyzerStore'
import type { AnalysisResult } from '../types'

const mockResult: AnalysisResult = {
  format_detected: 'json',
  summary: {
    total_lines: 10,
    error_count: 1,
    warn_count: 2,
    info_count: 6,
    debug_count: 1,
    time_range: ['2024-01-01T00:00:00Z', '2024-01-01T01:00:00Z'],
    top_sources: ['api'],
  },
  entries: [],
  patterns: [],
  time_series: [],
  bucket_interval: '1m',
}

beforeEach(() => {
  useAnalyzerStore.getState().reset()
})

describe('analyzerStore', () => {
  it('starts idle', () => {
    const s = useAnalyzerStore.getState()
    expect(s.status).toBe('idle')
    expect(s.result).toBeNull()
    expect(s.error).toBeNull()
  })

  it('transitions idle → uploading → done on success', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: async () => mockResult,
    }))

    const states: string[] = []
    const unsub = useAnalyzerStore.subscribe((s) => states.push(s.status))

    await useAnalyzerStore.getState().upload(new File([''], 'test.log'))

    unsub()
    expect(states).toContain('uploading')
    expect(states).toContain('done')
    expect(useAnalyzerStore.getState().result).toEqual(mockResult)

    vi.unstubAllGlobals()
  })

  it('transitions idle → uploading → error on fetch failure', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: false,
      json: async () => ({ error: 'bad file' }),
    }))

    await useAnalyzerStore.getState().upload(new File([''], 'test.log'))

    const s = useAnalyzerStore.getState()
    expect(s.status).toBe('error')
    expect(s.error).toBe('bad file')

    vi.unstubAllGlobals()
  })

  it('transitions to idle on abort', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(
      Object.assign(new DOMException('aborted', 'AbortError')),
    ))

    await useAnalyzerStore.getState().upload(new File([''], 'test.log'))

    expect(useAnalyzerStore.getState().status).toBe('idle')
    vi.unstubAllGlobals()
  })

  it('cancel aborts and resets to idle', async () => {
    let rejectFn!: (e: unknown) => void
    vi.stubGlobal('fetch', vi.fn().mockReturnValue(
      new Promise((_, r) => { rejectFn = r }),
    ))

    const uploadPromise = useAnalyzerStore.getState().upload(new File([''], 'test.log'))
    expect(useAnalyzerStore.getState().status).toBe('uploading')

    rejectFn(Object.assign(new DOMException('aborted', 'AbortError')))
    await uploadPromise

    expect(useAnalyzerStore.getState().status).toBe('idle')
    vi.unstubAllGlobals()
  })

  it('setFilters merges partial filters', () => {
    useAnalyzerStore.getState().setFilters({ search: 'error' })
    expect(useAnalyzerStore.getState().filters.search).toBe('error')
    expect(useAnalyzerStore.getState().filters.levels).toEqual([])
  })

  it('reset clears all state', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: async () => mockResult,
    }))
    await useAnalyzerStore.getState().upload(new File([''], 'test.log'))

    useAnalyzerStore.getState().reset()
    const s = useAnalyzerStore.getState()
    expect(s.status).toBe('idle')
    expect(s.result).toBeNull()
    expect(s.filters.search).toBe('')

    vi.unstubAllGlobals()
  })
})
