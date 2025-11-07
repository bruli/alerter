package main

import (
	"encoding/json"
	"fmt"
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

	resource := fmt.Sprintf("testing at %s", time.Now().String())

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
	}
	nc.Close()
}
