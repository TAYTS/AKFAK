FROM golang:latest

RUN apt-get update -y && \
  apt-get install vim -y

WORKDIR /go/src/AKFAK

COPY ["go.mod", "go.sum", "./"]

RUN go mod download


COPY . .

RUN go install ./...
