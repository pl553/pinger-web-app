# syntax=docker/dockerfile:1

FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /pinger-app

CMD "/pinger-app" "--backend_url" "http://localhost:$BACKEND_PORT" "--ping_period_sec" "$PING_PERIOD_SEC"
