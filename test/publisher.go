package main

import (
	"log"
	"time"

	"github.com/bruli/alerter/internal/domain/message"
	"github.com/bruli/pinger/pkg/events"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = nc.Drain()
		nc.Close()
	}()

	resource := "testing"

	msg := events.PingResult{
		Resource: resource,
		Status:   message.FailedStatus,
	}

	data, err := proto.Marshal(&msg)
	if err != nil {
		log.Fatal(err)
	}

	var n int
	i := 5
	for range i {
		err = nc.Publish(events.PingSubjet, data)
		if err != nil {
			log.Fatal(err)
		}
		n++
		time.Sleep(time.Second)
	}

	readyMsg := events.PingResult{
		Resource: resource,
		Status:   "ok",
	}

	data, err = proto.Marshal(&readyMsg)
	if err != nil {
		log.Fatal(err)
	}

	no := 10
	for range no {
		if err = nc.Publish(events.PingSubjet, data); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
	}

	nc.Close()
}
