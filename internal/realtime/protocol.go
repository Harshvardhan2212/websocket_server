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
	Name        RoleName                `json:"name"`
	Permissions map[Permission]struct{} `json:"permissions"`
}

type RoleName string

var (
	Admin     RoleName = "admin"
	Modarator RoleName = "modarator"
	Member    RoleName = "member"
	Guest     RoleName = "guest"
)

var Roles = map[RoleName]*Role{
	Admin: {
		Name: Admin,
		Permissions: map[Permission]struct{}{
			CanSend:          {},
			CanCreateChannel: {},
			CanDeleteChannel: {},
			CanInviteUser:    {},
			CanKick:          {},
			CanMute:          {},
		},
	},

	Modarator: {
		Name: Modarator,
		Permissions: map[Permission]struct{}{
			CanSend:       {},
			CanInviteUser: {},
			CanKick:       {},
			CanMute:       {},
		},
	},
	Member: {
		Name: Member,
		Permissions: map[Permission]struct{}{
			CanSend: {},
		},
	},
	Guest: {
		Name:        Guest,
		Permissions: map[Permission]struct{}{},
	},
}
