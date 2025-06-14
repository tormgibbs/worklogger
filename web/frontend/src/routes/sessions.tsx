import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import type { Session } from '@/data/sessions'
import { columns } from '@/sessions/columns'
import { DataTable } from '@/sessions/data-table'
import { createFileRoute } from '@tanstack/react-router'
import { Download } from 'lucide-react'
import { useEffect, useMemo, useState } from 'react'

export const Route = createFileRoute('/sessions')({
  component: RouteComponent,
})

function RouteComponent() {
  const [sessions, setSessions] = useState<Session[]>([])
  const [loading, setLoading] = useState(true)
  const [searchTerm, setSearchTerm] = useState('')
  const [statusFilter, setStatusFilter] = useState<string>('all')

  useEffect(() => {
    const fetchSessions = async () => {
      try {
        const res = await fetch('http://localhost:8080/api/sessions')
        const data = await res.json()
        setSessions(data)
      } catch (err) {
        console.error('Failed to fetch sessions:', err)
      } finally {
        setLoading(false)
      }
    }

    fetchSessions()
  }, [])

  const filteredSessions = useMemo(() => {
    return sessions.filter((session) => {
      // Filter by search term (task)
      const matchesSearch = searchTerm === '' ||
        session.task.toLowerCase().includes(searchTerm.toLowerCase())

      // Filter by status
      const matchesStatus = statusFilter === 'all' || session.status === statusFilter

      return matchesSearch && matchesStatus
    })
  }, [sessions, searchTerm, statusFilter])

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(e.target.value)
  }

  const handleStatusChange = (value: string) => {
    setStatusFilter(value)
  }

  return (
    <div className='flex flex-1 flex-col p-4'>
      <div className='flex flex-row items-center justify-between mb-4'>
        <div>
          <p className='text-3xl font-bold'>Sessions</p>
          <p className='text-gray-500'>View your work sessions</p>
        </div>
        <Button>
          <Download />
          Export CSV
        </Button>
      </div>

      <div className='flex flex-row items-center gap-4 mb-4'>
        <Input
          placeholder='Search by task...'
          value={searchTerm}
          onChange={handleSearchChange}
        />
        <Select value={statusFilter} onValueChange={handleStatusChange}>
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Filter by status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Sessions</SelectItem>
            <SelectItem value="ended">Ended</SelectItem>
            <SelectItem value="in_progress">In Progress</SelectItem>
            <SelectItem value="paused">Paused</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {(searchTerm || statusFilter !== 'all') && (
        <div className="mb-4 text-sm text-gray-600">
          Showing {filteredSessions.length} of {sessions.length} sessions
          {searchTerm && ` matching "${searchTerm}"`}
          {statusFilter !== 'all' && ` with status "${statusFilter}"`}
        </div>
      )}

      <div className="">
        {loading ? (
          <p className="text-sm text-muted-foreground">Loading...</p>
        ) : (
          <DataTable columns={columns} data={filteredSessions} />
        )}
      </div>
    </div>
  )
}
