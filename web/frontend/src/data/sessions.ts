export interface Session {
  id: number
  task: string
  startTime: string
  endTime: string
  duration: string
  status: 'in_progress' | 'paused' | 'ended'
}

export const sessions: Session[] = [
  {
    id: 1,
    task: 'Write API docs',
    startTime: '2025-06-01T09:00:00Z',
    endTime: '2025-06-01T10:30:00Z',
    duration: '1h 30m',
    status: 'ended',
  },
  {
    id: 2,
    task: 'Fix login bug',
    startTime: '2025-06-01T11:00:00Z',
    endTime: '2025-06-01T12:15:00Z',
    duration: '1h 15m',
    status: 'ended',
  },
  {
    id: 3,
    task: 'Refactor dashboard layout',
    startTime: '2025-06-01T13:00:00Z',
    endTime: '2025-06-01T13:45:00Z',
    duration: '45m',
    status: 'paused',
  },
  {
    id: 4,
    task: 'Design new onboarding flow',
    startTime: '2025-06-01T14:00:00Z',
    endTime: '2025-06-01T15:30:00Z',
    duration: '1h 30m',
    status: 'ended',
  },
  {
    id: 5,
    task: 'Review PRs',
    startTime: '2025-06-01T16:00:00Z',
    endTime: '2025-06-01T16:20:00Z',
    duration: '20m',
    status: 'in_progress',
  },
  {
    id: 6,
    task: 'Write unit tests',
    startTime: '2025-06-01T17:00:00Z',
    endTime: '2025-06-01T18:10:00Z',
    duration: '1h 10m',
    status: 'ended',
  },
]
