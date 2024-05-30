FROM golang:1.21.1-alpine3.17

WORKDIR /app

COPY . .

COPY ./.env /app/.env

RUN go mod download && \
    go build -o main

CMD ["./main"]