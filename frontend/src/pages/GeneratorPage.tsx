import { ConfigForm } from '../components/generator/ConfigForm'
import { LivePreview } from '../components/generator/LivePreview'
import { ActionBar } from '../components/generator/ActionBar'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'

export function GeneratorPage() {
  return (
    <div className="max-w-7xl mx-auto px-6 py-8 flex flex-col gap-6">
      <div>
        <h1
          className="text-3xl font-bold tracking-tight"
          style={{ fontFamily: 'var(--font-heading)' }}
        >
          Log Generator
        </h1>
        <p className="text-sm text-muted-foreground mt-1">
          Configure and stream synthetic log data for testing and development
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Config panel */}
        <Card>
          <CardHeader>
            <CardTitle className="text-sm">Configuration</CardTitle>
          </CardHeader>
          <CardContent>
            <ConfigForm />
          </CardContent>
        </Card>

        {/* Preview panel */}
        <div className="lg:col-span-2 flex flex-col gap-4">
          <LivePreview />
          <ActionBar />
        </div>
      </div>
    </div>
  )
}
