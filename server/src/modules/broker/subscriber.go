package broker

import (
	"fmt"
)

type Subscriber struct {
	Channel             chan string
	isProcessingMessage bool
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
			s.isProcessingMessage = true
			executer(msg)
			s.isProcessingMessage = false
		case <-s.Unsubscribe:
			fmt.Println("Unsubscribed.")
			return
		}
	}
}
