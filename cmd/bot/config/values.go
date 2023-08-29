package config

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

	// EmailFrom is the email address to send emails from.
	EmailFrom string

	// EmailPassword is the password for the email address.
	EmailPassword string

	// EmailHost is the host for the email address.
	EmailHost string

	// EmailPort is the port for the email address.
	EmailPort string

	// MonitoringPort is the port for the monitoring server.
	MonitoringPort string
)
