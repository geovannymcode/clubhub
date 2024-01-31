FROM golang:1.21.6-alpine as builder
# Whois is required for logic the business.
RUN apk add whois
# Create folder app
RUN mkdir /app
# Create and change to the app directory.
WORKDIR /app

COPY run.sh /run.sh
RUN chmod +x /run.sh
# Copy go.mod and if present go.sum.
COPY go.mod go.sum ./
# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

COPY . .
# Build the Go app
RUN go build -v -o main ./cmd/http

# Final stage
FROM alpine
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

CMD ["/app/main"]