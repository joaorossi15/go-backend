# Build stage
FROM golang:1.23-alpine AS builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

LABEL maintainer="Joao Rossi <joaorossiborba@gmail.com>"

WORKDIR /app

COPY go.mod ./

# COPY go.sum ./ 2>/dev/null || true

RUN go mod download

COPY . .

RUN go build -o main ./cmd


# Runtime stage

FROM alpine:latest

RUN apk add --no-cache ca-certificates && \
    adduser -D appuser

WORKDIR /app

COPY --from=builder /app/main /app/main

USER appuser

EXPOSE 8080

CMD ["./main"]
