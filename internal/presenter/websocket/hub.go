// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file at https://github.com/gorilla/websocket.

package websocket

import (
	"github.com/pagient/pagient-server/internal/model"
	"github.com/pagient/pagient-server/internal/presenter/renderer"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	transmit chan *Message

	// stop running hub
	stop chan struct{}

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
		clients:    make(map[*Client]bool),
		transmit:   make(chan *Message),
	}
}

// Run performs client registration and message broadcasts in a goroutine
func (h *Hub) Run(stop <-chan struct{}) {
	go func() {
		for {
			select {
			case client := <-h.Register:
				h.clients[client] = true
			case client := <-h.Unregister:
				if _, ok := h.clients[client]; ok {
					delete(h.clients, client)
					close(client.send)
				}
			case message := <-h.transmit:
				for client := range h.clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			case <-stop:
				return
			}
		}
	}()
}

// broadcast creates a message and adds it to the broadcast channel
func (h *Hub) broadcast(msgType MessageType, data interface{}) {
	msg := &Message{
		Type: msgType,
		Data: data,
	}

	h.transmit <- msg
}

// NotifyNewPatient broadcasts a notification about a new patient
func (h *Hub) NotifyNewPatient(patient *model.Patient) {
	h.broadcast(MessageTypePatientAdd, renderer.NewPatientResponse(patient))
}

// NotifyUpdatedPatient broadcasts a notification about an updated patient
func (h *Hub) NotifyUpdatedPatient(patient *model.Patient) {
	h.broadcast(MessageTypePatientUpdate, renderer.NewPatientResponse(patient))
}

// NotifyDeletedPatient broadcasts a notification about a deleted patient
func (h *Hub) NotifyDeletedPatient(patient *model.Patient) {
	h.broadcast(MessageTypePatientDelete, renderer.NewPatientResponse(patient))
}

// DisconnectClient disconnects a client by token signature
func (h *Hub) DisconnectClient(id uint) {
	for client := range h.clients {
		if client.id == id {
			h.Unregister <- client
		}
	}
}
