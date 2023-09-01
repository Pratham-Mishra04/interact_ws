package ws

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/Pratham-Mishra04/interactWS/helpers"
	"github.com/Pratham-Mishra04/interactWS/initializers"
)

type Manager struct {
	clients ClientList
	sync.RWMutex

	handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}

	m.setupEventHandlers()
	return m
}

func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = SendMessageHandler
	// m.handlers[EventChangeChat] = ChatRoomHandler
	m.handlers[MeTyping] = MeTypingHandler
	m.handlers[MeStopTyping] = MeStopTypingHandler
	m.handlers[ChatSetup] = ChatSetupHandler
	m.handlers[EventSendNotification] = NotificationHandler
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("invalid event requested")
	}
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		helpers.LogError("Error upgrading the connection. ", err)
		return
	}

	params := r.URL.Query()
	userID := params.Get("userID")

	client := NewClient(conn, m, userID)

	m.addClient(client)

	if initializers.CONFIG.ENV == initializers.DevelopmentEnv {
		fmt.Println("New Connect established for user: " + userID)
	}

	go client.readMessages()
	go client.writeMessages()
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
