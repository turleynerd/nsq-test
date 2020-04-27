package main

import (
	"crypto/tls"
	"fmt"
	"log"
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

	msgs := 0
	// Create a new producer so that we can put messages back on the queue
	nsqProducer, err := nsq.NewProducer("nsqd:4150", nsqConfig)
	if err != nil {
		panic(err)
	}

	// Publish a new message every second
	for x := 1; x < 1001; x++ {
		err := nsqProducer.Publish("events", []byte(`hello, world`))
		if err != nil {
			log.Panicln(fmt.Sprintf("Send Failed: %v messages sent.", msgs))
			panic(err)
		}
		msgs = x
		time.Sleep(time.Millisecond)
	}
	time.Sleep(2 * time.Second)
	log.Panicln(fmt.Sprintf("Send Complete: %v messages sent.", msgs))
}
