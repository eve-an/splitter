FROM golang:1.25 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o splitter ./cmd/server

# Run stage
FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/splitter .
EXPOSE 8080
ENTRYPOINT ["./splitter"]

