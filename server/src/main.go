package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/matheusm25/thoth/src/modules/websocket"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	broadcast := make(chan string)
	websocket.InitWebsocketServer(func(msg string) {
		broadcast <- msg
	})
}
