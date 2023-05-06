package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type message struct {
	Message string `json:"message"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan message)

func main() {
	// Serve the WebSocket endpoint
	http.HandleFunc("/ws", handleWebSocket)

	// Start the message broadcaster
	go broadcaster()

	// Start the server
	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Add the client to the list of connected clients
	clients[conn] = true

	// Handle incoming messages from the client
	for {
		var msg message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			delete(clients, conn)
			break
		}
		// Send the message to the broadcast channel
		broadcast <- msg
	}
}

func broadcaster() {
	for {
		// Get the next message from the broadcast channel
		msg := <-broadcast

		// Send the message to all connected clients
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println(err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
