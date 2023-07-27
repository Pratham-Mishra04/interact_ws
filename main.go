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
	initializers.AddLogger()
}

func main() {
	defer initializers.LoggerCleanUp()
	setupAPI()

	fmt.Println("Running the Server on PORT " + initializers.CONFIG.PORT)
	log.Fatal(http.ListenAndServe(":"+initializers.CONFIG.PORT, nil))
}

func setupAPI() {
	manager := ws.NewManager()
	http.HandleFunc("/ws", manager.ServeWS)
}
