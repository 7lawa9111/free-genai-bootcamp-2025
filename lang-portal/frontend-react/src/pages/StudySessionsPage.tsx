import React, { useEffect, useState } from 'react'
import StudySessionsTable, { StudySessionSortKey } from '@/components/StudySessionsTable'

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

export default function StudySessionsPage() {
  const [sessions, setSessions] = useState<Session[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [sortKey, setSortKey] = useState<StudySessionSortKey>('start_time')
  const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('desc')

  useEffect(() => {
    const createSession = async (groupId: string, sessionId: string) => {
      try {
        const response = await fetch('http://localhost:5001/api/study_sessions', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            group_id: parseInt(groupId),
            study_activity_id: 1, // You'll need to determine which activity this is for
          }),
        })
        if (!response.ok) {
          throw new Error('Failed to create session')
        }
        const data = await response.json()
        setSessions(prev => [...prev, data])
      } catch (err) {
        console.error('Error creating session:', err)
      }
    }

    // Check URL parameters
    const params = new URLSearchParams(window.location.search)
    const groupId = params.get('group_id')
    const sessionId = params.get('session_id')
    if (groupId && sessionId) {
      createSession(groupId, sessionId)
    }
  }, [])

  useEffect(() => {
    const fetchSessions = async () => {
      try {
        console.log('Fetching sessions...')
        const response = await fetch('http://localhost:5001/api/study_sessions')
        console.log('Response status:', response.status)
        
        if (!response.ok) {
          const errorText = await response.text()
          console.error('Error response:', errorText)
          throw new Error(`Failed to fetch sessions: ${response.status} ${errorText}`)
        }
        
        const data = await response.json()
        console.log('Fetched sessions:', data)
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
    setSortKey(key)
    setSortDirection(prevDirection => prevDirection === 'asc' ? 'desc' : 'asc')
  }

  if (loading) return null
  if (error) return (
    <div className="text-red-500 text-center py-4">
      {error}
    </div>
  )
  if (!sessions) return null

  return (
    <StudySessionsTable 
      sessions={sessions}
      sortKey={sortKey}
      sortDirection={sortDirection}
      onSort={handleSort}
    />
  )
} 