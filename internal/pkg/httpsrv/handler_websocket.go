package httpsrv

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"goapp/internal/pkg/watcher"

	"github.com/gorilla/websocket"
)

// checkSameOrigin returns true if the origin set and equal to the request host.
func checkSameOrigin(r *http.Request) bool {
	origin := r.Header["Origin"]
	if len(origin) == 0 {
		return true
	}
	u, err := url.Parse(origin[0])
	if err != nil {
		return false
	}
	return strings.EqualFold(u.Host, r.Host)
}

func (s *Server) handlerWebSocket(w http.ResponseWriter, r *http.Request) {
	// Create and start a watcher.
	var watch = watcher.New()
	if err := watch.Start(); err != nil {
		s.error(w, http.StatusInternalServerError, fmt.Errorf("failed to start watcher: %w", err))
		return
	}
	defer watch.Stop()

	s.addWatcher(watch)
	defer s.removeWatcher(watch)

	// Start WS.
	var upgrader = websocket.Upgrader{
		// Validate request origin to prevent CSRF attacks
		CheckOrigin: checkSameOrigin,
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.error(w, http.StatusInternalServerError, fmt.Errorf("failed to upgrade connection: %w", err))
		return
	}
	defer func() { _ = c.Close() }()

	log.Printf("websocket started for watcher %s\n", watch.GetWatcherId())
	defer func() {
		log.Printf("websocket stopped for watcher %s\n", watch.GetWatcherId())
	}()

	// Read done.
	readDoneCh := make(chan struct{})

	// All done.
	doneCh := make(chan struct{})
	defer close(doneCh)

	go func() {
		defer close(readDoneCh)
		for {
			select {
			default:
				_, p, err := c.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
						log.Printf("failed to read message: %v\n", err)
					}
					return
				}
				var m watcher.CounterReset
				if err := json.Unmarshal(p, &m); err != nil {
					log.Printf("failed to unmarshal message: %v\n", err)
					continue
				}
				watch.ResetCounter()
			case <-doneCh:
				return
			case <-s.quitChannel:
				return
			}
		}
	}()

	for {
		select {
		case cv := <-watch.Recv():
			data, _ := json.Marshal(cv)
			err = c.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("failed to write message: %v\n", err)
				}
				return
			}
		case <-readDoneCh:
			return
		case <-s.quitChannel:
			return
		}
	}
}
