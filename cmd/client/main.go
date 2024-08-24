package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"goapp/internal/pkg/watcher"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gorilla/websocket"
)

func main() {
	// Parse flags
	var (
		address     = flag.String("a", "localhost:8080", "Host address for connecting to WebSocket server")
		connections = flag.Int("n", 1, "Number of parallel WebSocket connections to open")
	)
	flag.Parse()

	// URL of the WS server
	u, err := url.Parse(fmt.Sprintf("ws://%v/goapp/ws", *address))
	if err != nil {
		log.Fatal("Invalid server URL:", err)
	}

	// Channel to signal termination
	stopChan := make(chan struct{})

	// Handle interrupt signals
	go handleInterrupt(stopChan)

	// Wait for all connections
	var wg sync.WaitGroup
	wg.Add(*connections)

	// Open specified number of WS connections
	for i := 0; i < *connections; i++ {
		go func(id int) {
			defer wg.Done()
			connectWebSocket(u, id, stopChan)
		}(i)
	}

	// Wait for all connections to complete
	wg.Wait()
}

func connectWebSocket(u *url.URL, id int, stopChan <-chan struct{}) {
	// Connect to server
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Printf("[conn #%d] Failed to connect: %v\n", id, err)
		return
	}

	// Close the WebSocket connection gracefully
	defer func() {
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		conn.Close()
	}()

	for {
		select {
		case <-stopChan:
			// Exit and close the connection
			return
		default:
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
}

// handleInterrupt waits until a SIGINT or SIGTERM signal is received and closes
// the stopChan channel.
func handleInterrupt(stopChan chan struct{}) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	close(stopChan)
}
