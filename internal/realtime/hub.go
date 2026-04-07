package realtime

import (
	"github.com/google/uuid"
)

type Hub struct {
	// NOTE: for listing channels and subscrib a users
	channels map[uuid.UUID]map[uuid.UUID]*Client

	// Inbound messages from the clients.
	broadcast chan *Message

	registerClient chan *Register

	// Unregister requests from clients.
	unregisterClient chan *Client

	registerChannel chan uuid.UUID

	unregisterChannel chan uuid.UUID
}

func NewHub() *Hub {
	return &Hub{
		broadcast:         make(chan *Message),
		registerClient:    make(chan *Register),
		unregisterClient:  make(chan *Client),
		channels:          make(map[uuid.UUID]map[uuid.UUID]*Client),
		registerChannel:   make(chan uuid.UUID),
		unregisterChannel: make(chan uuid.UUID),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case channel := <-h.registerChannel:
			h.channels[channel] = make(map[uuid.UUID]*Client)

		case channel := <-h.unregisterChannel:
			for _, client := range h.channels[channel] {
				close(client.send)
			}
			delete(h.channels, channel)

		case register := <-h.registerClient:
			channID := register.ChannelID
			if _, ok := h.channels[channID]; !ok {
				h.channels[channID] = make(map[uuid.UUID]*Client)
			}
			h.channels[channID][register.Client.ID] = register.Client

		case client := <-h.unregisterClient:
			for _, clients := range h.channels {
				if _, ok := clients[client.ID]; ok {
					delete(clients, client.ID)
					close(client.send)
				}
			}

		case message := <-h.broadcast:
			clients := h.channels[message.ChannelID]
			for _, client := range clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(clients, client.ID)
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

func (h *Hub) KickClient(channID, clientID uuid.UUID) bool {
	client, ok := h.channels[channID][clientID]
	h.unregisterClient <- client
	return ok
}

func (h *Hub) MuteClient(channID, clientID uuid.UUID) bool {
	client, ok := h.channels[channID][clientID]
	delete(client.Role.Permissions, CanSend)
	return ok
}
