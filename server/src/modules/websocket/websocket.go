package websocket

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

type HandleMessagesFunction func(string)

var Clients = make(map[*websocket.Conn]bool)

func InitWebsocketServer(handleMessages HandleMessagesFunction) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebsocketConnection(w, r, handleMessages)
	})

	fmt.Printf("Server started on port %v", os.Getenv("PORT"))
	err := http.ListenAndServe(fmt.Sprintf(":%v", os.Getenv("PORT")), nil)
	if err != nil {
		panic("Error starting server: " + err.Error())
	}
}

func handleWebsocketConnection(w http.ResponseWriter, r *http.Request, handleMessages HandleMessagesFunction) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	Clients[conn] = true

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			delete(Clients, conn)
			return
		}

		if messageType == websocket.TextMessage {
			handleMessages(string(msg))
		} else {
			fmt.Printf("Type message '%v' not implemented yet,\n", messageType)
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func TestWebsocket() {
	var broadcast = make(chan string)

	InitWebsocketServer(func(msg string) {
		broadcast <- msg
	})

	go func() {
		for {
			msg := <-broadcast

			for client := range Clients {
				err := client.WriteJSON(msg)
				if err != nil {
					fmt.Println(err)
					client.Close()
					delete(Clients, client)
				}
			}
		}
	}()

}
