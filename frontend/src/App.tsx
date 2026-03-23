import { createBrowserRouter, RouterProvider, Navigate } from 'react-router-dom'
import { Layout } from './components/shared/Layout'
import { AnalyzerPage } from './pages/AnalyzerPage'
import { GeneratorPage } from './pages/GeneratorPage'

const router = createBrowserRouter([
  {
    path: '/',
    element: <Layout />,
    children: [
      { index: true, element: <Navigate to="/analyze" replace /> },
      { path: 'analyze', element: <AnalyzerPage /> },
      { path: 'generate', element: <GeneratorPage /> },
    ],
  },
])

export default function App() {
  return <RouterProvider router={router} />
}
