package realtime

import "github.com/google/uuid"

type Message struct {
	ChannelID uuid.UUID `json:"channelId"`
	Payload   string    `json:"payload"`
}

type Permission int

const (
	CanSend Permission = iota + 1
	CanCreateChannel
	CanDeleteChannel
	CanInviteUser
	CanKick
	CanMute
)

type Role struct {
	Name        string                  `json:"name"`
	Permissions map[Permission]struct{} `json:"permissions"`
}

var Roles = map[string]*Role{
	"admin": {
		Name: "admin",
		Permissions: map[Permission]struct{}{
			CanSend:          {},
			CanCreateChannel: {},
			CanDeleteChannel: {},
			CanInviteUser:    {},
			CanKick:          {},
			CanMute:          {},
		},
	},

	"modarator": {
		Name: "modarator",
		Permissions: map[Permission]struct{}{
			CanSend:       {},
			CanInviteUser: {},
			CanKick:       {},
			CanMute:       {},
		},
	},
	"member": {
		Name: "member",
		Permissions: map[Permission]struct{}{
			CanSend: {},
		},
	},
	"guest": {
		Name:        "guest",
		Permissions: map[Permission]struct{}{},
	},
}
