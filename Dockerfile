FROM golang:1.17.3 AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o readimension -v -ldflags '-s -w' *.go

FROM debian:stable-slim
RUN apt-get update
RUN apt-get install -y sqlite3
RUN apt-get clean
WORKDIR /data/readimension
COPY config.json.sample config.json
WORKDIR /app
COPY --from=builder /app/readimension .
CMD ["./readimension", "--env", "development", "--path", "/data/readimension"]
