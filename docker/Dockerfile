FROM golang:1.16-alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 5001

RUN go build ./cmd/api/main.go

ENV config=docker

CMD ["./main"]