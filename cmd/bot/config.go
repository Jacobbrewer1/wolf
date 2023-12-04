package main

import (
	"log/slog"
	"os"

	"github.com/Jacobbrewer1/wolf/pkg/dataaccess"
	"github.com/Jacobbrewer1/wolf/pkg/dataaccess/connection"
	"github.com/Jacobbrewer1/wolf/pkg/logging"
)

const (
	// AppName is the name of the application.
	AppName = "wolf"

	// EnvBotToken is the environment variable for the bot token.
	EnvBotToken = `BOT_TOKEN`

	// EnvApplicationId is the environment variable for the application ID.
	EnvApplicationId = `APPLICATION_ID`

	// EnvMongoUri is the environment variable for the MongoDB URI.
	EnvMongoUri = `MONGO_URI`

	// EnvMonitoringPort is the environment variable for the monitoring port.
	EnvMonitoringPort = `MONITORING_PORT`
)

var (
	// BotToken is the token for the bot.
	BotToken string

	// ApplicationId is the ID of the application.
	ApplicationId string

	// MongoUri is the URI for the MongoDB database.
	MongoUri string

	// MonitoringPort is the port for the monitoring server.
	MonitoringPort string
)

func parseConfig() {
	if envBT := os.Getenv(EnvBotToken); envBT != "" {
		slog.Debug("Found bot token in environment", slog.String("key", EnvBotToken))
		BotToken = envBT
	}

	if envAppId := os.Getenv(EnvApplicationId); envAppId != "" {
		slog.Debug("Found application ID in environment", slog.String("key", EnvApplicationId))
		ApplicationId = envAppId
	}

	if envMongoUri := os.Getenv(EnvMongoUri); envMongoUri != "" {
		slog.Debug("Found MongoDB URI in environment", slog.String("key", EnvMongoUri))
		MongoUri = envMongoUri
	}

	if envMonitoringPort := os.Getenv(EnvMonitoringPort); envMonitoringPort != "" {
		slog.Debug("Found monitoring port in environment", slog.String("key", EnvMonitoringPort))
		MonitoringPort = envMonitoringPort
	} else {
		// Default to 8080 if not provided.
		MonitoringPort = "8080"
		slog.Info("No monitoring port provided in environment, defaulting to 8080", slog.String("key", EnvMonitoringPort))
	}

	if BotToken != "" &&
		ApplicationId != "" &&
		MongoUri != "" {

		// All required environment variables have been provided.
		slog.Debug("All required environment variables have been provided")
		connectMongo()
		return
	}

	slog.Error("Not all required environment variables have been provided", slog.String(logging.KeyError, "Incomplete configuration"))
	os.Exit(1)
}

func connectMongo() {
	mongoConn := new(connection.MongoDB)
	mongoConn.ConnectionString = MongoUri

	db, err := mongoConn.Connect()
	if err != nil {
		slog.Error("Error connecting to mongo", slog.String(logging.KeyError, err.Error()))
		os.Exit(1)
	} else if db == nil {
		slog.Error("MongoDB came back nil", slog.String(logging.KeyError, "MongoDB came back nil"))
		os.Exit(1)
	}

	dataaccess.MongoDB = db
	slog.Debug("Connected to MongoDB", slog.String("key", EnvMongoUri))
}
