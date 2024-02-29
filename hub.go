package main

type Hub struct {
	broadcast chan []byte
	clients   map[*Client]bool
}

func newHub() *Hub {
	return &Hub{
		broadcast: make(chan []byte),
		clients:   make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for message := range h.broadcast {
		for client := range h.clients {
			client.send <- message
		}
	}
}
