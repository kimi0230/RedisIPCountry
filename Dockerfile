FROM golang:1.16.6-alpine3.14

ENV GO111MODULE on

WORKDIR /src/app

ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY . .




