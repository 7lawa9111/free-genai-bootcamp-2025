import { Link, useLocation } from 'react-router-dom'
import { Button } from '@/components/ui/button'

type ActivityProps = {
  activity: {
    id: number
    title: string
    preview_url: string
    launch_url: string
  }
  groupId: string
}

export default function StudyActivity({ activity, groupId }: ActivityProps) {
  const location = useLocation();
  const searchParams = new URLSearchParams(location.search);
  const groupIdFromParams = searchParams.get('group_id');

  const getLaunchUrl = async () => {
    try {
      // Create a study session first
      const response = await fetch(`http://localhost:5001/api/study-sessions`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          group_id: parseInt(groupId),
          study_activity_id: activity.id
        })
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`Failed to create study session: ${errorText}`);
      }

      const data = await response.json();
      const sessionId = data.id;

      // Return the launch URL with both group_id and session_id
      const baseUrl = activity.launch_url;
      const url = new URL(baseUrl);
      url.searchParams.set('group_id', groupId);
      url.searchParams.set('session_id', sessionId.toString());
      return url.toString();
    } catch (error) {
      console.error('Error creating study session:', error);
      throw error;
    }
  };

  return (
    <div className="bg-sidebar rounded-lg shadow-md overflow-hidden">
      <img src={activity.preview_url} alt={activity.title} className="w-full h-48 object-cover" />
      <div className="p-4">
        <h3 className="text-xl font-semibold mb-2">{activity.title}</h3>
        <div className="flex justify-between">
          <Button 
            disabled={!groupId}
            asChild
          >
            <a 
              href="#"
              onClick={async (e) => {
                e.preventDefault();
                try {
                  const url = await getLaunchUrl();
                  window.open(url, '_blank');
                } catch (error) {
                  console.error('Launch failed:', error);
                }
              }}
            >
              Launch
            </a>
          </Button>
          <Button asChild variant="outline">
            <Link to={`/study-activities/${activity.id}`}>
              View
            </Link>
          </Button>
        </div>
      </div>
    </div>
  );
}