import { StrictMode } from 'react'
import ReactDOM from 'react-dom/client'
import { RouterProvider, createRouter, NotFoundRoute } from '@tanstack/react-router'
import { QueryClient, QueryClientProvider, } from '@tanstack/react-query';
import './index.css'
import { Route } from './routes/__root.tsx'
import AuthProvider, { useAuth } from './context/authContext.tsx';

// Import the generated route tree
import { routeTree } from './routeTree.gen'

const notFoundRoute = new NotFoundRoute({
    getParentRoute: () => Route,
    component: () => "404 Not Found",
})
// Create a new router instance
const router = createRouter({
    routeTree,
    notFoundRoute,
    defaultPreload: 'intent',
    context: {
        auth: undefined!
    }
})
const queryClient = new QueryClient();

// Register the router instance for type safety
declare module '@tanstack/react-router' {
    interface Register {
        router: typeof router
    }
}

// Render the app
const rootElement = document.getElementById('root')!
if (!rootElement.innerHTML) {
    const auth = useAuth()
    const root = ReactDOM.createRoot(rootElement)
    root.render(
        <StrictMode>
            <AuthProvider>
                <QueryClientProvider client={queryClient}>
                    <RouterProvider router={router} context={{ auth }} />
                </QueryClientProvider>
            </AuthProvider>
        </StrictMode>
    )
}
