import { describe, expect, it, vi } from 'vitest'
import { streamGenerate } from './generate'
import type { GenerateConfig } from '../types'

const cfg: GenerateConfig = {
  format: 'json',
  total_lines: 10,
  levels: { error: 0.05, warn: 0.15, info: 0.70, debug: 0.10 },
  start: '2024-01-01T00:00:00Z',
  end: '2024-01-02T00:00:00Z',
}

function makeStream(chunks: string[]): ReadableStream<Uint8Array> {
  const encoder = new TextEncoder()
  return new ReadableStream({
    start(controller) {
      for (const chunk of chunks) {
        controller.enqueue(encoder.encode(chunk))
      }
      controller.close()
    },
  })
}

describe('streamGenerate', () => {
  it('yields batch and done events', async () => {
    const body = makeStream([
      'event: batch\ndata: {"lines":["line1","line2"]}\n\n',
      'event: done\ndata: {"totalLines":2}\n\n',
    ])

    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: true, body }))

    const events = []
    for await (const e of streamGenerate(cfg, new AbortController().signal)) {
      events.push(e)
    }

    expect(events).toHaveLength(2)
    expect(events[0]).toEqual({ type: 'batch', lines: ['line1', 'line2'] })
    expect(events[1]).toEqual({ type: 'done', totalLines: 2 })

    vi.unstubAllGlobals()
  })

  it('handles chunks split across SSE boundaries', async () => {
    // Simulate chunk arriving mid-event-block
    const body = makeStream([
      'event: batch\ndata: {"lines":["a"]}\n',
      '\nevent: done\ndata: {"totalLines":1}\n\n',
    ])

    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: true, body }))

    const events = []
    for await (const e of streamGenerate(cfg, new AbortController().signal)) {
      events.push(e)
    }

    expect(events).toHaveLength(2)
    expect(events[0].type).toBe('batch')
    expect(events[1].type).toBe('done')

    vi.unstubAllGlobals()
  })

  it('throws on non-ok response', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: false,
      json: async () => ({ error: 'invalid_config' }),
    }))

    const gen = streamGenerate(cfg, new AbortController().signal)
    await expect(gen.next()).rejects.toThrow('invalid_config')

    vi.unstubAllGlobals()
  })

  it('stops on abort', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(
      new DOMException('aborted', 'AbortError'),
    ))

    const controller = new AbortController()
    controller.abort()

    const gen = streamGenerate(cfg, controller.signal)
    await expect(gen.next()).rejects.toThrow()

    vi.unstubAllGlobals()
  })
})
