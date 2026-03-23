import { beforeEach, describe, expect, it, vi } from 'vitest'
import { useGeneratorStore } from './generatorStore'

function makeStream(chunks: string[]): ReadableStream<Uint8Array> {
  const encoder = new TextEncoder()
  return new ReadableStream({
    start(controller) {
      for (const chunk of chunks) controller.enqueue(encoder.encode(chunk))
      controller.close()
    },
  })
}

beforeEach(() => {
  useGeneratorStore.getState().reset()
  vi.unstubAllGlobals()
})

describe('generatorStore', () => {
  it('starts idle with default config', () => {
    const s = useGeneratorStore.getState()
    expect(s.status).toBe('idle')
    expect(s.lines).toEqual([])
    expect(s.error).toBeNull()
    expect(s.config.format).toBe('json')
  })

  it('accumulates lines from batch events', async () => {
    const body = makeStream([
      'event: batch\ndata: {"lines":["a","b"]}\n\n',
      'event: batch\ndata: {"lines":["c"]}\n\n',
      'event: done\ndata: {"totalLines":3}\n\n',
    ])
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: true, body }))

    await useGeneratorStore.getState().generate()

    const s = useGeneratorStore.getState()
    expect(s.status).toBe('done')
    expect(s.lines).toEqual(['a', 'b', 'c'])
  })

  it('transitions idle → generating → done', async () => {
    const body = makeStream(['event: done\ndata: {"totalLines":0}\n\n'])
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: true, body }))

    const statuses: string[] = []
    const unsub = useGeneratorStore.subscribe((s) => statuses.push(s.status))
    await useGeneratorStore.getState().generate()
    unsub()

    expect(statuses).toContain('generating')
    expect(statuses).toContain('done')
  })

  it('transitions to error on fetch failure', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: false,
      json: async () => ({ error: 'invalid_config' }),
    }))

    await useGeneratorStore.getState().generate()

    const s = useGeneratorStore.getState()
    expect(s.status).toBe('error')
    expect(s.error).toBe('invalid_config')
  })

  it('transitions to idle on abort', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(
      new DOMException('aborted', 'AbortError'),
    ))

    await useGeneratorStore.getState().generate()
    expect(useGeneratorStore.getState().status).toBe('idle')
  })

  it('setConfig merges partial config', () => {
    useGeneratorStore.getState().setConfig({ format: 'plaintext' })
    expect(useGeneratorStore.getState().config.format).toBe('plaintext')
    expect(useGeneratorStore.getState().config.total_lines).toBe(1000)
  })

  it('reset clears lines and status', async () => {
    const body = makeStream([
      'event: batch\ndata: {"lines":["x"]}\n\n',
      'event: done\ndata: {"totalLines":1}\n\n',
    ])
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: true, body }))
    await useGeneratorStore.getState().generate()

    useGeneratorStore.getState().reset()
    const s = useGeneratorStore.getState()
    expect(s.status).toBe('idle')
    expect(s.lines).toEqual([])
  })
})
