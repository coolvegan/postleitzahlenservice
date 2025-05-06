# Build-Stage
FROM golang:1.24.2 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./cmd/server

# Runtime-Stage
FROM alpine:latest
WORKDIR /app/server
RUN apk --no-cache add ca-certificates
RUN mkdir -p /app/server/certs 
COPY --from=builder /app/server /app/server
COPY --from=builder /app/certs /app/server/certs/
EXPOSE 50051
ENTRYPOINT ["/app/server"]