import { useEffect, useState } from 'react'
import { NavLink } from 'react-router-dom'
import { Sun, Moon } from 'lucide-react'
import { Separator } from '@/components/ui/separator'

function useTheme() {
  const [dark, setDark] = useState(() => {
    const stored = localStorage.getItem('theme')
    return stored ? stored === 'dark' : true
  })

  useEffect(() => {
    document.documentElement.classList.toggle('dark', dark)
    localStorage.setItem('theme', dark ? 'dark' : 'light')
  }, [dark])

  return { dark, toggle: () => setDark((v) => !v) }
}

export function Navigation() {
  const { dark, toggle } = useTheme()

  return (
    <nav className="flex items-center gap-1 border-b border-border bg-background px-6 h-10 shrink-0 sticky top-0 z-10">
      <div className="flex items-center gap-1.5 mr-4">
        <span className="size-5 rounded bg-primary flex items-center justify-center shrink-0">
          <svg viewBox="0 0 16 16" fill="none" className="size-3" aria-hidden="true">
            <rect x="2" y="3" width="12" height="1.5" rx="0.75" fill="currentColor" className="text-primary-foreground" />
            <rect x="2" y="7" width="8" height="1.5" rx="0.75" fill="currentColor" className="text-primary-foreground" />
            <rect x="2" y="11" width="10" height="1.5" rx="0.75" fill="currentColor" className="text-primary-foreground" />
          </svg>
        </span>
        <span className="font-semibold text-sm tracking-tight text-foreground">
          LogScope
        </span>
      </div>

      <Separator orientation="vertical" className="h-5 mx-2" />

      <NavLink
        to="/analyze"
        className={({ isActive }) =>
          `px-3 py-1.5 text-sm transition-colors duration-100 font-medium border-b-2 ${
            isActive
              ? 'border-primary text-foreground'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          }`
        }
      >
        Analyzer
      </NavLink>
      <NavLink
        to="/generate"
        className={({ isActive }) =>
          `px-3 py-1.5 text-sm transition-colors duration-100 font-medium border-b-2 ${
            isActive
              ? 'border-primary text-foreground'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          }`
        }
      >
        Generator
      </NavLink>

      <button
        onClick={toggle}
        className="ml-auto text-muted-foreground hover:text-foreground transition-colors"
        aria-label={dark ? 'Switch to light mode' : 'Switch to dark mode'}
      >
        {dark ? <Sun className="size-4" /> : <Moon className="size-4" />}
      </button>
    </nav>
  )
}
