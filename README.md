# How would I improve my SOLUTIUON if I had time

1) I would move frontend/ folder to a separate repo, for sure.
2) Add delete handler and repo func, and also Decrement for Prometheus metric.
3) Swagger

## Go Activity Tracker

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

## Tests

- go test (or make test)
- go test -v -tags=integration ./internal/repository (make integration) (integration tests)
- make benchmark
- linters: golangci-lint run --config .golangci.yml
- make all-tests

## Daily Job Description

The job runs every 4 hours, queries events from the last 4 hours, groups by user_id, counts them, and inserts into 'stats' table with period_start timestamp. It's not strictly daily but produces periodic aggregates that can be summed for daily stats.

## Metrics

- We can check on port: <http://localhost:8080/metrics>

## Prometheus

- Open in your browser: <http://localhost:9090/>

## Grafana

- Open Grafana: <http://localhost:3000/>
- Default login: admin / admin
- You’ll be prompted to change the password
- Go to Settings → Data Sources → Add data source
- Choose Prometheus: in the URL field, enter: <http://prometheus:9090>
- Click Save & Test — should show Data source is working

## .env.example

See .env.example for DB_DSN
