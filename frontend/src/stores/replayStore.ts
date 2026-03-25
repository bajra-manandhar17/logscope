import { create } from 'zustand'
import type { ReplayMode, ReplaySpeed } from '../types'
import { useAnalyzerStore } from './analyzerStore'

interface ReplayState {
  mode: ReplayMode
  speed: ReplaySpeed
  currentTime: string | null
  startTime: string | null
  endTime: string | null
  progress: number // 0–1
  // actions
  start: (startTime: string, endTime: string) => void
  pause: () => void
  resume: () => void
  seek: (time: string) => void
  setSpeed: (speed: ReplaySpeed) => void
  tick: (time: string, progress: number) => void
  stop: () => void
}

export const useReplayStore = create<ReplayState>((set, get) => ({
  mode: 'idle',
  speed: 1,
  currentTime: null,
  startTime: null,
  endTime: null,
  progress: 0,

  start: (startTime, endTime) => {
    set({
      mode: 'playing',
      startTime,
      endTime,
      currentTime: startTime,
      progress: 0,
    })
  },

  pause: () => {
    if (get().mode === 'playing') set({ mode: 'paused' })
  },

  resume: () => {
    if (get().mode === 'paused') set({ mode: 'playing' })
  },

  seek: (time) => {
    const { startTime, endTime } = get()
    if (!startTime || !endTime) return
    const startMs = new Date(startTime).getTime()
    const endMs = new Date(endTime).getTime()
    const timeMs = new Date(time).getTime()
    const span = endMs - startMs
    const progress = span > 0 ? Math.min(1, Math.max(0, (timeMs - startMs) / span)) : 0
    set({ currentTime: time, progress })
  },

  setSpeed: (speed) => set({ speed }),

  tick: (time, progress) => set({ currentTime: time, progress }),

  stop: () =>
    set({
      mode: 'idle',
      currentTime: null,
      startTime: null,
      endTime: null,
      progress: 0,
    }),
}))

// Auto-stop when analyzer result is cleared
useAnalyzerStore.subscribe((state) => {
  if (state.result === null) {
    useReplayStore.getState().stop()
  }
})
