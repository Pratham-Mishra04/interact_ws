package ws

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/Pratham-Mishra04/interactWS/initializers"
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
	m.handlers[EventSendMessage] = SendMessageHandler
	m.handlers[EventChangeChat] = ChatRoomHandler
	m.handlers[MeTyping] = MeTypingHandler
	m.handlers[MeStopTyping] = MeStopTypingHandler
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
	case initializers.CONFIG.FRONTEND_URL:
		return true
	case initializers.CONFIG.DEV_URL:
		return initializers.CONFIG.ENV == "dev"
	default:
		return false
	}
}
