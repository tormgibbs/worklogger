import type { Session } from '@/data/sessions'
import type { ColumnDef } from '@tanstack/react-table'

export const columns: ColumnDef<Session>[] = [
  {
    accessorKey: 'task',
    header: 'Task',
  },
  {
    accessorKey: 'startTime',
    header: 'Start Time',
    cell: ({ row }) =>
      new Date(row.getValue('startTime')).toLocaleString('en-US', {
        month: 'short',
        day: 'numeric',
        hour: 'numeric',
        minute: '2-digit',
        hour12: true,
      }),
  },
  {
    accessorKey: 'endTime',
    header: 'End Time',
    cell: ({ row }) =>
      new Date(row.getValue('startTime')).toLocaleString('en-US', {
        month: 'short',
        day: 'numeric',
        hour: 'numeric',
        minute: '2-digit',
        hour12: true,
      }),
  },
  {
    accessorKey: 'duration',
    header: 'Duration',
  },
  {
    accessorKey: 'status',
    header: 'Status',
    cell: ({ row }) => {
      const rawStatus = (row.getValue('status') as string) ?? ''

      const statusStyles: Record<string, string> = {
        in_progress: 'bg-blue-300 border p-1 px-2 text-sm rounded-full',
        paused: 'bg-gray-200 border p-1 px-2 text-sm rounded-full',
        ended: 'bg-green-300 border p-1 px-2 text-sm rounded-full',
      }
      const style = statusStyles[rawStatus] || 'text-gray-600'
      const label = rawStatus.replaceAll('_', ' ').toUpperCase()

      return <span className={style}>{label}</span>
    },
  },
]
