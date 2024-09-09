package client

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	send     chan []byte
	receive  chan []byte
	done     chan struct{}
	OnMessage func([]byte)
}

func NewClient(url string) (*Client, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	client := &Client{
		conn:    conn,
		send:    make(chan []byte),
		receive: make(chan []byte),
		done:    make(chan struct{}),
	}

	go client.readPump()
	go client.writePump()

	return client, nil
}

func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
		close(c.done)
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			return
		}
		if c.OnMessage != nil {
			c.OnMessage(message)
		}
		c.receive <- message
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		select {
		case message := <-c.send:
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("write error:", err)
				return
			}
		case <-c.done:
			return
		}
	}
}

func (c *Client) Send(message []byte) {
	c.send <- message
}

func (c *Client) Receive() <-chan []byte {
	return c.receive
}

func (c *Client) Close() {
	close(c.done)
	c.conn.Close()
}
