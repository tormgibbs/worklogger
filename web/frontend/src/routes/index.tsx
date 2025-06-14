import Sessions from '@/components/Sessions'
import StatCards from '@/components/StatsCards'
import TabbedStats from '@/components/TabbedStats'

import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  component: App,
})

function App() {
  return (
    <div className="flex flex-1 flex-col p-4">
      <div>
        <h1 className="text-2xl font-bold">Dashboard</h1>
        <p className="text-gray-600">
          Track your development time and boost productivity
        </p>
      </div>

      <StatCards />

      <div className="mt-10">
        <TabbedStats />
      </div>

      <div className="mt-10">
        <Sessions />
      </div>
    </div>
  )
}
