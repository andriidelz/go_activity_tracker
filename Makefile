test:
	go test ./... -v -timeout 30s

build:
	GOEXPERIMENT=greenteagc CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/main ./cmd