package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	nsq "github.com/nsqio/go-nsq"
)

type counter struct {
	counter    int
	counterTwo int
	lock       sync.Mutex
}

func main() {
	nsqConfig := nsq.NewConfig()
	nsqConfig.MaxInFlight = 100
	nsqConfig.TlsV1 = true
	nsqConfig.TlsConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	ticker := time.NewTicker(10 * time.Second)
	done := make(chan bool)

	msgConsumed := counter{}

	// Create a new consumer on the events topic
	nsqConsumer, err := nsq.NewConsumer("events", "events-consumer", nsqConfig)
	if err != nil {
		panic(err)
	}

	// Tell it which handler to run when a message comes in
	nsqConsumer.AddConcurrentHandlers(nsq.HandlerFunc(func(msg *nsq.Message) error {
		msgConsumed.lock.Lock()
		defer msgConsumed.lock.Unlock()
		fmt.Printf("Message: %+v \n", msg)
		msgConsumed.counterTwo++
		msg.Finish()
		msgConsumed.counter++
		return nil
	}), 50)

	// Print out our current count
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				log.Println(fmt.Sprintf("Consumed Messages: %v Finished Messages: %v at %v", msgConsumed.counterTwo, msgConsumed.counter, t))
			}
		}
	}()

	// Connect to nsqlookupd
	err = nsqConsumer.ConnectToNSQLookupd("nsqlookupd:4161")
	if err != nil {
		panic(err)
	}

	// We listen for stopping signals, and attempt to shut down gracefully.
	shutdown := make(chan os.Signal, 2)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGKILL)

	// We listen for a couple of messages in order to exit. If receive an interrupt or kill signal, we stop the
	// metrics server, and start to drain NSQ. Once NSQ is drained, we return from the main func.
	for {
		select {
		case <-nsqConsumer.StopChan:
			return
		case <-shutdown:
			nsqConsumer.Stop()
			ticker.Stop()
		}
	}
}
