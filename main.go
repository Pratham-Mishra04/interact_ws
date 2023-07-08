package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Pratham-Mishra04/interactWS/ws"
)

func main() {
	setupAPI()

	fmt.Println("Running the Server on PORT 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupAPI() {

	manager := ws.NewManager()

	http.HandleFunc("/ws", manager.ServeWS)
}
