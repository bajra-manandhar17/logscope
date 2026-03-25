import { useRef } from 'react'
import { Play, Pause, RotateCcw } from 'lucide-react'
import { useReplayStore } from '../../stores/replayStore'
import { Slider } from '@/components/ui/slider'
import type { ReplaySpeed } from '../../types'

const SPEEDS: ReplaySpeed[] = [1, 5, 10, 50]

function formatTimestamp(iso: string | null): string {
  if (!iso) return '--:--:--'
  const d = new Date(iso)
  if (isNaN(d.getTime())) return '--:--:--'
  return d.toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

export function ReplayControlBar() {
  const mode = useReplayStore((s) => s.mode)
  const speed = useReplayStore((s) => s.speed)
  const currentTime = useReplayStore((s) => s.currentTime)
  const progress = useReplayStore((s) => s.progress)
  const startTime = useReplayStore((s) => s.startTime)
  const endTime = useReplayStore((s) => s.endTime)
  const { pause, resume, seek, setSpeed, stop } = useReplayStore.getState()

  const wasPlayingRef = useRef(false)

  if (mode === 'idle') return null

  function handlePlayPause() {
    if (mode === 'playing') {
      pause()
    } else {
      resume()
    }
  }

  function handleSpeedCycle() {
    const idx = SPEEDS.indexOf(speed)
    setSpeed(SPEEDS[(idx + 1) % SPEEDS.length])
  }

  function handleSliderChange(value: number[]) {
    if (!startTime || !endTime) return
    const startMs = new Date(startTime).getTime()
    const endMs = new Date(endTime).getTime()
    const p = value[0] / 1000
    const timeMs = startMs + p * (endMs - startMs)
    seek(new Date(timeMs).toISOString())
  }

  function handlePointerDown() {
    wasPlayingRef.current = mode === 'playing'
    if (mode === 'playing') pause()
  }

  function handlePointerUp() {
    if (wasPlayingRef.current) resume()
  }

  return (
    <div className="bg-card border border-border rounded-lg px-4 py-2 flex items-center gap-3">
      {/* Play/Pause */}
      <button
        onClick={handlePlayPause}
        className="text-foreground hover:text-primary transition-colors shrink-0"
        aria-label={mode === 'playing' ? 'Pause' : 'Play'}
      >
        {mode === 'playing' ? <Pause className="size-4" /> : <Play className="size-4" />}
      </button>

      {/* Speed */}
      <button
        onClick={handleSpeedCycle}
        className="text-xs font-mono text-muted-foreground hover:text-foreground transition-colors shrink-0 w-10 text-left"
      >
        {speed}x
      </button>

      {/* Scrub bar */}
      <div
        className="flex-1"
        onPointerDown={handlePointerDown}
        onPointerUp={handlePointerUp}
      >
        <Slider
          value={[Math.round(progress * 1000)]}
          max={1000}
          step={1}
          onValueChange={handleSliderChange}
          className="w-full"
        />
      </div>

      {/* Timestamp */}
      <span className="text-xs font-mono text-muted-foreground shrink-0 w-20 text-right">
        {formatTimestamp(currentTime)}
      </span>

      {/* Reset */}
      <button
        onClick={stop}
        className="text-muted-foreground hover:text-foreground transition-colors shrink-0"
        aria-label="Reset replay"
      >
        <RotateCcw className="size-4" />
      </button>
    </div>
  )
}
