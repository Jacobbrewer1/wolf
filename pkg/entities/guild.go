package entities

// Guild is a configuration for a guild.
type Guild struct {
	// ID is the ID of the guild.
	ID string `json:"id" bson:"id"`

	// Ticketing is the ticketing configuration.
	Ticketing TicketingConfig `json:"ticketing" bson:"ticketing"`
}
