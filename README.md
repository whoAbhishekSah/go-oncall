# Go OnCall Scheduler

A simple Go web application for managing on-call schedules with Slack notifications.

## Features

- Create teams with list of users
- Create on-call schedules with rotation periods
- Automatic rotation with Slack notifications
- Web UI for managing teams and schedules

## Setup

### Option 1: Docker Compose (Recommended)

1. Clone the repository and navigate to the project directory

2. Copy the environment file:
```bash
cp .env.example .env
```

3. Edit `.env` file with your Slack credentials:
```bash
SLACK_TOKEN=xoxb-your-slack-bot-token-here
SLACK_CHANNEL=#oncall
```

4. Start the application with Docker Compose:
```bash
docker-compose up --build
```

5. Run database migrations:
```bash
docker-compose exec app ./migrate.sh docker
```

6. Visit `http://localhost:8080` to access the web interface.

### Option 2: Local Development

1. Install dependencies:
```bash
go mod tidy
```

2. Set up PostgreSQL database:
```bash
createdb oncall
```

3. Set environment variables:
```bash
export DATABASE_URL="postgres://localhost/oncall?sslmode=disable"
export SLACK_TOKEN="your-slack-bot-token"
export SLACK_CHANNEL="#oncall"
```

4. Run database migrations:
```bash
./migrate.sh prod
```

5. Run the application:
```bash
go run .
```

6. Visit `http://localhost:8080` to access the web interface.

## Usage

1. **Create a Team**: Add a team with a list of users
2. **Create a Schedule**: Set up an on-call schedule with:
   - Team ID
   - Start and end times
   - Rotation period in hours
   - List of participants from the team
3. **Automatic Rotation**: The system will automatically rotate on-call assignments and send Slack notifications

## Database Migrations

The application uses a migration script to set up the database schema:

- **For Docker Compose**: `./migrate.sh docker`
- **For Production**: `./migrate.sh prod` (requires `DATABASE_URL` environment variable)

The migration script will:
- Create all necessary tables with proper indexes
- Handle database connection waiting for Docker environments
- Log migration progress
- Track applied migrations in a `migration_log` table

## Environment Variables

- `DATABASE_URL`: PostgreSQL connection string
- `SLACK_TOKEN`: Slack bot token for notifications
- `SLACK_CHANNEL`: Slack channel for notifications (default: #oncall)

## Files Structure

- `main.go` - Web server and routing
- `models.go` - Data structures
- `database.go` - Database operations
- `handlers.go` - HTTP handlers and web UI
- `scheduler.go` - On-call rotation logic
- `slack.go` - Slack notifications
- `migrate.sh` - Database migration script
- `migrations/001_init.sql` - Initial database schema
- `docker-compose.yml` - Docker Compose configuration
- `Dockerfile` - Docker image configuration