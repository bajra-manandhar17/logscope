import type { AnalysisResult } from '../types'

export async function analyzeFile(
  file: File,
  format: string = 'auto',
  signal?: AbortSignal,
): Promise<AnalysisResult> {
  const form = new FormData()
  form.append('file', file)

  const url = format && format !== 'auto' ? `/api/analyze?format=${format}` : '/api/analyze'
  const res = await fetch(url, { method: 'POST', body: form, signal })

  const json = await res.json()
  if (!res.ok) {
    throw new Error(json.error ?? `HTTP ${res.status}`)
  }
  return json as AnalysisResult
}
