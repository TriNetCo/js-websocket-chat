package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestHandleWebSocket(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(handleWebSocket))
	defer server.Close()

	// Create a WebSocket connection to the mock server
	url := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatal("Failed to establish WebSocket connection:", err)
	}
	defer conn.Close()

	// Send a message to the server
	msg := message{Message: "Hello, server!"}
	err = conn.WriteJSON(msg)
	if err != nil {
		t.Fatal("Failed to send message to server:", err)
	}

	// Wait for the server to broadcast the message to all clients
	received := make(chan message)
	go func() {
		for {
			msg := message{}
			err := conn.ReadJSON(&msg)
			if err != nil {
				t.Fatal("Failed to read message from server:", err)
			}
			received <- msg
		}
	}()

	// Check that the message was broadcast correctly
	expected := msg
	actual := <-received
	if expected != actual {
		t.Errorf("Expected message %v but got %v", expected, actual)
	}
}

func TestBroadcaster(t *testing.T) {
	// Create a mock WebSocket client
	client := &websocket.Conn{}

	// Create a broadcast channel
	broadcast := make(chan message)

	// Broadcast a message to the client
	msg := message{Message: "Hello, client!"}
	broadcast <- msg

	// Wait for the message to be sent to the client
	go func() {
		actual := <-broadcast
		expected := msg
		if expected != actual {
			t.Errorf("Expected message %v but got %v", expected, actual)
		}
	}()

	// Close the broadcast channel
	close(broadcast)

	// Check that the client connection was closed
	_, _, err := client.ReadMessage()
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
}
