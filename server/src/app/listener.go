package app

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/matheusm25/thoth/src/modules/broker"
	IDUtils "github.com/matheusm25/thoth/src/utils/id"
)

type App struct {
	broker  *broker.Broker
	clients map[string]*ManagedConncetion
	mutex   sync.Mutex
}

type ManagedConncetion struct {
	ID     string
	Conn   *websocket.Conn
	Topics map[string]*broker.Subscriber
	Active bool
}

func NewApp(broker *broker.Broker) *App {
	return &App{
		broker:  broker,
		clients: make(map[string]*ManagedConncetion),
	}
}

func (app *App) HandleNewConnection(conn *websocket.Conn) {
	id := IDUtils.GenerateID()
	app.clients[id] = &ManagedConncetion{
		ID:     id,
		Conn:   conn,
		Topics: make(map[string]*broker.Subscriber),
		Active: true,
	}
	conn.WriteJSON(map[string]string{"id": id})
}

func (app *App) HandleNewMessage(msg string) {
	message := broker.Message{}
	err := json.Unmarshal([]byte(msg), &message)
	if err != nil {
		fmt.Printf("Error unmarshalling message: %v\n", err)
		return
	}

	_, err = message.Validate()
	if err != nil {
		fmt.Printf("Invalid message: %v\n", err)
		return
	}

	switch strings.ToUpper(message.MessageType) {
	case "PUBLISH":
		app.handlePublishMessage(&message)
	case "SUBSCRIBE":
		app.handleSubscribeMessage(&message)
	case "UNSUBSCRIBE":
		app.handleUnsubscribeMessage(&message)
	case "ACKNOWLEDGE":
		app.handleAcknowledgeMessage(&message)
	case "UNACKNOWLEDGE":
		app.handleUnacknowledgeMessage(&message)
	default:
		fmt.Printf("Invalid message type: %v\n", message.MessageType)
	}
}

func (app *App) handlePublishMessage(message *broker.Message) {
	app.mutex.Lock()
	defer app.mutex.Unlock()

	fmt.Printf("Publishing message: %v\n", message)

	app.broker.Publish(message)
}

func (app *App) handleSubscribeMessage(message *broker.Message) {
	app.mutex.Lock()
	defer app.mutex.Unlock()

	fmt.Printf("Subscribing message: %v\n", message)

	client := app.clients[message.ID]
	if client == nil {
		fmt.Printf("Client not found: %v\n", message.ID)
		return
	}

	if client.Topics[message.Topic] != nil {
		fmt.Printf("Client already subscribed to topic: %v\n", message.Topic)
		return
	}

	client.Topics[message.Topic] = app.broker.Subscribe(message.Topic, func(msg string) {
		app.clients[message.ID].Conn.WriteJSON(map[string]string{
			"Topic":   message.Topic,
			"Payload": msg,
		})
	})
}

func (app *App) handleUnsubscribeMessage(message *broker.Message) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	fmt.Printf("Unsubscribing message: %v\n", message)

	client := app.clients[message.ID]
	if client == nil {
		fmt.Printf("Client not found: %v\n", message.ID)
		return
	}

	if client.Topics[message.Topic] == nil {
		fmt.Printf("Client not subscribed to topic: %v\n", message.Topic)
		return
	}

	app.broker.Unsubscribe(message.Topic, client.Topics[message.Topic])
	delete(client.Topics, message.Topic)
}

func (app *App) handleAcknowledgeMessage(message *broker.Message) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	fmt.Printf("Acknowledging message: %v\n", message)

	client := app.clients[message.ID]
	if client == nil {
		fmt.Printf("Client not found: %v\n", message.ID)
		return
	}

	if client.Topics[message.Topic] == nil {
		fmt.Printf("Client not subscribed to topic: %v\n", message.Topic)
		return
	}

	client.Topics[message.Topic].IsProcessingMessage = false
}

func (app *App) handleUnacknowledgeMessage(message *broker.Message) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	fmt.Printf("Unacknowledging message: %v\n", message)

	client := app.clients[message.ID]
	if client == nil {
		fmt.Printf("Client not found: %v\n", message.ID)
		return
	}

	if client.Topics[message.Topic] == nil {
		fmt.Printf("Client not subscribed to topic: %v\n", message.Topic)
		return
	}

	client.Topics[message.Topic].IsProcessingMessage = false

	message.MessageType = "PUBLISH"
	app.broker.SetMessageOnHold(message)
}
