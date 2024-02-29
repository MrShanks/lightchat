package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	name string
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) reader() {
	defer func() {
		c.conn.Close()
		c.hub.clients[c] = false
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.hub.broadcast <- append([]byte(fmt.Sprintf("%s: ", c.name)), message...)
	}
}

func (c *Client) writer() {
	defer c.conn.Close()

	for message := range c.send {
		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			log.Println(err)
			return
		}
		w.Write(message)

		if err := w.Close(); err != nil {
			log.Fatal(err)
		}
	}
}
