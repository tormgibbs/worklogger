// hooks/useStatsData.ts
import { useState, useEffect } from 'react'
import {
  formatDailyStats,
  formatWeeklyStats,
  formatMonthlyStats,
} from '@/utils/statsFormatters'

type RawDaily = { date: string; hours: number; sessions: number }
type RawWeekly = { week_start: string; hours: number; sessions: number }
type RawMonthly = { month: string; hours: number; sessions: number }

export const useStatsData = (useMock: boolean) => {
  const [daily, setDaily] = useState({ hours: [], sessions: [] })
  const [weekly, setWeekly] = useState({ hours: [], sessions: [] })
  const [monthly, setMonthly] = useState({ hours: [], sessions: [] })
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (useMock) {
      const mockDaily: RawDaily[] = [
        { date: '2025-05-26', hours: 4, sessions: 2 },
        { date: '2025-05-27', hours: 3.5, sessions: 1 },
        { date: '2025-05-28', hours: 5.25, sessions: 3 },
        { date: '2025-05-29', hours: 0.0, sessions: 0 },
        { date: '2025-05-30', hours: 2.75, sessions: 1 },
        { date: '2025-05-31', hours: 1.0, sessions: 1 },
        { date: '2025-06-01', hours: 0.0, sessions: 0 },
      ]
      const mockWeekly: RawWeekly[] = [
        { week_start: '2025-05-06', hours: 23.5, sessions: 12 },
        { week_start: '2025-05-13', hours: 18.2, sessions: 9 },
      ]
      const mockMonthly: RawMonthly[] = [
        { month: '2025-04', hours: 78, sessions: 34 },
        { month: '2025-05', hours: 92, sessions: 45 },
      ]
      setDaily(formatDailyStats(mockDaily))
      setWeekly(formatWeeklyStats(mockWeekly))
      setMonthly(formatMonthlyStats(mockMonthly))
      setLoading(false)
    } else {
      const fetchAll = async () => {
        try {
          const [dRes, wRes, mRes] = await Promise.all([
            fetch('http://localhost:8080/api/stats/daily'),
            fetch('http://localhost:8080/api/stats/weekly'),
            fetch('http://localhost:8080/api/stats/monthly'),
          ])
          if (!dRes.ok || !wRes.ok || !mRes.ok) throw new Error('Failed to fetch')

          const [d, w, m] = await Promise.all([
            dRes.json() as Promise<RawDaily[]>,
            wRes.json() as Promise<RawWeekly[]>,
            mRes.json() as Promise<RawMonthly[]>,
          ])

          setDaily(formatDailyStats(d))
          setWeekly(formatWeeklyStats(w))
          setMonthly(formatMonthlyStats(m))
        } catch (err: any) {
          setError(err.message ?? 'Unknown error')
        } finally {
          setLoading(false)
        }
      }

      fetchAll()
    }
  }, [useMock])

  return {
    daily,
    weekly,
    monthly,
    loading,
    error,
  }
}
