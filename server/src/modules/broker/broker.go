package broker

import (
	"fmt"
	"slices"
	"sync"
	"time"

	SliceUtils "github.com/matheusm25/thoth/src/utils/slices"
)

type Broker struct {
	subscribers    map[string][]*Subscriber
	messagesOnHold map[string][]*Message
	mutex          sync.Mutex
}

func NewBroker() *Broker {
	return &Broker{
		subscribers:    make(map[string][]*Subscriber),
		messagesOnHold: make(map[string][]*Message),
	}
}

func (b *Broker) Subscribe(topic string, executer MessageConsumerFunction) *Subscriber {
	b.mutex.Lock()

	subscriber := &Subscriber{
		Channel:     make(chan string),
		Unsubscribe: make(chan bool),
	}

	go subscriber.OnMessage(executer)

	b.subscribers[topic] = append(b.subscribers[topic], subscriber)

	if messages, found := b.messagesOnHold[topic]; found && len(messages) > 0 {
		b.mutex.Unlock()
		for _, msg := range messages {
			b.Publish(msg)
			b.messagesOnHold[topic] = append(messages[:0], messages[1:]...)
		}

		return subscriber
	}

	b.mutex.Unlock()
	return subscriber
}

func (b *Broker) Unsubscribe(topic string, subscriber *Subscriber) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if subscribers, found := b.subscribers[topic]; found {
		for i, sub := range subscribers {
			if sub == subscriber {
				close(sub.Channel)
				b.subscribers[topic] = append(subscribers[:i], subscribers[i+1:]...)
				return
			}
		}
	}
}

func (b *Broker) Publish(message *Message) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	filteredSubscribers := SliceUtils.Filter(b.subscribers[message.Topic], func(s *Subscriber) bool {
		return !s.IsProcessingMessage
	})

	if len(filteredSubscribers) > 0 {
		select {
		case filteredSubscribers[0].Channel <- message.Payload:
		case <-time.After(time.Second):
			fmt.Printf("Subscriber slow. Unsubscribing from topic: %s\n", message.Topic)
			b.Unsubscribe(message.Topic, filteredSubscribers[0])
		}
	} else {
		b.SetMessageOnHold(message)
	}
}

func (b *Broker) SetMessageOnHold(message *Message) {
	b.messagesOnHold[message.Topic] = append(b.messagesOnHold[message.Topic], message)
	fmt.Printf("No subscribers found for topic: %s. Messages on hold: %v\n", message.Topic, len(b.messagesOnHold[message.Topic]))

	if len(b.messagesOnHold[message.Topic]) == 1 {
		go b.MessagesOnHoldRoutine(message.Topic)
	}
}

func (b *Broker) MessagesOnHoldRoutine(topic string) {
	for {
		if len(b.messagesOnHold[topic]) == 0 {
			return
		}

		filteredSubscribers := SliceUtils.Filter(b.subscribers[topic], func(s *Subscriber) bool {
			return !s.IsProcessingMessage
		})

		if len(filteredSubscribers) > 0 {
			b.Publish(b.messagesOnHold[topic][0])
			b.messagesOnHold[topic] = slices.Delete(b.messagesOnHold[topic], 0, 1)
		}

		time.Sleep(time.Second)
	}
}

func TestBroker() {
	broker := NewBroker()

	subscriber := broker.Subscribe("websocket", func(msg string) {
		fmt.Printf("Received: %v\n", msg)
	})

	broker.Publish(&Message{
		ID:          "1",
		Topic:       "websocket",
		MessageType: "PUBLISH",
		Payload:     "First message to fisrt listener",
	})
	broker.Publish(&Message{
		ID:          "1",
		Topic:       "websocket",
		MessageType: "PUBLISH",
		Payload:     "Second message to fisrt listener",
	})
	broker.Publish(&Message{
		ID:          "1",
		Topic:       "websocket-fake",
		MessageType: "PUBLISH",
		Payload:     "This message should not show",
	})

	time.Sleep(2 * time.Second)

	broker.Unsubscribe("websocket", subscriber)

	broker.Publish(&Message{
		ID:          "1",
		Topic:       "websocket",
		MessageType: "PUBLISH",
		Payload:     "This message won't be received",
	})

	time.Sleep(time.Second)
}
