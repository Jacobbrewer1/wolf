package dataaccess

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Jacobbrewer1/wolf/pkg/dataaccess/monitoring"
	"github.com/Jacobbrewer1/wolf/pkg/entities"
	"github.com/Jacobbrewer1/wolf/pkg/logging"
	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const ticketDalName = "ticket_dal"

var TicketDB TicketDal

type TicketDal interface {
	// SaveTicket saves a ticket.
	SaveTicket(ctx context.Context, ticket *entities.Ticket) error

	// GetTicket gets a ticket by name.
	GetTicket(ctx context.Context, guildID string, channelID string) (*entities.Ticket, error)

	// GetLatestTicket gets the latest ticket.
	GetLatestTicket(ctx context.Context, guildID string) (*entities.Ticket, error)
}

type ticketDalImpl struct {
	// l is the logger.
	l *slog.Logger

	// client is the database.
	client *mongo.Client
}

// NewTicketDal creates a new ticket data access layer.
func NewTicketDal() TicketDal {
	l := slog.Default().With(slog.String(logging.KeyDal, guildDalName))

	if MongoDB == nil {
		l.Warn("MongoDB is nil, this can cause a panic. Proceeding...")
	}

	return &ticketDalImpl{
		l:      l,
		client: MongoDB,
	}
}

func (d *ticketDalImpl) SaveTicket(ctx context.Context, ticket *entities.Ticket) error {
	// Get the guild collection.
	collection := d.client.Database(mongoDatabase).Collection("tickets")

	// Start the prometheus metrics.
	monitoring.MongoTotalRequests.WithLabelValues(ticketDalName, "save_ticket", mongoDatabase, "tickets").Inc()
	t := prometheus.NewTimer(monitoring.MongoLatency.WithLabelValues(ticketDalName, "save_ticket", mongoDatabase, "tickets"))
	defer t.ObserveDuration()

	// Save the guild.
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, bson.M{"guild_id": ticket.GuildID, "channel_id": ticket.ChannelID}, bson.M{"$set": ticket}, opts)
	if err != nil {
		return fmt.Errorf("error updating guild: %w", err)
	}
	return nil
}

func (d *ticketDalImpl) GetTicket(ctx context.Context, guildID string, channelID string) (*entities.Ticket, error) {
	// Get the guild collection.
	collection := d.client.Database(mongoDatabase).Collection("tickets")

	// Start the prometheus metrics.
	monitoring.MongoTotalRequests.WithLabelValues(ticketDalName, "get_ticket", mongoDatabase, "tickets").Inc()
	t := prometheus.NewTimer(monitoring.MongoLatency.WithLabelValues(ticketDalName, "get_ticket", mongoDatabase, "tickets"))
	defer t.ObserveDuration()

	// Get the ticket.
	var ticket entities.Ticket
	err := collection.FindOne(ctx, bson.M{
		"guild_id":   guildID,
		"channel_id": channelID,
		"deleted":    false,
	}).Decode(&ticket)
	if err != nil {
		return nil, fmt.Errorf("error getting ticket: %w", err)
	}

	return &ticket, nil
}

func (d *ticketDalImpl) GetLatestTicket(ctx context.Context, guildID string) (*entities.Ticket, error) {
	// Get the guild collection.
	collection := d.client.Database(mongoDatabase).Collection("tickets")

	// Start the prometheus metrics.
	monitoring.MongoTotalRequests.WithLabelValues(ticketDalName, "get_latest_ticket", mongoDatabase, "tickets").Inc()
	t := prometheus.NewTimer(monitoring.MongoLatency.WithLabelValues(ticketDalName, "get_latest_ticket", mongoDatabase, "tickets"))
	defer t.ObserveDuration()

	// Set the options to get the latest ticket.
	opts := options.FindOne()
	opts.SetSort(bson.M{"created_at": -1})

	// Get the ticket.
	var ticket entities.Ticket
	err := collection.FindOne(ctx, bson.M{"guild_id": guildID}, opts).Decode(&ticket)
	if err != nil {
		return nil, fmt.Errorf("error getting ticket: %w", err)
	}

	return &ticket, nil
}
