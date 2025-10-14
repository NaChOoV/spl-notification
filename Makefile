build:
	go build -o cmd/spl-notification cmd/spl-notification/main.go

build-arm64:
	GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o bin/spl-notification-linux-arm64 ./cmd/spl-notification/main.go
build-amd64:
	go build -ldflags="-w -s" -o bin/spl-notification-linux-amd64 ./cmd/spl-notification/main.go

run:
	go run cmd/spl-notification/main.go

test:
	go test ./internal/service/... -v