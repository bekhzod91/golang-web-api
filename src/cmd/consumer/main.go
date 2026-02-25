package main

import (
	"fmt"
	pulsar2 "github.com/apache/pulsar-client-go/pulsar"
	"github.com/myproject/api/pkg/pulsar"
	"time"
)

func consumer1(c pulsar2.ConsumerMessage) {
	for true {
		fmt.Printf("Consumer %s received a message, msgId: %v, content: '%s'\n",
			c.Consumer.Name(), c.Message.ID(), string(c.Message.Payload()))
		time.Sleep(1 * time.Second)
	}
}

func main() {
	funcs := pulsar.Topics{"event": consumer1}
	tenants := []string{"tenant1", "tenant2"}
	err := pulsar.NewConsumer(tenants, funcs)
	if err != nil {
		panic(fmt.Sprintf("consumer %s", err.Error()))
	}
}
