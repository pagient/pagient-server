// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file at https://github.com/gorilla/websocket.

package websocket

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
}

// NewHub creates and returns a new websocket hub
func NewHub() *Hub {
	return &Hub{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		clients:    make(map[*Client]bool),
	}
}

// Run performs client registration and message broadcasts in a goroutine
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// Broadcast creates a message and adds it to the broadcast channel
func (h *Hub) Broadcast(msgType MessageType, data interface{}) error {
	patientMsg := &Message{
		Type: msgType,
		Data: data,
	}

	msg, err := json.Marshal(patientMsg)
	if err != nil {
		return errors.Wrap(err, "json marshal failed")
	}

	h.broadcast <- msg

	return nil
}
