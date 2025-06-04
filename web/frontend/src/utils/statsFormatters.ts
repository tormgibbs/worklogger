// utils/statsFormatters.ts

// For daily stats
export const formatDailyStats = (data: { date: string; hours: number; sessions: number }[]) => ({
  hours: data.map(d => ({
    day: new Date(d.date).toLocaleDateString('en-US', { weekday: 'short' }), // e.g. "Mon"
    hours: d.hours,
  })),
  sessions: data.map(d => ({
    day: new Date(d.date).toLocaleDateString('en-US', { weekday: 'short' }),
    sessions: d.sessions,
  })),
})

// For weekly stats
export const formatWeeklyStats = (data: { week_start: string; hours: number; sessions: number }[]) => ({
  hours: data.map(w => ({
    day: new Date(w.week_start).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
    }), // e.g. "May 6"
    hours: w.hours,
  })),
  sessions: data.map(w => ({
    day: new Date(w.week_start).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
    }),
    sessions: w.sessions,
  })),
})

// For monthly stats
export const formatMonthlyStats = (data: { month: string; hours: number; sessions: number }[]) => ({
  hours: data.map(m => ({
    day: new Date(`${m.month}-01`).toLocaleDateString('en-US', {
      month: 'short',
      year: 'numeric',
    }), // e.g. "May 2025"
    hours: m.hours,
  })),
  sessions: data.map(m => ({
    day: new Date(`${m.month}-01`).toLocaleDateString('en-US', {
      month: 'short',
      year: 'numeric',
    }),
    sessions: m.sessions,
  })),
})
