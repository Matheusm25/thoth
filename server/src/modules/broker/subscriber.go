package broker

import (
	"fmt"
)

type Subscriber struct {
	Channel             chan string
	IsProcessingMessage bool
	Unsubscribe         chan bool
}

type MessageConsumerFunction func(string)

func (s *Subscriber) OnMessage(executer MessageConsumerFunction) {
	for {
		select {
		case msg, ok := <-s.Channel:
			if !ok {
				fmt.Println("Subscriber channel closed.")
				return
			}
			s.IsProcessingMessage = true
			executer(msg)
		case <-s.Unsubscribe:
			fmt.Println("Unsubscribed.")
			return
		}
	}
}
