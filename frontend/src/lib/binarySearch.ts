import type { LogEntry, TimeBucket } from '../types'

/**
 * Returns the index of the first element whose time is strictly greater than targetMs.
 * Equivalent to std::upper_bound. Returns arr.length if all elements ≤ targetMs.
 */
export function upperBoundByTime<T>(
  arr: T[],
  targetMs: number,
  getTimeMs: (item: T) => number,
): number {
  let lo = 0
  let hi = arr.length
  while (lo < hi) {
    const mid = (lo + hi) >>> 1
    if (getTimeMs(arr[mid]) <= targetMs) {
      lo = mid + 1
    } else {
      hi = mid
    }
  }
  return lo
}

/**
 * Returns the index of the last entry whose timestamp ≤ time (exclusive upper bound - 1).
 * Returns -1 if no entries have timestamps or all are after time.
 * Entries without timestamps are treated as always-included by consumers.
 */
export function entryIndexAtTime(entries: LogEntry[], time: string): number {
  const targetMs = new Date(time).getTime()
  if (isNaN(targetMs)) return entries.length - 1

  // Find first entry with a valid timestamp to determine if sorted search applies
  const idx = upperBoundByTime(
    entries,
    targetMs,
    (e) => (e.timestamp ? new Date(e.timestamp).getTime() : -Infinity),
  )
  return idx - 1
}

/**
 * Returns the index of the last bucket whose timestamp ≤ time.
 * Returns -1 if no buckets or all are after time.
 */
export function bucketIndexAtTime(buckets: TimeBucket[], time: string): number {
  const targetMs = new Date(time).getTime()
  if (isNaN(targetMs)) return buckets.length - 1

  const idx = upperBoundByTime(
    buckets,
    targetMs,
    (b) => new Date(b.timestamp).getTime(),
  )
  return idx - 1
}
