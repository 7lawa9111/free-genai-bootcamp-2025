import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { Button } from '@/components/ui/button'

type DashboardData = {
  study_progress: {
    total_studied: number
    total_words: number
    mastery_progress: number
  }
  latest_session: {
    activity_name: string
    group_name: string
    date: string
    correct_count: number
    wrong_count: number
  } | null
  quick_stats: {
    success_rate: number
    study_streak: number
    active_groups: number
  }
}

export default function Dashboard() {
  const [data, setData] = useState<DashboardData | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchDashboard = async () => {
      try {
        const response = await fetch('http://localhost:5001/api/dashboard')
        if (!response.ok) {
          throw new Error('Failed to fetch dashboard data')
        }
        const data = await response.json()
        console.log('Dashboard data:', data)  // Debug log
        setData(data)
      } catch (err) {
        console.error('Error fetching dashboard:', err)
        setError(err instanceof Error ? err.message : 'Failed to load dashboard')
      } finally {
        setLoading(false)
      }
    }

    fetchDashboard()
    
    // Poll for updates every 5 seconds
    const interval = setInterval(fetchDashboard, 5000)
    return () => clearInterval(interval)
  }, [])

  if (loading) return <div>Loading...</div>
  if (error) return <div>Error: {error}</div>
  if (!data) return <div>No data available</div>

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">Dashboard</h1>
        <Button asChild>
          <Link to="/study-activities">Start Studying →</Link>
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* Last Study Session */}
        <div className="bg-sidebar rounded-lg p-6">
          <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
            <span className="i-lucide-clock text-xl" />
            Last Study Session
          </h2>
          {data.latest_session ? (
            <>
              <div>{data.latest_session.activity_name}</div>
              <div className="text-gray-500 dark:text-gray-400">{data.latest_session.date}</div>
              <div className="mt-2">
                <span className="text-green-500">✓ {data.latest_session.correct_count} correct</span>
                {' '}
                <span className="text-red-500">✗ {data.latest_session.wrong_count} wrong</span>
              </div>
              <Link to={`/groups/${data.latest_session.group_name}`} className="text-blue-500 hover:underline mt-2 block">
                View Group →
              </Link>
            </>
          ) : (
            <div className="text-gray-500">No study sessions yet</div>
          )}
        </div>

        {/* Study Progress */}
        <div className="bg-sidebar rounded-lg p-6">
          <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
            <span className="i-lucide-activity text-xl" />
            Study Progress
          </h2>
          <div className="text-2xl font-bold">
            {data.study_progress.total_studied} / {data.study_progress.total_words}
          </div>
          <div className="mt-2">Total Words Studied</div>
          <div className="mt-4">
            <div>Mastery Progress</div>
            <div className="w-full bg-gray-200 rounded-full h-2.5 dark:bg-gray-700 mt-2">
              <div 
                className="bg-blue-600 h-2.5 rounded-full" 
                style={{ width: `${data.study_progress.mastery_progress}%` }}
              />
            </div>
            <div className="text-right mt-1">{data.study_progress.mastery_progress}%</div>
          </div>
          <Link to="/words" className="text-blue-500 hover:underline mt-4 block">
            Browse word groups →
          </Link>
        </div>

        {/* Quick Stats */}
        <div className="bg-sidebar rounded-lg p-6">
          <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
            <span className="i-lucide-trophy text-xl" />
            Quick Stats
          </h2>
          <div className="space-y-4">
            <div>
              <div className="text-gray-500">Success Rate</div>
              <div className="text-2xl font-bold">{data.quick_stats.success_rate}%</div>
            </div>
            <div>
              <div className="text-gray-500">Study Sessions</div>
              <div className="text-2xl font-bold">{data.latest_session ? 1 : 0}</div>
            </div>
            <div>
              <div className="text-gray-500">Active Groups</div>
              <div className="text-2xl font-bold">{data.quick_stats.active_groups}</div>
            </div>
            <div>
              <div className="text-gray-500">Study Streak</div>
              <div className="text-2xl font-bold">{data.quick_stats.study_streak} days</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}