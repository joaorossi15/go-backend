# Start from golang:1.12-alpine base image
FROM golang:1.22-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

LABEL maintainer="Joao Rossi <joaorossiborba@gmail.com>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
