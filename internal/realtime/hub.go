// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package realtime

import (
	"log"

	"github.com/google/uuid"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	// clients map[*Client]bool  // NOTE: work only for a single braodcast
	channels map[uuid.UUID]map[*Client]bool // NOTE: for listing channels and subscrib a users

	// Inbound messages from the clients.
	broadcast chan *Message

	// Register requests from the clients.
	register chan *Register

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Message),
		register:   make(chan *Register),
		unregister: make(chan *Client),
		// clients:    make(map[*Client]bool),
		channels: make(map[uuid.UUID]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case register := <-h.register:
			channID := register.ChannelID
			if _, ok := h.channels[channID]; !ok {
				h.channels[channID] = make(map[*Client]bool)
			}
			h.channels[channID][register.Client] = true

			log.Println("register is called for", h.channels)
			log.Printf("register is called for %+v", register.Client)
			// h.clients[client] = true
		case client := <-h.unregister:
			for _, clients := range h.channels {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
				}
			}
			log.Println("unregister is called for", client)
		case message := <-h.broadcast:
			clients := h.channels[message.ChannelID]
			log.Printf("message: %+v\n", message)
			for client := range clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(clients, client)
					log.Printf("deleted a client %+v\n", client)
				}
			}
		}
	}
}
