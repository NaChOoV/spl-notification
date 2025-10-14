FROM golang:1.25.2-alpine AS build
RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN ARCH=$(uname -m) && \
    if [ "$ARCH" = "aarch64" ]; then \
    GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o main ./cmd/spl-notification/main.go; \
    else \
    go build -ldflags="-w -s" -o main ./cmd/spl-notification/main.go; \
    fi

FROM golang:1.25.2-alpine
WORKDIR /app
COPY --from=build /app/main .

# Turso Migrations
COPY --from=build /app/migrations/*.sql /app/migrations/

EXPOSE 8000

CMD ["./main"]