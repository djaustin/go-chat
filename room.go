package main

import (
	"log"
	"net/http"
	"os"

	"github.com/djaustin/go-chat/trace"
	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

type room struct {
	// forward is a channel for sending incoming messages to sent to all clients
	forward chan []byte
	// clients contains all clients currently connected to the room
	clients map[*client]bool
	// join is a channel for clients attempting to join the room
	join chan *client
	// leave is a channel for clients attempting to leave the room
	leave chan *client
	// tracer receives trace information of activity in the room
	tracer trace.Tracer
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// a client is waiting to join the room
			r.clients[client] = true
			r.tracer.Trace("Client entering room")
		case client := <-r.leave:
			// a client is waiting to leave the room
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client leaving room")
		case msg := <-r.forward:
			// a message has been received from a client. Distribute it
			for client := range r.clients {
				select {
				case client.send <- msg:
					r.tracer.Trace("Sent message ", string(msg), " to client")
					// we did it!
				default:
					// send failed
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace("Send to client failed. Cleaning up.")
				}
			}
		}
	}
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}
	client := &client{
		socket: socket,
		room:   r,
		send:   make(chan []byte, messageBufferSize),
	}

	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.New(os.Stdout),
	}
}
