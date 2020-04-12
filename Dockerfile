FROM golang:latest

WORKDIR /go/src/AKFAK

RUN go mod init AKFAK && \
  go mod tidy && \
  apt-get update -y && \
  apt-get install vim -y

COPY . .

RUN go install ./...
