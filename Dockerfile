FROM golang:1.12

COPY . /go/src/github.com/turleynerd/nsq-test

RUN go install github.com/turleynerd/nsq-test/cmd/consumer
RUN go install github.com/turleynerd/nsq-test/cmd/producer
