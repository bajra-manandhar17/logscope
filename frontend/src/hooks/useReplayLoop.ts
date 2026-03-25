import { useEffect, useRef } from 'react'
import { useReplayStore } from '../stores/replayStore'

const TICK_INTERVAL_MS = 66 // ~15fps

export function useReplayLoop() {
  const mode = useReplayStore((s) => s.mode)
  const speed = useReplayStore((s) => s.speed)
  const startTime = useReplayStore((s) => s.startTime)
  const endTime = useReplayStore((s) => s.endTime)
  const { tick, pause } = useReplayStore.getState()

  const rafRef = useRef<number | null>(null)
  const lastFrameTimeRef = useRef<number>(0)
  const lastTickTimeRef = useRef<number>(0)
  const currentMsRef = useRef<number>(0)

  useEffect(() => {
    if (mode !== 'playing' || !startTime || !endTime) {
      if (rafRef.current !== null) {
        cancelAnimationFrame(rafRef.current)
        rafRef.current = null
      }
      return
    }

    const startMs = new Date(startTime).getTime()
    const endMs = new Date(endTime).getTime()

    if (endMs <= startMs) {
      // Guard: complete instantly — endTime is non-null (guarded above)
      tick(endTime!, 1)
      pause()
      return
    }

    // timeScaleMs: how many log-ms advance per real-ms at 1x
    // At 1x, full span completes in 60s → scale = span / 60_000
    const timeScaleMs = (endMs - startMs) / 60_000

    // Initialize currentMs to match current progress (handles resume after seek)
    const currentProgress = useReplayStore.getState().progress
    currentMsRef.current = startMs + currentProgress * (endMs - startMs)
    lastFrameTimeRef.current = performance.now()
    lastTickTimeRef.current = performance.now()

    function frame(now: number) {
      const deltaMs = now - lastFrameTimeRef.current
      lastFrameTimeRef.current = now

      currentMsRef.current = Math.min(
        endMs,
        currentMsRef.current + deltaMs * speed * timeScaleMs,
      )

      const progress = (currentMsRef.current - startMs) / (endMs - startMs)
      const clampedProgress = Math.min(1, progress)

      // Throttle tick calls to ~15fps
      if (now - lastTickTimeRef.current >= TICK_INTERVAL_MS) {
        lastTickTimeRef.current = now
        const currentIso = new Date(currentMsRef.current).toISOString()
        tick(currentIso, clampedProgress)
      }

      if (currentMsRef.current >= endMs) {
        // endTime is non-null (guarded at effect start)
        tick(endTime!, 1)
        pause()
        return
      }

      rafRef.current = requestAnimationFrame(frame)
    }

    rafRef.current = requestAnimationFrame(frame)

    return () => {
      if (rafRef.current !== null) {
        cancelAnimationFrame(rafRef.current)
        rafRef.current = null
      }
    }
  }, [mode, speed, startTime, endTime])
}
