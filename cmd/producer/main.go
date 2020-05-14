package main

import (
	"crypto/tls"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	shutdown := make(chan os.Signal, 2)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// We listen for a couple of messages in order to exit. If receive an interrupt or kill signal, we stop the
	// metrics server, and start to drain NSQ. Once NSQ is drained, we return from the main func.
	go func() {
		for {
			select {
			case <-shutdown:
				log.Println("SIGTERM REC: Shutting down.")
				nsqProducer.Stop()
				log.Printf("Send Complete: %v messages sent.", msgs)
			}
		}
	}()

	// Publish a new message every Millisecond
	for x := 1; x < 1001; x++ {
		err := nsqProducer.Publish("events", []byte(`hello, world`))
		if err != nil {
			log.Printf("Send Failed: %v messages sent.", msgs)
			panic(err)
		}
		msgs = x
		time.Sleep(time.Millisecond)
	}

	for x := 1; x < 1001; x++ {
		err := nsqProducer.Publish("events2", []byte(`hello, world`))
		if err != nil {
			log.Printf("Send Failed: %v messages sent.", msgs)
			panic(err)
		}
		msgs = x
		time.Sleep(time.Millisecond)
	}

	time.Sleep(2 * time.Second)
	log.Printf("Send Complete: %v messages sent.", msgs)

}
