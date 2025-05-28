FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o job_solution ./cmd/api/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/job_solution .
COPY --from=builder /app/migrations ./migrations

ENV GIN_MODE=release

EXPOSE 8080

CMD ["./job_solution"]