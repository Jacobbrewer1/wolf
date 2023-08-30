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

const guildDalName = "guild_dal"

type IGuildDal interface {
	// SaveGuild saves a guild.
	SaveGuild(guild *entities.Guild) error

	// GetGuildByID gets a guild by ID.
	GetGuildByID(id string) (*entities.Guild, error)
}

type guildDal struct {
	// ctx is the context.
	ctx context.Context

	// l is the logger.
	l *slog.Logger

	// client is the database.
	client *mongo.Client
}

// NewGuildDal creates a new guild data access layer.
func NewGuildDal(ctx context.Context, logger *slog.Logger) IGuildDal {
	// If the context is nil, create a new one.
	if ctx == nil {
		ctx = context.Background()
	}

	l := logger.With(slog.String(logging.KeyDal, guildDalName))

	if MongoDB == nil {
		l.Warn("MongoDB is nil, this can cause a panic. Proceeding...")
	}

	return &guildDal{
		ctx:    ctx,
		l:      l,
		client: MongoDB,
	}
}

func (g *guildDal) SaveGuild(guild *entities.Guild) error {
	// Get the guild collection.
	collection := g.client.Database(mongoDatabase).Collection("guilds")

	// Start the prometheus metrics.
	monitoring.MongoTotalRequests.WithLabelValues(guildDalName, "save_guild_config", mongoDatabase, "guilds").Inc()
	t := prometheus.NewTimer(monitoring.MongoLatency.WithLabelValues(guildDalName, "save_guild_config", mongoDatabase, "guilds"))
	defer t.ObserveDuration()

	// Save the guild.
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(g.ctx, bson.M{"id": guild.ID}, bson.M{"$set": guild}, opts)
	if err != nil {
		return fmt.Errorf("error updating guild: %w", err)
	}
	return nil
}

// GetGuildByID gets a guild by ID.
func (g *guildDal) GetGuildByID(id string) (*entities.Guild, error) {
	// Get the guild collection.
	collection := g.client.Database(mongoDatabase).Collection("guilds")

	// Start the prometheus metrics.
	monitoring.MongoTotalRequests.WithLabelValues(guildDalName, "get_guild_by_id", mongoDatabase, "guilds").Inc()
	t := prometheus.NewTimer(monitoring.MongoLatency.WithLabelValues(guildDalName, "get_guild_by_id", mongoDatabase, "guilds"))
	defer t.ObserveDuration()

	// Get the guild.
	guild := new(entities.Guild)

	err := collection.FindOne(g.ctx, bson.M{"id": id}).Decode(guild)
	if err != nil {
		return nil, fmt.Errorf("error getting guild: %w", err)
	}
	return guild, nil
}
