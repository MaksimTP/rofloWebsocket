package endpoints

import (
	"fmt"
	"log"
	"main/server/socket_server"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	writeWait  = 10 * time.Second
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

var server = &socket_server.Server{Subscriptions: make(socket_server.Subscription)}

func HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	log.Println(err)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed upgrading connection"))
		return
	}
	defer conn.Close()

	// create new client id
	clientID := uuid.New().String()
	server.Subscribe(conn, clientID, socket_server.BaseTopic)
	server.Send(conn, fmt.Sprintf("Server: welcome message, id: %s", clientID))

	done := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go writePump(conn, clientID, done, &wg)

	wg.Wait()

}

func ChangeSomething(w http.ResponseWriter, r *http.Request) {
	log.Println("something changed")
	data := []byte(fmt.Sprintf("changing something, %s", time.Now()))
	server.Publish(socket_server.BaseTopic, data)
	w.Write(data)
}

func writePump(conn *websocket.Conn, clientID string, done <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait))

			if err != nil {
				server.RemoveClient(clientID)
				return
			}
		case <-done:
			return
		}
	}
}
