package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"goapp/internal/pkg/watcher"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

const serverURL = "ws://localhost:8080/goapp/ws"

func main() {
	// Parse flags
	var (
		connections = flag.Int("n", 1, "Number of parallel WebSocket connections to open")
	)
	flag.Parse()
	_ = connections

	// URL of the WS server
	u, err := url.Parse(serverURL)
	if err != nil {
		log.Fatal("Invalid server URL:", err)
	}

	// TODO: Implement parallel connections
	// Connect to WS
	connectWebSocket(u, 0)
}

func connectWebSocket(u *url.URL, id int) {
	// Connect to server
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Printf("[conn #%d] Failed to connect: %v\n", id, err)
		return
	}
	defer conn.Close()

	for {
		// Read message from server
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("[conn #%d] Failed to read message: %v\n", id, err)
			return
		}

		// Parse message as Counter
		var counter watcher.Counter
		err = json.Unmarshal(message, &counter)
		if err != nil {
			fmt.Printf("[conn #%d] Failed to parse message: %v\n", id, err)
			return
		}
		fmt.Printf("[conn #%d] iteration: %v, value: %v\n", id, counter.Iteration, counter.Value)
	}
}
