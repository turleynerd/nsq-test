FROM golang:1.12

COPY . /go/src/github.com/RedVentures/nsq-tls

RUN go install github.com/RedVentures/nsq-tls/cmd/consumer
RUN go install github.com/RedVentures/nsq-tls/cmd/producer