
import { createFileRoute } from '@tanstack/react-router'
export const Route = createFileRoute('/')({
  component: Landing
})
function Landing() {
  return (
    <div className="bg-zinc-950 min-h-screen py-2 px-2">
      <div className="container mx-auto rounded-md bg-zinc-900 flex flex-col items-center ">
        <h2 className="font-sans text-zinc-200 tracking-tight text-3xl text-center font-bold ">
          Landing Page
        </h2>
      </div>
    </div>
  );
}
