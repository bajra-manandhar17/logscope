import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Download, Clipboard, ArrowRight, Check } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { Tooltip, TooltipTrigger, TooltipContent } from '@/components/ui/tooltip'
import { useGeneratorStore } from '../../stores/generatorStore'
import { useAnalyzerStore } from '../../stores/analyzerStore'

export function ActionBar() {
  const lines = useGeneratorStore((s) => s.lines)
  const status = useGeneratorStore((s) => s.status)
  const config = useGeneratorStore((s) => s.config)
  const analyzerStatus = useAnalyzerStore((s) => s.status)
  const navigate = useNavigate()

  const [copied, setCopied] = useState(false)

  const disabled = status !== 'done' || lines.length === 0
  const sending = analyzerStatus === 'uploading'

  function download() {
    const text = lines.join('\n')
    const ext = config.format === 'json' ? 'json' : 'log'
    const blob = new Blob([text], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `generated.${ext}`
    a.click()
    URL.revokeObjectURL(url)
  }

  async function copy() {
    await navigator.clipboard.writeText(lines.join('\n'))
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  function sendToAnalyzer() {
    useGeneratorStore.getState().sendToAnalyzer(navigate)
  }

  if (lines.length === 0) return null

  return (
    <>
      <Separator />
      <div className="flex flex-wrap items-center gap-3">
        <Tooltip>
          <TooltipTrigger asChild>
            <Button variant="outline" size="sm" onClick={download} disabled={disabled}>
              <Download />
              Download
            </Button>
          </TooltipTrigger>
          <TooltipContent>Download as .{config.format === 'json' ? 'json' : 'log'} file</TooltipContent>
        </Tooltip>

        <Tooltip>
          <TooltipTrigger asChild>
            <Button variant="outline" size="sm" onClick={copy} disabled={disabled}>
              {copied ? <Check className="text-green-500" /> : <Clipboard />}
              {copied ? 'Copied!' : 'Copy'}
            </Button>
          </TooltipTrigger>
          <TooltipContent>Copy all lines to clipboard</TooltipContent>
        </Tooltip>

        <Tooltip>
          <TooltipTrigger asChild>
            <Button size="sm" onClick={sendToAnalyzer} disabled={disabled || sending} className="shadow-sm shadow-primary/20">
              {sending ? (
                <>
                  <span className="size-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
                  Sending…
                </>
              ) : (
                <>
                  <ArrowRight />
                  Send to Analyzer
                </>
              )}
            </Button>
          </TooltipTrigger>
          <TooltipContent>Send generated logs to the Analyzer</TooltipContent>
        </Tooltip>
      </div>
    </>
  )
}
