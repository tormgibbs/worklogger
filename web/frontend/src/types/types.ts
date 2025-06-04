export interface SummaryMetric {
  value: number
  change: number
}

export interface SummaryData {
  today_hours: SummaryMetric
  week_hours: SummaryMetric
  sessions_today: SummaryMetric
  productivity_score: SummaryMetric
}