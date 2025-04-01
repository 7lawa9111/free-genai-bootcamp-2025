import { ThemeProvider } from "@/components/theme-provider"
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import AppSidebar from '@/components/Sidebar'
import Breadcrumbs from '@/components/Breadcrumbs'
import AppRouter from '@/components/AppRouter'
import { NavigationProvider } from '@/context/NavigationContext'
import StudySessionsPage from '@/pages/StudySessionsPage'
import StudyActivities from '@/pages/StudyActivities'
import StudyActivityShow from '@/pages/StudyActivityShow'

import {
  SidebarInset,
  SidebarProvider
} from "@/components/ui/sidebar"

const API_BASE_URL = 'http://localhost:5001';

export const fetchStudyActivities = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/api/study_activities`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Error fetching study activities:', error);
    throw error;
  }
};

export default function App() {
  return (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <NavigationProvider>
        <Router>
          <SidebarProvider>
            <AppSidebar />
            <SidebarInset>
              <Breadcrumbs />
              <AppRouter />
            </SidebarInset>
          </SidebarProvider>  
        </Router>
      </NavigationProvider>
    </ThemeProvider>
  )
}