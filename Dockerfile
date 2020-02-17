FROM golang:1.13

WORKDIR /app

ADD . .
RUN go mod download


#RUN go build -o main .
