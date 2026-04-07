// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package realtime

import (
	"github.com/google/uuid"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// NOTE: for listing channels and subscrib a users
	channels map[uuid.UUID]map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan *Message

	// Register requests from the clients.
	register chan *Register

	// Unregister requests from clients.
	unregister chan *Client

	registerChannel chan uuid.UUID

	unregisterChannel chan uuid.UUID
}

func NewHub() *Hub {
	return &Hub{
		broadcast:         make(chan *Message),
		register:          make(chan *Register),
		unregister:        make(chan *Client),
		channels:          make(map[uuid.UUID]map[*Client]bool),
		registerChannel:   make(chan uuid.UUID),
		unregisterChannel: make(chan uuid.UUID),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case channel := <-h.registerChannel:
			h.channels[channel] = make(map[*Client]bool)

		case channel := <-h.unregisterChannel:
			for client := range h.channels[channel] {
				close(client.send)
			}
			delete(h.channels, channel)

		case register := <-h.register:
			channID := register.ChannelID
			if _, ok := h.channels[channID]; !ok {
				h.channels[channID] = make(map[*Client]bool)
			}
			h.channels[channID][register.Client] = true

		case client := <-h.unregister:
			for _, clients := range h.channels {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
				}
			}

		case message := <-h.broadcast:
			clients := h.channels[message.ChannelID]
			for client := range clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(clients, client)
				}
			}
		}
	}
}

func (h *Hub) CreateChannel() uuid.UUID {
	channID := uuid.New()
	h.registerChannel <- channID
	return channID
}

func (h *Hub) DeleteChannel(channID uuid.UUID) bool {
	h.unregisterChannel <- channID
	return true
}
