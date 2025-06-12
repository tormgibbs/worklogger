import type { Session } from "@/data/sessions"
import { useEffect, useState } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "./ui/card"
import { DataTable } from "@/sessions/data-table"
import { columns } from "@/sessions/columns"

const Sessions = () => {
  const [sessions, setSessions] = useState<Session[]>([])
  const [loading, setLoading] = useState(true)

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

  return (
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
          {loading ? (
            <p className="text-sm text-muted-foreground">Loading...</p>
          ) : (
            <DataTable columns={columns} data={sessions} />
          )}
        </div>
      </CardContent>
    </Card>

  )
}

export default Sessions