FROM golang:1.13.8-alpine3.11

WORKDIR /app

ADD . .

RUN go mod download