import React, { useEffect, useState } from 'react'
import { useParams, useLocation } from 'react-router-dom'
import { useNavigation } from '@/context/NavigationContext'
import StudySessionsTable from '@/components/StudySessionsTable'
import Pagination from '@/components/Pagination'

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

type StudyActivity = {
  id: number
  preview_url: string
  title: string
  description: string
  launch_url: string
}

type PaginatedSessions = {
  items: Session[]
  total: number
  page: number
  per_page: number
  total_pages: number
}

const ITEMS_PER_PAGE = 10

export default function StudyActivityShow() {
  const { id } = useParams<{ id: string }>()
  const location = useLocation()
  const searchParams = new URLSearchParams(location.search)
  const groupId = searchParams.get('group_id')
  const { setCurrentStudyActivity } = useNavigation()
  const [activity, setActivity] = useState<StudyActivity | null>(null)
  const [sessionData, setSessionData] = useState<PaginatedSessions | null>(null)
  const [currentPage, setCurrentPage] = useState(1)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchActivity = async () => {
      try {
        const response = await fetch(`http://localhost:5001/api/study-activities/${id}`)
        if (!response.ok) {
          throw new Error('Failed to fetch activity details')
        }
        const data = await response.json()
        setActivity(data)
        setCurrentStudyActivity(data)
        
        const sessionsResponse = await fetch(
          `http://localhost:5001/api/study_activities/${id}/sessions?page=${currentPage}&per_page=${ITEMS_PER_PAGE}`
        )
        if (!sessionsResponse.ok) {
          throw new Error('Failed to fetch sessions')
        }
        const sessionsData = await sessionsResponse.json()
        setSessionData({
          items: sessionsData.items.map((item: any) => ({
            id: item.id,
            group_name: item.group_name,
            group_id: item.group_id,
            activity_id: item.activity_id,
            activity_name: item.activity_name,
            start_time: item.start_time,
            end_time: item.end_time,
            review_items_count: item.review_items_count
          })),
          total: sessionsData.total,
          page: sessionsData.page,
          per_page: sessionsData.per_page,
          total_pages: sessionsData.total_pages
        })
      } catch (err) {
        console.error('Error fetching activity:', err)
        setError(err instanceof Error ? err.message : 'Failed to load activity')
      } finally {
        setLoading(false)
      }
    }

    fetchActivity()
  }, [id, currentPage, setCurrentStudyActivity])

  // Clean up when unmounting
  useEffect(() => {
    return () => {
      setCurrentStudyActivity(null)
    }
  }, [setCurrentStudyActivity])

  const handleLaunchActivity = async (groupId: number) => {
    try {
      console.log('Creating session for group:', groupId, 'activity:', activity.id)
      
      // First create a new study session
      const sessionResponse = await fetch('http://localhost:5001/api/study_sessions', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          group_id: groupId,
          study_activity_id: activity.id
        }),
      })

      if (!sessionResponse.ok) {
        const errorText = await sessionResponse.text()
        throw new Error(`Failed to create session: ${errorText}`)
      }

      const sessionData = await sessionResponse.json()
      console.log('Created session:', sessionData)

      // Then redirect to the activity URL with the session parameters
      const redirectUrl = `${activity.launch_url}?group_id=${groupId}&session_id=${sessionData.id}`
      console.log('Redirecting to:', redirectUrl)
      window.location.href = redirectUrl

    } catch (err) {
      console.error('Error launching activity:', err)
      setError(err instanceof Error ? err.message : 'Failed to create study session')
    }
  }

  if (loading) {
    return <div className="text-center py-4">Loading...</div>
  }

  if (error || !activity) {
    return <div className="text-red-500 text-center py-4">{error || 'Activity not found'}</div>
  }

  return (
    <div className="container mx-auto p-4">
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
        <h1 className="text-2xl font-bold mb-4">{activity.title}</h1>
        <img 
          src={activity.preview_url} 
          alt={activity.title} 
          className="w-full h-64 object-cover rounded-lg mb-4"
        />
        <p className="text-gray-600 dark:text-gray-300 mb-4">{activity.description}</p>
        <div className="flex justify-end">
          <button 
            onClick={() => handleLaunchActivity(groupId)}
            className="bg-primary text-primary-foreground px-4 py-2 rounded hover:bg-primary/90"
          >
            Launch Activity
          </button>
        </div>
      </div>

      {sessionData && sessionData.items.length > 0 && (
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 mt-4">
          <h2 className="text-xl font-semibold mb-4">Study Sessions</h2>
          <StudySessionsTable sessions={sessionData.items} />
          {sessionData.total_pages > 1 && (
            <div className="mt-4">
              <Pagination
                currentPage={currentPage}
                totalPages={sessionData.total_pages}
                onPageChange={setCurrentPage}
              />
            </div>
          )}
        </div>
      )}
    </div>
  )
}