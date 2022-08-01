# syntax=docker/dockerfile:1
FROM golang:1.19-rc-bullseye

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY config.txt ./

RUN go build -o /random

EXPOSE 8080

CMD [ "/random" ]
