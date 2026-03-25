import { Outlet } from 'react-router-dom'
import { TooltipProvider } from '@/components/ui/tooltip'
import { Navigation } from './Navigation'

export function Layout() {
  return (
    <TooltipProvider>
      <div className="min-h-screen flex flex-col bg-background text-foreground antialiased">
        <Navigation />
        <main className="flex-1 overflow-auto">
          <Outlet />
        </main>
      </div>
    </TooltipProvider>
  )
}
