FROM golang:1.25.2-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-w -s" -o main ./cmd/spl-notification/main.go

FROM golang:1.25.2-alpine
WORKDIR /app
COPY --from=build /app/main .

# Turso Migrations
COPY --from=build /app/migrations/*.sql /app/migrations/

EXPOSE 8000

CMD ["./main"]