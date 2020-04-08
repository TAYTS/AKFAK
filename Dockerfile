FROM golang:latest

RUN go get google.golang.org/grpc

RUN apt-get update -y && apt-get install vim -y

WORKDIR /go/src/AKFAK

COPY . .

CMD [ "go", "run", "/go/src/AKFAK/broker/server/server.go" ]
