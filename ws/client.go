package ws

import (
	"encoding/json"
	"time"

	"github.com/Pratham-Mishra04/interactWS/config"
	"github.com/Pratham-Mishra04/interactWS/helpers"
	"github.com/gorilla/websocket"
)

var (
	pongWait     = config.PONG_WAIT
	pingInterval = config.PING_INTERVAL
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	chatID     string
	egress     chan Event
}

func NewClient(conn *websocket.Conn, manager *Manager, chatID string) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		chatID:     chatID,
		egress:     make(chan Event),
	}
}

func (c *Client) readMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		helpers.LogError("Error setting connection deadline", err)
		return
	}

	c.connection.SetReadLimit(1024)

	c.connection.SetPongHandler(c.pongHandler)

	for {
		_, payload, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				helpers.LogError("Error reading a message, abnormal connection closure", err)
			}
			break
		}

		var request Event

		if err := json.Unmarshal(payload, &request); err != nil {
			helpers.LogError("Error UnMarshalling an event", err)
			break
		}

		if err := c.manager.routeEvent(request, c); err != nil {
			helpers.LogError("Error Processing an event", err)
		}
	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	ticker := time.NewTicker(pingInterval)

	for {
		select {
		case message, ok := <-c.egress:
			if !ok { // Some error with the channel
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					helpers.LogError("Error writing a message, abnormal connection closure", err)
				}
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				helpers.LogError("Error marshalling the data", err)
				return
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				helpers.LogError("Error sending messages", err)
			}

		case <-ticker.C:
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
				helpers.LogWarn("Ticker Write Message Error", err)
				return
			}
		}
	}
}

func (c *Client) pongHandler(pongMsg string) error {
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}
