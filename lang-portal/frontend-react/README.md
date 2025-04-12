# Japanese Learning App Frontend

A React-based frontend for the Japanese language learning application, built with TypeScript, Vite, and Tailwind CSS.

## Features

- Interactive dashboard with study progress tracking
- Word group management
- Study session tracking
- Real-time progress updates
- Responsive design with dark mode support

## Prerequisites

- Node.js (v18 or higher)
- npm or yarn
- Backend server running (see Backend Setup)

## Quick Start

1. Install dependencies:
```bash
npm install
```

2. Start the development server:
```bash
npm run dev
```

The application will be available at `http://localhost:5173`

## Backend Setup

The frontend requires the Flask backend server running on port 5001. To set up the backend:

1. Navigate to the backend directory:
```bash
cd ../backend-flask
```

2. Create and activate a virtual environment:
```bash
python -m venv venv
source venv/bin/activate  # On Windows use: venv\Scripts\activate
```

3. Install dependencies:
```bash
pip install -r requirements.txt
```

4. Start the backend server:
```bash
python app.py
```

The backend will run on `http://localhost:5001`

## Available API Endpoints

The frontend interacts with the following API endpoints:

### Dashboard
- `GET /api/dashboard` - Fetch dashboard data including study progress and stats

### Study Sessions
- `GET /api/study-sessions` - List all study sessions
- `POST /api/study-sessions` - Create a new study session
- `GET /api/study-sessions/:id` - Get specific session details

### Words and Groups
- `GET /api/words` - List all vocabulary words
- `GET /api/groups` - List all word groups
- `POST /api/groups` - Create a new word group
- `GET /api/groups/:id` - Get specific group details

### Study Activities
- `GET /api/study-activities` - List available study activities
- `POST /api/study-activities/:id/start` - Start a study activity
- `POST /api/study-activities/:id/complete` - Complete a study activity

## Development Commands

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run lint` - Run ESLint
- `npm run preview` - Preview production build

## Environment Configuration

The application automatically detects the environment and configures API endpoints:
- Development: `http://localhost:5001`
- Gitpod: Automatically uses the correct workspace URL

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request
