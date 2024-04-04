package socket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"strings"
)

type Hub struct {
	clients map[*Client]bool

	broadcast chan []byte

	register chan *Client

	unregister chan *Client
}

var h *Hub

func init() {
	h = &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func GetHub() *Hub {
	return h
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
		case message := <-h.broadcast:
			msg := strings.Split(string(message), " ")
			for client := range h.clients {
				if client.Channel == msg[0] {
					err := client.Conn.WriteMessage(websocket.TextMessage, []byte(msg[1]))
					if err != nil {
						fmt.Println("websocket发送消息失败")
						return
					}
				}
			}
		}
	}
}
