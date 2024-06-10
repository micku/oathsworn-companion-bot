FROM arm64v8/golang:1.22.4 AS builder

RUN apt-get update && apt-get install -y gcc-aarch64-linux-gnu

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64 go build -o /oath-bot

FROM arm64v8/debian:12-slim

WORKDIR /

COPY --from=builder /oath-bot /oath-bot

RUN apt-get update && apt-get install -y ca-certificates

ENV TG_TOKEN=""
CMD ["/oath-bot"]
