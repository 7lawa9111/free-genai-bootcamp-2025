import { useState, useEffect } from 'react'
import StudySessionsTable from '@/components/StudySessionsTable'

type Session = {
  id: number
  group_name: string
  group_id: number
  activity_id: number
  activity_name: string
  start_time: string
  end_time: string
  review_items_count: number
}

type StudySessionSortKey = 'start_time' | 'group_name' | 'activity_name' | 'review_items_count'

export default function SessionsPage() {
  const [sessions, setSessions] = useState<Session[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [sortKey, setSortKey] = useState<StudySessionSortKey>('start_time')
  const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('desc')

  useEffect(() => {
    const fetchSessions = async () => {
      try {
        const response = await fetch('http://localhost:5001/api/study-sessions')
        if (!response.ok) {
          throw new Error('Failed to fetch sessions')
        }
        const data = await response.json()
        setSessions(data.items)
      } catch (err) {
        console.error('Error fetching sessions:', err)
        setError(err instanceof Error ? err.message : 'Failed to load sessions')
      } finally {
        setLoading(false)
      }
    }

    fetchSessions()
  }, [])

  const handleSort = (key: StudySessionSortKey) => {
    setSortDirection(current => current === 'asc' ? 'desc' : 'asc')
    setSortKey(key)
  }

  if (loading) return <div className="text-center py-4">Loading...</div>
  if (error) return <div className="text-red-500 text-center py-4">{error}</div>

  return (
    <div className="space-y-8">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold text-gray-800 dark:text-white">Study Sessions</h1>
      </div>

      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden">
        <div className="p-6">
          <StudySessionsTable 
            sessions={sessions}
            sortKey={sortKey}
            sortDirection={sortDirection}
            onSort={handleSort}
          />
        </div>
      </div>
    </div>
  )
} 