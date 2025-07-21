package main

import (
	"context"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
)

func main() {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: "pulsar://localhost:6650",
	})
	defer client.Close()

	if err != nil {
		panic(err)
	}

	for i := 0; i < 100; i++ {
		producer, err := client.CreateProducer(pulsar.ProducerOptions{
			Topic: fmt.Sprintf("persistent://tenant-%d/default/event", 1),
		})
		defer producer.Close()

		if err != nil {
			panic(err)
		}
		producer.Send(context.Background(), &pulsar.ProducerMessage{Payload: []byte(fmt.Sprintf("hello from producer %d", 1))})
	}

	for i := 0; i < 5; i++ {
		producer, err := client.CreateProducer(pulsar.ProducerOptions{
			Topic: fmt.Sprintf("persistent://tenant-%d/default/event", 2),
		})
		defer producer.Close()

		if err != nil {
			panic(err)
		}
		producer.Send(context.Background(), &pulsar.ProducerMessage{Payload: []byte(fmt.Sprintf("hello from producer %d", i))})
	}
}
