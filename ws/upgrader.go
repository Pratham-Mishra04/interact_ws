package ws

import (
	"net/http"

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
