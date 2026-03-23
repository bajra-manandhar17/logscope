import { useRef, useState } from 'react'
import { Upload } from 'lucide-react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { useAnalyzerStore } from '../../stores/analyzerStore'

const MAX_SIZE_BYTES = 100 * 1024 * 1024 // 100 MB
const ACCEPTED_TYPES = ['.log', '.txt', '.json', 'text/plain', 'application/json']

function validateFile(file: File): string | null {
  if (file.size > MAX_SIZE_BYTES) return 'File exceeds 100 MB limit'
  const ext = file.name.slice(file.name.lastIndexOf('.')).toLowerCase()
  const ok = ACCEPTED_TYPES.some((t) => t === ext || file.type === t)
  if (!ok) return 'Unsupported file type — use .log, .txt, or .json'
  return null
}

export function FileUpload() {
  const upload = useAnalyzerStore((s) => s.upload)
  const status = useAnalyzerStore((s) => s.status)
  const error = useAnalyzerStore((s) => s.error)
  const cancel = useAnalyzerStore((s) => s.cancel)

  const inputRef = useRef<HTMLInputElement>(null)
  const [dragOver, setDragOver] = useState(false)
  const [validationError, setValidationError] = useState<string | null>(null)

  const uploading = status === 'uploading'

  function handleFile(file: File) {
    const err = validateFile(file)
    if (err) { setValidationError(err); return }
    setValidationError(null)
    upload(file)
  }

  function onDrop(e: React.DragEvent) {
    e.preventDefault()
    setDragOver(false)
    const file = e.dataTransfer.files[0]
    if (file) handleFile(file)
  }

  function onInputChange(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0]
    if (file) handleFile(file)
    e.target.value = ''
  }

  const displayError = validationError ?? error

  return (
    <div className="flex flex-col items-center gap-4 w-full max-w-lg">
      <div
        role="button"
        tabIndex={0}
        aria-label="Drop log file here or click to browse"
        onClick={() => !uploading && inputRef.current?.click()}
        onKeyDown={(e) => e.key === 'Enter' && !uploading && inputRef.current?.click()}
        onDragOver={(e) => { e.preventDefault(); setDragOver(true) }}
        onDragLeave={() => setDragOver(false)}
        onDrop={onDrop}
        className={cn(
          'relative w-full rounded-xl border p-10 text-center cursor-pointer transition-all duration-150 select-none overflow-hidden',
          dragOver
            ? 'border-primary bg-primary/8 ring-2 ring-primary/20'
            : 'border-border/50 bg-muted/20 hover:border-primary/40 hover:bg-primary/5',
          uploading && 'pointer-events-none opacity-60',
        )}
      >
        {uploading && (
          <div className="absolute top-0 left-0 right-0 h-0.5 bg-muted/30 overflow-hidden">
            <div className="h-full w-1/3 bg-primary rounded-full animate-[shimmer_1.5s_ease-in-out_infinite]" />
          </div>
        )}
        <div className="mx-auto mb-3 bg-muted/50 p-3 rounded-full w-fit">
          <Upload className="size-6 text-muted-foreground" />
        </div>
        <p className="text-sm font-medium">
          {uploading ? 'Analyzing…' : 'Drop a log file here'}
        </p>
        <p className="mt-1 text-xs text-muted-foreground">
          {uploading ? 'Processing your file' : '.log · .txt · .json — up to 100 MB'}
        </p>
      </div>

      <input
        ref={inputRef}
        type="file"
        accept=".log,.txt,.json,text/plain,application/json"
        className="hidden"
        onChange={onInputChange}
      />

      {uploading ? (
        <Button variant="outline" size="sm" onClick={cancel}>
          Cancel
        </Button>
      ) : (
        <Button variant="outline" size="sm" onClick={() => inputRef.current?.click()}>
          Browse files
        </Button>
      )}

      {displayError && (
        <p role="alert" className="text-sm text-destructive text-center">
          {displayError}
        </p>
      )}
    </div>
  )
}
