# Language Learning Portal Backend

A Go backend service for a language learning application.

## Run

```sh
go run cmd/server/main.go
```

## Test Code

When running tests, use test environment for the go app:
```sh
APP_ENV=test go run cmd/server/main.go
```

Running a test:
```sh
cd api_tests
bundle exec rspec spec/words_spec.rb
```

Running all tests:
```sh
cd api_tests
bundle exec rspec
```

## Kill if already running

If the port is already in use from running go app prior you can kill the process:
```sh
lsof -ti:8080 | xargs kill -9
```

## Database Management

The application uses SQLite for data storage. Test and development environments use separate database files.

### Database Initialization
The database is automatically initialized when the server starts. Test data is automatically loaded in test environment.

### Manual Database Reset
You can use the API endpoints to reset the database:

```sh
# Reset study history only
curl -X POST http://localhost:8080/api/reset_history

# Full system reset
curl -X POST http://localhost:8080/api/full_reset
```

## API Endpoints

### Words
- GET `/api/words` - List all words
- GET `/api/words/:id` - Get specific word

### Groups
- GET `/api/groups` - List all groups
- GET `/api/groups/:id` - Get specific group
- GET `/api/groups/:id/words` - List words in a group
- GET `/api/groups/:id/study_sessions` - List study sessions for a group

### Study Sessions
- GET `/api/study_sessions` - List all study sessions
- GET `/api/study_sessions/:id` - Get specific study session
- POST `/api/study_sessions/:id/words/:word_id/review` - Record word review

### Study Activities
- GET `/api/study_activities/:id` - Get specific study activity
- GET `/api/study_activities/:id/study_sessions` - List sessions for an activity
- POST `/api/study_activities` - Create new study activity

### Dashboard
- GET `/api/dashboard/quick-stats` - Get dashboard statistics
- GET `/api/dashboard/study_progress` - Get study progress

### Running mage commands

```sh
go run github.com/magefile/mage@latest testdb
go run github.com/magefile/mage@latest dbinit
go run github.com/magefile/mage@latest seed
```