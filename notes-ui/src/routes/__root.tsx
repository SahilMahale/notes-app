import type { AppContext } from '@/context/authContext'
import {  createRootRouteWithContext, Link, Outlet } from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/router-devtools'

interface RouterAppContext {
    auth: AppContext
}
export const Route = createRootRouteWithContext<RouterAppContext>()({
  component: () => (
    <div className="min-h-screen bg-zinc-950">
      <nav className="border-b border-zinc-800 bg-zinc-900 p-4">
        <div className="container mx-auto flex gap-4">
          <Link 
            to="/" 
            className="text-zinc-300 hover:text-white"
            activeProps={{ className: 'font-bold text-white' }}
          >
            Home
          </Link>
          {/* Add more navigation links as needed */}
        </div>
      </nav>
      
      <main>
        <Outlet />
      </main>
      
      <TanStackRouterDevtools />
    </div>
  ),
})
