package dataaccess

import (
	"context"
	"fmt"

	"github.com/Jacobbrewer1/wolf/pkg/dataaccess/monitoring"
	"github.com/Jacobbrewer1/wolf/pkg/entities"
	"github.com/Jacobbrewer1/wolf/pkg/logging"
	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slog"
)

const ticketDalName = "ticket_dal"

type ITicketDal interface {
	// SaveTicket saves a ticket.
	SaveTicket(ticket *entities.Ticket) error

	// GetTicket gets a ticket by name.
	GetTicket(guildID string, channelID string) (*entities.Ticket, error)

	// GetLatestTicket gets the latest ticket.
	GetLatestTicket(guildID string) (*entities.Ticket, error)
}

type ticketDal struct {
	// ctx is the context.
	ctx context.Context

	// l is the logger.
	l *slog.Logger

	// client is the database.
	client *mongo.Client
}

// NewTicketDal creates a new ticket data access layer.
func NewTicketDal(ctx context.Context, logger *slog.Logger) ITicketDal {
	// If the context is nil, create a new one.
	if ctx == nil {
		ctx = context.Background()
	}

	l := logger.With(slog.String(logging.KeyDal, guildDalName))

	if MongoDB == nil {
		l.Warn("MongoDB is nil, this can cause a panic. Proceeding...")
	}

	return &ticketDal{
		ctx:    ctx,
		l:      l,
		client: MongoDB,
	}
}

func (d *ticketDal) SaveTicket(ticket *entities.Ticket) error {
	// Get the guild collection.
	collection := d.client.Database(mongoDatabase).Collection("tickets")

	// Start the prometheus metrics.
	monitoring.MongoTotalRequests.WithLabelValues(ticketDalName, "save_ticket", mongoDatabase, "tickets").Inc()
	t := prometheus.NewTimer(monitoring.MongoLatency.WithLabelValues(ticketDalName, "save_ticket", mongoDatabase, "tickets"))
	defer t.ObserveDuration()

	// Save the guild.
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(d.ctx, bson.M{"guild_id": ticket.GuildID, "channel_id": ticket.ChannelID}, bson.M{"$set": ticket}, opts)
	if err != nil {
		return fmt.Errorf("error updating guild: %w", err)
	}
	return nil
}

func (d *ticketDal) GetTicket(guildID string, channelID string) (*entities.Ticket, error) {
	// Get the guild collection.
	collection := d.client.Database(mongoDatabase).Collection("tickets")

	// Start the prometheus metrics.
	monitoring.MongoTotalRequests.WithLabelValues(ticketDalName, "get_ticket", mongoDatabase, "tickets").Inc()
	t := prometheus.NewTimer(monitoring.MongoLatency.WithLabelValues(ticketDalName, "get_ticket", mongoDatabase, "tickets"))
	defer t.ObserveDuration()

	// Get the ticket.
	var ticket entities.Ticket
	err := collection.FindOne(d.ctx, bson.M{
		"guild_id":   guildID,
		"channel_id": channelID,
	}).Decode(&ticket)
	if err != nil {
		return nil, fmt.Errorf("error getting ticket: %w", err)
	}

	return &ticket, nil
}

func (d *ticketDal) GetLatestTicket(guildID string) (*entities.Ticket, error) {
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
	err := collection.FindOne(d.ctx, bson.M{"guild_id": guildID}, opts).Decode(&ticket)
	if err != nil {
		return nil, fmt.Errorf("error getting ticket: %w", err)
	}

	return &ticket, nil
}
