# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /opt/src

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /opt/bin/server

# TODO set environment variables?

EXPOSE 8080

CMD [ "/opt/bin/server" ]
