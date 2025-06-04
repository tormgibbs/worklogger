import type { SummaryData } from '@/types/types'
import { useEffect, useState } from 'react'

export const useSummaryData = (useMock: boolean) => {
  const [data, setData] = useState<SummaryData | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (useMock) {
      const mockData: SummaryData = {
        today_hours: { value: 4.5, change: 12 },
        week_hours: { value: 23.5, change: 8 },
        sessions_today: { value: 5, change: 5 },
        productivity_score: { value: 87, change: 3 },
      }
      setData(mockData)
      setLoading(false)
    } else {
      fetch('http://localhost:8080/api/summary')
        .then((res) => {
          if (!res.ok) throw new Error('Network response not ok')
          return res.json()
        })
        .then((realData: SummaryData) => {
          setData(realData)
          setLoading(false)
        })
        .catch((err: Error) => {
          setError(err.message)
          setLoading(false)
        })
    }
  }, [useMock])

  return { data, loading, error }
}


  const stats = [
    {
      title: "Today's Hours",
      main: '4.5',
      subtitle: 'Total hours tracked today',
      change: '12% from last period',
    },
    {
      title: 'Week Total',
      main: '23.5',
      subtitle: 'Hours tracked this week',
      change: '8% from last period',
    },
    {
      title: 'Sessions Today',
      main: '5',
      subtitle: 'Number of work sessions',
      change: '5% from last period',
    },
    {
      title: 'Productivity Score',
      main: '87%',
      subtitle: 'Based on focus time',
      change: '3% from last period',
    },
  ]
