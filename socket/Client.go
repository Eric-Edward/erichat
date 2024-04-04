package socket

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type Client struct {
	Hub     *Hub
	Uid     string
	Cid     string
	Channel string
	Conn    *websocket.Conn
	Send    chan []byte
	Quit    chan bool
}

func NewWSClient(Uid, Cid, Channel string, Conn *websocket.Conn) *Client {
	c := &Client{
		Hub:     GetHub(),
		Uid:     Uid,
		Cid:     Cid,
		Channel: Channel,
		Conn:    Conn,
		Send:    make(chan []byte),
		Quit:    make(chan bool),
	}
	c.Hub.register <- c
	go c.RMessage()
	go c.SMessage()
	return c
}

func (c *Client) RMessage() {
	defer func() {
		c.Hub.unregister <- c
		close(c.Quit)
		_ = c.Conn.Close()
	}()
	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println("读取数据失败")
			return
		}

		select {
		case c.Send <- p:
		case <-c.Quit:
			return
		}
	}
}

func (c *Client) SMessage() {
	defer func() {
		c.Hub.unregister <- c
		close(c.Quit)
		_ = c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			{
				if !ok {
					fmt.Println("从数据库读取信息失败读取信息失败")
					return
				}
				c.Hub.broadcast <- []byte(c.Channel + " " + string(message))
			}
		case <-c.Quit:
			fmt.Println("用户主动断开websocket连接")
			return
		}
	}
}
