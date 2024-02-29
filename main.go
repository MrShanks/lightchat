package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var guestNumber int

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

// wsHandler upgrades the http to a webSocket conn
// it accept messages and returns them back
func wsHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	guestNumber++
	name := fmt.Sprintf("Guest %d", guestNumber)
	client := &Client{name: name, hub: hub, conn: conn, send: make(chan []byte, 256)}
	hub.clients[client] = true

	go client.writer()
	go client.reader()
}

func main() {
	hub := newHub()
	go hub.run()
	server := &http.Server{
		Addr:              ":9090",
		ReadHeaderTimeout: 3 * time.Second,
	}

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { wsHandler(hub, w, r) })
	log.Fatal(server.ListenAndServe())
}
