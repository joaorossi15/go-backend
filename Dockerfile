# Build stage
FROM golang:1.23-alpine AS builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

LABEL maintainer="Joao Rossi <joaorossiborba@gmail.com>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd

EXPOSE 8080

CMD ["./main"]
