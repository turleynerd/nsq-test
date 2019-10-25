package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	nsq "github.com/nsqio/go-nsq"
)

func main() {
	nsqConfig := nsq.NewConfig()
	nsqConfig.MaxInFlight = 100
	nsqConfig.TlsV1 = true
	nsqConfig.TlsConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	// Create a new consumer on the events topic
	nsqConsumer, err := nsq.NewConsumer("events", "events-consumer", nsqConfig)
	if err != nil {
		panic(err)
	}

	// Tell it which handler to run when a message comes in
	nsqConsumer.AddConcurrentHandlers(nsq.HandlerFunc(func(msg *nsq.Message) error {
		fmt.Println(string(msg.Body))
		msg.Finish()
		return nil
	}), 50)

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
		}
	}
}
