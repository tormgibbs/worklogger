
export interface DailyStat {
  day: string
  hours: number
  sessions?: number
}

export interface DailySessionStat {
  day: string
  sessions: number
}

export interface WeeklyStat {
  week: string
  hours: number
}

export interface WeeklySessionStat {
  week: string
  sessions: number
}

export interface MonthlyStat {
  month: string
  hours: number
  sessions?: number
}

export interface MonthlySessionStat {
  month: string
  sessions: number
}

export const dailyHours: DailyStat[] = [
  { day: 'Mon', hours: 4.0 },
  { day: 'Tue', hours: 3.5 },
  { day: 'Wed', hours: 5.25 },
  { day: 'Thu', hours: 0.0 },
  { day: 'Fri', hours: 2.75 },
  { day: 'Sat', hours: 1.0 },
  { day: 'Sun', hours: 0.0 },
]

// [
//   { "date": "2025-05-26", "hours": 3.5, "sessions": 2 },
//   { "date": "2025-05-27", "hours": 5, "sessions": 3 },
// ]

export const dailySessions: DailySessionStat[] = [
  { day: 'Mon', sessions: 2 },
  { day: 'Tue', sessions: 1 },
  { day: 'Wed', sessions: 3 },
  { day: 'Thu', sessions: 0 },
  { day: 'Fri', sessions: 1 },
  { day: 'Sat', sessions: 1 },
  { day: 'Sun', sessions: 0 },
]

export const weeklyHours: WeeklyStat[] = [
  { week: 'May 1–5', hours: 14.5 },
  { week: 'May 6–12', hours: 21.0 },
  { week: 'May 13–19', hours: 18.25 },
  { week: 'May 20–26', hours: 25.75 },
  { week: 'May 27–31', hours: 9.0 },
]

// [
//   { "week_start": "2025-05-06", "hours": 23.5, "sessions": 12 },
//   { "week_start": "2025-05-13", "hours": 19, "sessions": 10 },
// ]

export const weeklySessions: WeeklySessionStat[] = [
  { week: 'May 1–5', sessions: 5 },
  { week: 'May 6–12', sessions: 8 },
  { week: 'May 13–19', sessions: 6 },
  { week: 'May 20–26', sessions: 10 },
  { week: 'May 27–31', sessions: 3 },
]

export const monthlyStats0: MonthlyStat[] = [
  { month: 'Jan', hours: 42.5, sessions: 15 },
  { month: 'Feb', hours: 38.0, sessions: 12 },
  { month: 'Mar', hours: 54.25, sessions: 20 },
  { month: 'Apr', hours: 47.0, sessions: 18 },
  { month: 'May', hours: 88.5, sessions: 28 },
]

// [
//   { "month": "2025-05", "hours": 92, "sessions": 45 },
//   { "month": "2025-04", "hours": 88, "sessions": 42 },
// ]

export const monthlySessions: MonthlySessionStat[] = [
  { "month": "Jan", "sessions": 15 },
  { "month": "Feb", "sessions": 12 },
  { "month": "Mar", "sessions": 20 },
  { "month": "Apr", "sessions": 18 },
  { "month": "May", "sessions": 28 }
]

export const dailyStats = [
  { date: '2025-05-26', hours: 4.0, sessions: 2 },
  { date: '2025-05-27', hours: 3.5, sessions: 1 },
  { date: '2025-05-28', hours: 5.25, sessions: 3 },
  { date: '2025-05-29', hours: 0, sessions: 0 },
  { date: '2025-05-30', hours: 2.75, sessions: 1 },
  { date: '2025-05-31', hours: 1, sessions: 1 },
  { date: '2025-06-01', hours: 0, sessions: 0 },
]

export const weeklyStats = [
  { week_start: '2025-05-05', hours: 14.5, sessions: 5 },
  { week_start: '2025-05-12', hours: 21.0, sessions: 8 },
  { week_start: '2025-05-19', hours: 15, sessions: 8 },
  { week_start: '2025-05-26', hours: 23.5, sessions: 12 },
  { week_start: '2025-06-02', hours: 18.75, sessions: 9 },
]

export const monthlyStats = [
  { month: '2025-04', hours: 78, sessions: 33 },
  { month: '2025-05', hours: 92, sessions: 45 },
  { month: '2025-06', hours: 60, sessions: 25 },
  { month: '2025-07', hours: 75, sessions: 30 },
  { month: '2025-08', hours: 80, sessions: 35 },
]