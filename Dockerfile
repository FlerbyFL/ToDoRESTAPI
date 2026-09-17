# --- Build stage ---
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server ./cmd/server

# --- Final stage ---
FROM alpine:3.20

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/bin/server .

USER appuser
EXPOSE 8080

ENTRYPOINT ["./server"]