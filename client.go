package main

import "github.com/gorilla/websocket"

type client struct {
	// Storage for websocket when HTTP connection has been upgraded
	socket *websocket.Conn
	//  Channel for sending messages
	send chan []byte
	// The room the client is chatting in
	room *room
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
