import type { GenerateConfig } from '../types'

export interface BatchEvent {
  type: 'batch'
  lines: string[]
}

export interface DoneEvent {
  type: 'done'
  totalLines: number
}

export type SSEEvent = BatchEvent | DoneEvent

// Streams SSE events from POST /api/generate.
// Yields events via async generator; respects AbortSignal.
export async function* streamGenerate(
  config: GenerateConfig,
  signal: AbortSignal,
): AsyncGenerator<SSEEvent> {
  const res = await fetch('/api/generate', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(config),
    signal,
  })

  if (!res.ok) {
    const json = await res.json().catch(() => ({}))
    throw new Error(json.error ?? `HTTP ${res.status}`)
  }

  if (!res.body) throw new Error('No response body')

  const reader = res.body.getReader()
  const decoder = new TextDecoder()
  let buf = ''

  try {
    while (true) {
      const { done, value } = await reader.read()
      if (done) break

      buf += decoder.decode(value, { stream: true })
      const events = buf.split('\n\n')
      buf = events.pop() ?? ''

      for (const block of events) {
        const parsed = parseSSEBlock(block)
        if (parsed) yield parsed
      }
    }
  } finally {
    reader.releaseLock()
  }
}

function parseSSEBlock(block: string): SSEEvent | null {
  let eventType = ''
  let dataLine = ''

  for (const line of block.split('\n')) {
    if (line.startsWith('event: ')) eventType = line.slice(7).trim()
    else if (line.startsWith('data: ')) dataLine = line.slice(6).trim()
  }

  if (!eventType || !dataLine) return null

  try {
    const payload = JSON.parse(dataLine)
    if (eventType === 'batch') return { type: 'batch', lines: payload.lines ?? [] }
    if (eventType === 'done') return { type: 'done', totalLines: payload.totalLines ?? 0 }
  } catch {
    // malformed data — skip
  }
  return null
}
