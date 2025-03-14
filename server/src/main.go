package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/matheusm25/thoth/src/app"
	"github.com/matheusm25/thoth/src/modules/broker"
	"github.com/matheusm25/thoth/src/modules/db/models"
	websocket_server "github.com/matheusm25/thoth/src/modules/websocket"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbMessages, err := models.GetAllMessages()
	if err != nil {
		log.Fatalf("Error getting messages: %v", err)
	}

	previousMessages := make([]*broker.Message, 0)
	for _, dbMessage := range dbMessages {
		previousMessages = append(previousMessages, &broker.Message{
			ID:          dbMessage.ID,
			Topic:       dbMessage.Topic,
			MessageType: "PUBLISH",
			Payload:     dbMessage.Payload,
		})
	}

	applicationBroker := broker.NewBroker()
	applicationBroker.BootstrapMessagesOnHold(previousMessages)
	application := app.NewApp(applicationBroker)

	websocket_server.InitWebsocketServer(application.HandleNewMessage, application.HandleNewConnection)
}
