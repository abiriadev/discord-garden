FROM golang:1.22.1-alpine3.19

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

ENTRYPOINT go run .
