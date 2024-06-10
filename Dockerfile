# syntax=docker/dockerfile:1

FROM golang:1.22.4

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /oath-bot

ENV TG_TOKEN=""

CMD ["/oath-bot"]
