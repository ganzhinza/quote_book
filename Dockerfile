FROM golang:1.24 AS builder

WORKDIR /app

COPY . .
RUN go mod download

WORKDIR /app/cmd/

ENV CONFIG_PATH=../config/config.json

RUN GOOS=linux go build -o app

EXPOSE 8080
CMD ["./app"]