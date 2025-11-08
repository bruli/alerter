package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/bruli/alerter/internal/domain/message"
	infranats "github.com/bruli/alerter/internal/infra/nats"
	"github.com/nats-io/nats.go"
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

	msg := infranats.Message{
		Resource: resource,
		Status:   message.FailedStatus,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}

	var n int
	i := 5
	for range i {
		err = nc.Publish(infranats.PingSubject, data)
		if err != nil {
			log.Fatal(err)
		}
		n++
		time.Sleep(time.Second)
	}

	readyMsg := infranats.Message{
		Resource: resource,
		Status:   "ok",
	}

	data, err = json.Marshal(readyMsg)
	if err != nil {
		log.Fatal(err)
	}

	no := 10
	for range no {
		if err = nc.Publish(infranats.PingSubject, data); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
	}

	nc.Close()
}
