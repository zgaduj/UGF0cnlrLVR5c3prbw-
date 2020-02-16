FROM golang:1.13

#RUN apt update && apk upgrade && \
#    apk add --no-cache bash git openssh

#RUN apt update && apt install y git

# Set the Current Working Directory inside the container
WORKDIR /app

ADD go.mod go.sum ./

RUN go mod download

ADD . .

#RUN go build -o main .
