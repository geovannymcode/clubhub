FROM golang:1.27.3-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /clubhub ./cmd/http

FROM alpine:3.24
RUN apk add --no-cache ca-certificates && addgroup -S app && adduser -S -G app app
WORKDIR /app
COPY --from=builder /clubhub /app/clubhub
USER app:app
EXPOSE 8001
CMD ["/app/clubhub"]
