# Stage 1: Build the go app
FROM golang:1.14 AS builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY go.sum .
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o main .

# Stage 2: Run the built file
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=builder /build/main .
EXPOSE 8080
CMD ["/app/main"]
