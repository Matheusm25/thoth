package main

import (
	"fmt"
	"time"

	"github.com/matheusm25/thoth/src/modules/broker"
)

func main() {
	broker := broker.NewBroker()

	subscriber := broker.Subscribe("websocket", func(msg string) {
		fmt.Printf("Received: %v\n", msg)
	})

	broker.Publish("websocket", "First message to first listener")
	broker.Publish("websocket", "Second message to fisrt listener")
	broker.Publish("websocket-fake", "This message should not show")

	time.Sleep(2 * time.Second)

	broker.Unsubscribe("websocket", subscriber)

	broker.Publish("websocket", "This message won't be received.")

	time.Sleep(time.Second)
}
