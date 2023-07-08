package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     checkOrigin,
	}
)

type Manager struct {
	clients ClientList
	sync.RWMutex

	handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		clients: make(ClientList),

		handlers: make(map[string]EventHandler),
	}

	m.setupEventHandlers()
	return m
}

func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = SendMessage
	m.handlers[EventChangeChat] = ChatRoomHandler
}

func ChatRoomHandler(event Event, c *Client) error {
	var changeChatEvent ChangeChatEvent

	if err := json.Unmarshal(event.Payload, &changeChatEvent); err != nil {
		return fmt.Errorf("Bad Payload :%v", err)
	}

	c.chatID = changeChatEvent.ID

	return nil
}

func SendMessage(event Event, c *Client) error {
	var chatEvent SendMessageEvent

	if err := json.Unmarshal(event.Payload, &chatEvent); err != nil {
		return fmt.Errorf("Bad Payload :%v", err)
	}

	var broadMessage NewMessageEvent

	broadMessage.CreatedAt = time.Now()
	broadMessage.Content = chatEvent.Content
	broadMessage.UserID = chatEvent.UserID
	broadMessage.User = chatEvent.User
	broadMessage.ChatID = chatEvent.ChatID
	broadMessage.Read = false

	data, err := json.Marshal(broadMessage)

	if err != nil {
		return fmt.Errorf("Failed to Marshall BroadCast Message")
	}

	outgoingEvent := Event{
		Type:    EventNewMessage,
		Payload: data,
	}

	for client := range c.manager.clients {
		if client.chatID == broadMessage.ChatID {
			client.egress <- outgoingEvent
		}
	}

	return nil

}

func (m *Manager) routeEvent(event Event, c *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("No such Event Type present.")
	}
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {

	log.Println("Connected to WS")

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	params := r.URL.Query()
	chatID := params.Get("chatID")

	log.Print(chatID)

	client := NewClient(conn, m, chatID)

	m.addClient(client)

	go client.readMessages()
	go client.writeMessages()

	// conn.Close()
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
	}
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	switch origin {
	case "http://localhost:3000":
		return true
	default:
		return false
	}
}
