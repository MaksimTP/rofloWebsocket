package socket_server

import (
	"sync"

	"github.com/gorilla/websocket"
)

var BaseTopic = "dbOnChange"

type Subscription map[string]Client

type Client map[string]*websocket.Conn

type Server struct {
	Subscriptions Subscription
}

func (s *Server) Send(conn *websocket.Conn, message string) {
	conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (s *Server) SendWithWait(conn *websocket.Conn, message string, wg *sync.WaitGroup) {
	conn.WriteMessage(websocket.TextMessage, []byte(message))

	wg.Done()
}

func (s *Server) RemoveClient(clientID string) {
	for _, client := range s.Subscriptions {
		delete(client, clientID)
	}
}

func (s *Server) Publish(topic string, message []byte) {
	if _, exist := s.Subscriptions[topic]; !exist {
		return
	}

	client := s.Subscriptions[topic]

	var wg sync.WaitGroup

	for _, conn := range client {
		wg.Add(1)

		go s.SendWithWait(conn, string(message), &wg)
	}

	wg.Wait()
}

func (s *Server) Subscribe(conn *websocket.Conn, clientID string, topic string) {
	if _, exist := s.Subscriptions[topic]; exist {
		client := s.Subscriptions[topic]

		if _, subbed := client[clientID]; subbed {
			return
		}

		client[clientID] = conn
		return
	}

	newClient := make(Client)
	s.Subscriptions[topic] = newClient

	s.Subscriptions[topic][clientID] = conn
}

func (s *Server) Unsubscribe(clientID string, topic string) {
	if _, exist := s.Subscriptions[topic]; exist {
		client := s.Subscriptions[topic]

		delete(client, clientID)
	}
}
