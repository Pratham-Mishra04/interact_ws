package ws

import (
	"net/http"

	"github.com/Pratham-Mishra04/interactWS/config"
	"github.com/Pratham-Mishra04/interactWS/initializers"
	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  config.READ_BUFFER_SIZE,
		WriteBufferSize: config.WRITE_BUFFER_SIZE,
		CheckOrigin:     checkOrigin,
	}
)

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	switch origin {
	case initializers.CONFIG.FRONTEND_URL:
		return true
	case initializers.CONFIG.DEV_URL:
		return initializers.CONFIG.ENV == "development"
	default:
		return false
	}
}
