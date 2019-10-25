package main

import (
	"crypto/tls"
	"time"

	nsq "github.com/nsqio/go-nsq"
)

func main() {
	nsqConfig := nsq.NewConfig()
	nsqConfig.MaxInFlight = 100
	nsqConfig.TlsV1 = true
	nsqConfig.TlsConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	// Create a new producer so that we can put messages back on the queue
	nsqProducer, err := nsq.NewProducer("nsqd:4150", nsqConfig)
	if err != nil {
		panic(err)
	}

	// Publish a new message every second
	for {
		err := nsqProducer.Publish("events", []byte(`hello, world`))
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}

}
