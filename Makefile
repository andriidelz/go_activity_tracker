test:
	go test ./... -v -timeout 30s

integration:
	@echo "Running integration tests..."
	TEST_DATABASE_DSN="postgres://postgres:postgres@localhost:5432/activity_test?sslmode=disable" \
	go test -v -tags=integration ./internal/repository

all-tests: test integration
	@echo "All tests completed ✅"

metrics:
	go run ./cmd/main.go

build:
	mkdir -p bin
	GOEXPERIMENT=greenteagc CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/main ./cmd

clean:
	rm -rf ./bin/main