import { beforeEach, describe, expect, it } from 'vitest'
import { useReplayStore } from '../replayStore'

const START = '2024-01-01T00:00:00Z'
const END = '2024-01-01T01:00:00Z'
const MID = '2024-01-01T00:30:00Z'

beforeEach(() => {
  useReplayStore.getState().stop()
})

describe('replayStore', () => {
  it('starts idle', () => {
    const s = useReplayStore.getState()
    expect(s.mode).toBe('idle')
    expect(s.currentTime).toBeNull()
    expect(s.progress).toBe(0)
  })

  it('start sets mode to playing and initializes times', () => {
    useReplayStore.getState().start(START, END)
    const s = useReplayStore.getState()
    expect(s.mode).toBe('playing')
    expect(s.startTime).toBe(START)
    expect(s.endTime).toBe(END)
    expect(s.currentTime).toBe(START)
    expect(s.progress).toBe(0)
  })

  it('pause transitions playing → paused', () => {
    useReplayStore.getState().start(START, END)
    useReplayStore.getState().pause()
    expect(useReplayStore.getState().mode).toBe('paused')
  })

  it('pause does nothing when idle', () => {
    useReplayStore.getState().pause()
    expect(useReplayStore.getState().mode).toBe('idle')
  })

  it('resume transitions paused → playing', () => {
    useReplayStore.getState().start(START, END)
    useReplayStore.getState().pause()
    useReplayStore.getState().resume()
    expect(useReplayStore.getState().mode).toBe('playing')
  })

  it('resume does nothing when idle', () => {
    useReplayStore.getState().resume()
    expect(useReplayStore.getState().mode).toBe('idle')
  })

  it('seek updates currentTime and computes progress', () => {
    useReplayStore.getState().start(START, END)
    useReplayStore.getState().seek(MID)
    const s = useReplayStore.getState()
    expect(s.currentTime).toBe(MID)
    expect(s.progress).toBeCloseTo(0.5, 2)
  })

  it('seek clamps progress to [0, 1]', () => {
    useReplayStore.getState().start(START, END)
    useReplayStore.getState().seek('2020-01-01T00:00:00Z') // way before start
    expect(useReplayStore.getState().progress).toBe(0)
    useReplayStore.getState().seek('2030-01-01T00:00:00Z') // way after end
    expect(useReplayStore.getState().progress).toBe(1)
  })

  it('tick updates currentTime and progress', () => {
    useReplayStore.getState().start(START, END)
    useReplayStore.getState().tick(MID, 0.5)
    const s = useReplayStore.getState()
    expect(s.currentTime).toBe(MID)
    expect(s.progress).toBe(0.5)
  })

  it('setSpeed updates speed', () => {
    useReplayStore.getState().setSpeed(10)
    expect(useReplayStore.getState().speed).toBe(10)
  })

  it('stop resets all state to idle', () => {
    useReplayStore.getState().start(START, END)
    useReplayStore.getState().tick(MID, 0.5)
    useReplayStore.getState().stop()
    const s = useReplayStore.getState()
    expect(s.mode).toBe('idle')
    expect(s.currentTime).toBeNull()
    expect(s.startTime).toBeNull()
    expect(s.endTime).toBeNull()
    expect(s.progress).toBe(0)
  })
})
