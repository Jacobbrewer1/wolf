package config

import (
	"flag"
	"os"

	"github.com/Jacobbrewer1/wolf/pkg/dataaccess"
	"github.com/Jacobbrewer1/wolf/pkg/dataaccess/connection"
	"github.com/Jacobbrewer1/wolf/pkg/logging"
	"golang.org/x/exp/slog"
)

func Parse(l *slog.Logger) {
	flagBT := flag.String(EnvBotToken, "", "(Optional) The token for the bot.")
	flagAppId := flag.String(EnvApplicationId, "", "(Optional) The ID of the application.")
	flagMongoUri := flag.String(EnvMongoUri, "", "(Optional) The URI for the MongoDB database.")

	if envBT := os.Getenv(EnvBotToken); envBT != "" {
		l.Debug("Found bot token in environment", slog.String("key", EnvBotToken))
		BotToken = envBT
	}

	if envAppId := os.Getenv(EnvApplicationId); envAppId != "" {
		l.Debug("Found application ID in environment", slog.String("key", EnvApplicationId))
		ApplicationId = envAppId
	}

	if envMongoUri := os.Getenv(EnvMongoUri); envMongoUri != "" {
		l.Debug("Found MongoDB URI in environment", slog.String("key", EnvMongoUri))
		MongoUri = envMongoUri
	}

	if envMonitoringPort := os.Getenv(EnvMonitoringPort); envMonitoringPort != "" {
		l.Debug("Found monitoring port in environment", slog.String("key", EnvMonitoringPort))
		MonitoringPort = envMonitoringPort
	} else {
		// Default to 8080 if not provided.
		MonitoringPort = "8080"

		l.Info("No monitoring port provided in environment, defaulting to 8080", slog.String("key", EnvMonitoringPort))
	}

	if BotToken != "" &&
		ApplicationId != "" &&
		MongoUri != "" &&
		EmailFrom != "" &&
		EmailPassword != "" &&
		EmailHost != "" &&
		EmailPort != "" {

		// All required environment variables have been provided.
		l.Debug("All required environment variables have been provided")
		connectMongo(l)
		return
	}

	flag.Parse()

	if *flagBT != "" {
		BotToken = *flagBT
	} else {
		l.Error("No bot token provided in environment or flag", slog.String("key", EnvBotToken))
		os.Exit(1)
	}

	if *flagAppId != "" {
		ApplicationId = *flagAppId
	} else {
		l.Error("No application ID provided in environment or flag", slog.String("key", EnvApplicationId))
		os.Exit(1)
	}

	if *flagMongoUri != "" {
		MongoUri = *flagMongoUri
		connectMongo(l)
	} else {
		l.Error("No MongoDB URI provided in environment or flag", slog.String("key", EnvMongoUri))
		os.Exit(1)
	}
}

func connectMongo(l *slog.Logger) {
	mongoConn := new(connection.MongoDB)
	mongoConn.ConnectionString = MongoUri

	db, err := mongoConn.Connect()
	if err != nil {
		l.Error("Error connecting to mongo", slog.String(logging.KeyError, err.Error()))
		os.Exit(1)
	} else if db == nil {
		l.Error("MongoDB came back nil", slog.String(logging.KeyError, "MongoDB came back nil"))
		os.Exit(1)
	}

	dataaccess.MongoDB = db

	l.Debug("Connected to MongoDB", slog.String("key", EnvMongoUri))
}
