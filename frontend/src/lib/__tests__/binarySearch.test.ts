import { describe, expect, it } from 'vitest'
import { upperBoundByTime, entryIndexAtTime, bucketIndexAtTime } from '../binarySearch'
import type { LogEntry, TimeBucket } from '../../types'

function makeEntry(timestamp: string, line = 1): LogEntry {
  return { timestamp, level: 'info', message: '', source: '', raw: '', line_number: line, entropy: 0 }
}

function makeBucket(timestamp: string): TimeBucket {
  return { timestamp, count: 1, error_count: 0 }
}

describe('upperBoundByTime', () => {
  it('returns 0 for empty array', () => {
    expect(upperBoundByTime([], 1000, (x: number) => x)).toBe(0)
  })

  it('returns length when all elements ≤ target', () => {
    const arr = [100, 200, 300]
    expect(upperBoundByTime(arr, 300, (x) => x)).toBe(3)
  })

  it('returns 0 when all elements > target', () => {
    const arr = [100, 200, 300]
    expect(upperBoundByTime(arr, 50, (x) => x)).toBe(0)
  })

  it('returns correct index at boundary', () => {
    const arr = [100, 200, 300, 400]
    expect(upperBoundByTime(arr, 200, (x) => x)).toBe(2)
    expect(upperBoundByTime(arr, 201, (x) => x)).toBe(2)
    expect(upperBoundByTime(arr, 199, (x) => x)).toBe(1)
  })

  it('handles single element', () => {
    expect(upperBoundByTime([100], 100, (x) => x)).toBe(1)
    expect(upperBoundByTime([100], 99, (x) => x)).toBe(0)
    expect(upperBoundByTime([100], 101, (x) => x)).toBe(1)
  })
})

describe('entryIndexAtTime', () => {
  const entries: LogEntry[] = [
    makeEntry('2024-01-01T00:01:00Z', 1),
    makeEntry('2024-01-01T00:02:00Z', 2),
    makeEntry('2024-01-01T00:03:00Z', 3),
  ]

  it('returns -1 when time is before all entries', () => {
    expect(entryIndexAtTime(entries, '2024-01-01T00:00:00Z')).toBe(-1)
  })

  it('returns last index when time is after all entries', () => {
    expect(entryIndexAtTime(entries, '2024-01-01T00:05:00Z')).toBe(2)
  })

  it('returns correct boundary index', () => {
    expect(entryIndexAtTime(entries, '2024-01-01T00:02:00Z')).toBe(1)
    expect(entryIndexAtTime(entries, '2024-01-01T00:02:30Z')).toBe(1)
  })

  it('returns length - 1 for invalid time', () => {
    expect(entryIndexAtTime(entries, 'not-a-date')).toBe(2)
  })

  it('handles empty array', () => {
    expect(entryIndexAtTime([], '2024-01-01T00:00:00Z')).toBe(-1)
  })

  it('entries without timestamps use -Infinity (treated as always before cutoff)', () => {
    const mixed: LogEntry[] = [
      makeEntry('', 1),
      makeEntry('2024-01-01T00:02:00Z', 2),
    ]
    // entry with no ts has getTime → -Infinity → always ≤ target
    expect(entryIndexAtTime(mixed, '2024-01-01T00:01:00Z')).toBeGreaterThanOrEqual(0)
  })
})

describe('bucketIndexAtTime', () => {
  const buckets: TimeBucket[] = [
    makeBucket('2024-01-01T00:00:00Z'),
    makeBucket('2024-01-01T00:01:00Z'),
    makeBucket('2024-01-01T00:02:00Z'),
  ]

  it('returns -1 when time is before all buckets', () => {
    expect(bucketIndexAtTime(buckets, '2023-12-31T23:59:00Z')).toBe(-1)
  })

  it('returns last index when time is after all buckets', () => {
    expect(bucketIndexAtTime(buckets, '2024-01-01T00:10:00Z')).toBe(2)
  })

  it('returns correct index at boundary', () => {
    expect(bucketIndexAtTime(buckets, '2024-01-01T00:01:00Z')).toBe(1)
    expect(bucketIndexAtTime(buckets, '2024-01-01T00:01:30Z')).toBe(1)
  })

  it('handles empty array', () => {
    expect(bucketIndexAtTime([], '2024-01-01T00:00:00Z')).toBe(-1)
  })

  it('returns length - 1 for invalid time', () => {
    expect(bucketIndexAtTime(buckets, 'bad')).toBe(2)
  })
})
