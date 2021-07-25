FROM golang:1.14
MAINTAINER Tae ho Kim <sjabber91@gmail.com>

ENV DEBIAN_FRONTEND=nointeractive

RUN mkdir -p /api
WORKDIR /api

COPY ./redteam/go.mod .
COPY ./redteam/go.mod .
RUN go mod download

COPY ./redteam .
RUN go build -o ./redteam ./main.go

EXPOSE 5000

ENTRYPOINT ["./redteam"]
