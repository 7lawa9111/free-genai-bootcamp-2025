import { useEffect, useState } from 'react'
import StudyActivity from '@/components/StudyActivity'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

// Define the API response type
type ActivityResponse = {
  id: number
  name: string
  url: string
  preview_url: string
  launch_url?: string
  title?: string
  description?: string
}

// Define the type expected by StudyActivity component
type ActivityCardProps = {
  id: number
  title: string
  preview_url: string
  launch_url: string
}

type Group = {
  id: number
  name: string
}

export default function StudyActivities() {
  const [activities, setActivities] = useState<ActivityResponse[]>([])
  const [groups, setGroups] = useState<Group[]>([])
  const [selectedGroupId, setSelectedGroupId] = useState<string>('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const loadData = async () => {
      try {
        // Fetch activities
        const activitiesResponse = await fetch('http://localhost:5001/api/study-activities')
        if (!activitiesResponse.ok) {
          throw new Error('Failed to fetch activities')
        }
        const activitiesData = await activitiesResponse.json()
        setActivities(activitiesData)

        // Fetch groups - update the endpoint
        const groupsResponse = await fetch('http://localhost:5001/api/word-groups')
        if (!groupsResponse.ok) {
          throw new Error('Failed to fetch groups')
        }
        const groupsData = await groupsResponse.json()
        setGroups(groupsData.items) // Note: API returns {items: Group[]}

        setLoading(false)
      } catch (err) {
        console.error('Error:', err)
        setError(err instanceof Error ? err.message : 'Failed to load data')
        setLoading(false)
      }
    }
    loadData()
  }, [])

  if (loading) {
    return <div className="text-center">Loading study activities...</div>
  }

  if (error) {
    return <div className="text-red-500">Error: {error}</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">Study Activities</h1>
        <Select value={selectedGroupId} onValueChange={setSelectedGroupId}>
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Select a group" />
          </SelectTrigger>
          <SelectContent>
            {groups.map(group => (
              <SelectItem key={group.id} value={group.id.toString()}>
                {group.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {activities.map((activity) => (
          <StudyActivity 
            key={activity.id} 
            activity={{
              id: activity.id,
              title: activity.title || activity.name,
              preview_url: activity.preview_url,
              launch_url: activity.launch_url || activity.url
            }}
            groupId={selectedGroupId}
          />
        ))}
      </div>
    </div>
  )
}