build:
	go build -o cmd/spl-notification cmd/spl-notification/main.go

run:
	go run cmd/spl-notification/main.go

test:
	go test ./internal/service/... -v