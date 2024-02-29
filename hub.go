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

// run loops trough the broadcast chan and sends the message to all the hub subs
func (h *Hub) run() {
	for message := range h.broadcast {
		for client := range h.clients {
			client.send <- message
		}
	}
}
