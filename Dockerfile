FROM golang:1.26.1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
ENV GOPROXY=https://goproxy.io,direct
RUN go mod download
COPY . .
RUN go build -o server ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 4000
CMD ["./server"]