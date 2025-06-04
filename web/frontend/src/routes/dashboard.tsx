import TabbedStats from '@/components/TabbedStats'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { sessions } from '@/data/sessions'

import { useSummaryData } from '@/hooks/useSummaryData'
import { columns } from '@/sessions/columns'
import { DataTable } from '@/sessions/data-table'
import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'

export const Route = createFileRoute('/dashboard')({
  component: RouteComponent,
})

function RouteComponent() {
  const [useMock, setUseMock] = useState(false)
  const { data, loading, error } = useSummaryData(useMock)

  if (loading) return <p>Loading...</p>
  if (error) return <p>Error: {error}</p>
  if (!data) return null

  const stats = [
    {
      title: "Today's Hours",
      main: data.today_hours.value.toString(),
      subtitle: 'Total hours tracked today',
      change: `${data.today_hours.change}% from last period`,
    },
    {
      title: 'Week Total',
      main: data.week_hours.value.toString(),
      subtitle: 'Hours tracked this week',
      change: `${data.week_hours.change}% from last period`,
    },
    {
      title: 'Sessions Today',
      main: data.sessions_today.value.toString(),
      subtitle: 'Number of work sessions',
      change: `${data.sessions_today.change}% from last period`,
    },
    {
      title: 'Productivity Score',
      main: `${data.productivity_score.value}%`,
      subtitle: 'Based on focus time',
      change: `${data.productivity_score.change}% from last period`,
    },
  ]

  return (
    <div className="flex flex-1 flex-col p-4">
      <div>
        <h1 className="text-2xl font-bold">Dashboard</h1>
        <p className="text-gray-600">
          Track your development time and boost productivity
        </p>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4 mt-10">
        {stats.map(({ title, main, subtitle, change }) => (
          <Card key={title}>
            <CardHeader>
              <CardTitle className="text-lg font-semibold">{title}</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-bold">{main}</p>
              <p className="text-gray-500 mt-1">{subtitle}</p>
              <p
                className={`text-sm mt-2 ${parseFloat(change) < 0 ? 'text-red-600' : 'text-green-600'
                  }`}
              >
                {change}
              </p>
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="mt-10">
        <TabbedStats />
      </div>

      <div className="mt-10">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg font-semibold">
              Recent Sessions
            </CardTitle>
            <CardDescription className="text-sm text-gray-500">
              View your recent work sessions and track progress
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="container mx-auto">
              <DataTable columns={columns} data={sessions} />
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
