import { useEffect, useState } from 'react';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from './ui/card'

interface SummaryData {
  today_hours: { value: number; change: number }
  week_hours: { value: number; change: number }
  sessions_today: { value: number; change: number }
  productivity_score: { value: number; change: number }
}

interface StatCard {
  title: string
  main: string
  subtitle: string
  change: string
}

export default function StatCards() {
  const [data, setData] = useState<SummaryData | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchSummary = async () => {
      try {
        const res = await fetch('http://localhost:8080/api/summary')
        if (!res.ok) throw new Error('Failed to fetch summary data')
        const json = await res.json()
        setData(json)
      } catch (err: any) {
        setError(err.message || 'Something went wrong')
      } finally {
        setLoading(false)
      }
    }

    fetchSummary()
  }, [])

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
  )
}
