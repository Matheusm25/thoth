package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/matheusm25/thoth/src/app"
	"github.com/matheusm25/thoth/src/modules/broker"
	websocket_server "github.com/matheusm25/thoth/src/modules/websocket"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	broker := broker.NewBroker()
	application := app.NewApp(broker)

	websocket_server.InitWebsocketServer(application.HandleNewMessage, application.HandleNewConnection)

}
