package realtime

import "github.com/google/uuid"

type Message struct {
	ChannelID uuid.UUID `json:"channelId"`
	Payload   string    `json:"payload"`
}
