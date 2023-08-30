package entities

import (
	"fmt"

	"github.com/Jacobbrewer1/wolf/pkg/custom"
)

// Ticket is a ticket.
type Ticket struct {
	// ID is the number of the ticket.
	// This is used to identify tickets with the users name.
	// For example, if the users username is "wolf", and the ticket ID is 1, the ticket name will be "1-wolf".
	ID int `json:"id" bson:"id"`

	// GuildID is the ID of the guild that the ticket is in.
	GuildID string `json:"guild_id" bson:"guild_id"`

	// ChannelID is the ID of the channel that the ticket is in.
	ChannelID string `json:"channel_id" bson:"channel_id"`

	// UserID is the ID of the user that created the ticket.
	UserID string `json:"user_id" bson:"user_id"`

	// Username is the username of the user that created the ticket.
	Username string `json:"username" bson:"username"`

	// SetupMessageID is the ID of the setup message.
	SetupMessageID string `json:"setup_message_id" bson:"setup_message_id"`

	// Claimed by is the ID of the user that claimed the ticket.
	ClaimedBy string `json:"claimed_by" bson:"claimed_by"`

	// ClosedBy is the ID of the user that closed the ticket.
	ClosedBy string `json:"closed_by" bson:"closed_by"`

	// Deleted is whether the ticket has been deleted.
	Deleted bool `json:"deleted" bson:"deleted"`

	// CreatedAt is the time that the ticket was created.
	CreatedAt custom.Datetime `json:"created_at" bson:"created_at"`
}

func (t *Ticket) Name() string {
	return fmt.Sprintf("%d-%s", t.ID, t.Username)
}
