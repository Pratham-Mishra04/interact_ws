package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Pratham-Mishra04/interactWS/initializers"
	"github.com/Pratham-Mishra04/interactWS/ws"
)

func init() {
	initializers.LoadEnv()
}

func main() {
	setupAPI()

	fmt.Println("Running the Server on PORT " + initializers.CONFIG.PORT)
	log.Fatal(http.ListenAndServe(":"+initializers.CONFIG.PORT, nil))
}

func setupAPI() {
	manager := ws.NewManager()
	http.HandleFunc("/ws", manager.ServeWS)
}
