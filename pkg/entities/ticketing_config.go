package entities

type TicketingConfig struct {
	// Enabled is whether ticketing is enabled.
	Enabled bool `json:"enabled" bson:"enabled"`

	// ChannelID is the ID of the channel that ticketing is enabled in.
	ChannelID string `json:"channel_id" bson:"channel_id"`

	// RoleID is the ID of the role that handles tickets.
	RoleID string `json:"role_id" bson:"role_id"`

	// OpenMessageID is the ID of the open ticket message.
	OpenMessageID string `json:"open_message_id" bson:"open_message_id"`

	// CreatedTicketsCategoryID is the ID of the category that created tickets are put in.
	CreatedTicketsCategoryID string `json:"created_tickets_category_id" bson:"created_tickets_category_id"`

	// ClaimedTicketsCategoryID is the ID of the category that claimed tickets are put in.
	ClaimedTicketsCategoryID string `json:"claimed_tickets_category_id" bson:"claimed_tickets_category_id"`

	// ClosedTicketsCategoryID is the ID of the category that closed tickets are put in.
	ClosedTicketsCategoryID string `json:"closed_tickets_category_id" bson:"closed_tickets_category_id"`
}
