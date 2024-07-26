package ws

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Pratham-Mishra04/interactWS/helpers"
	"github.com/Pratham-Mishra04/interactWS/initializers"
	"github.com/golang-jwt/jwt/v5"
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
	m.handlers[GetRead] = ReadMessageHandler
	m.handlers[SendUpdateMembershipEvent] = UpdateMembershipHandler
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
		helpers.LogError("Error upgrading the connection. ", err, "ServeWS")
		return
	}

	params := r.URL.Query()
	userID := params.Get("userID")
	token := params.Get("token")

	if err := verifyToken(token, userID); err != nil {
		helpers.LogWarn("Token verification failed: ", err, "ServeWS")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized: Token verification failed"))
	} else {
		client := NewClient(conn, m, userID)

		m.addClient(client)

		if initializers.CONFIG.ENV == initializers.DevelopmentEnv {
			fmt.Println("New Connect established for user: " + userID)
		}

		go client.readMessages()
		go client.writeMessages()
	}
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

func verifyToken(tokenString string, userID string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(initializers.CONFIG.JWT_SECRET), nil
	})

	if err != nil {
		return fmt.Errorf("token has expired")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return fmt.Errorf("token has expired")
		}

		tokenUserID, ok := claims["sub"].(string)
		if !ok {
			return fmt.Errorf("invalid user ID in token claims")
		}

		if userID != tokenUserID {
			return fmt.Errorf("invalid token for this user")
		}

		return nil
	} else {
		return fmt.Errorf("invalid Token")
	}
}
