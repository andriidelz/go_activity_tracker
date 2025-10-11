# Go Activity Tracker

A Go REST API for recording user activity events and aggregating statistics every 4 hours. Includes a minimal React client.

## Description

- API endpoints:
  - POST /events: Create event (body: {user_id, action, metadata}).
  - GET /events?user_id=42&start=2025-10-01T00:00:00Z&end=2025-10-10T00:00:00Z: Retrieve events.
- Background job: Runs every 4 hours via cron, aggregates event counts per user for the last 4 hours, saves to 'stats' table.
- DB: PostgreSQL.
- Client: React app to create/retrieve events.

## Run Instructions (Local)

1. Install Go 1.22+, Node.js.
2. Set up Postgres (or use Docker).
3. `go run cmd/main.go` for API.
4. In /client: `npm install && npm start` (runs on :3000).

## Run Instructions (Docker)

1. `docker-compose up --build`.
2. API at <http://localhost:8080>.
3. For client, build separately or run local.

## Sample Requests

- Create: curl -X POST <http://localhost:8080/events> -d '{"user_id":42,"action":"page_view","metadata":"{"page":"/home"}'}
- Retrieve: curl "<http://localhost:8080/events?user_id=42&start=2025-10-01T00:00:00Z&end=2025-10-10T00:00:00Z>"

## Daily Job Description

The job runs every 4 hours, queries events from the last 4 hours, groups by user_id, counts them, and inserts into 'stats' table with period_start timestamp. It's not strictly daily but produces periodic aggregates that can be summed for daily stats.

## Notes on Optional Parts

Grafana not implemented (nice-to-have). To add: Integrate Prometheus for metrics and Loki for logs.

## .env.example

See .env.example for DB_DSN.
