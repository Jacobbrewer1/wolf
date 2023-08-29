package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/Jacobbrewer1/discordgo"
	"github.com/Jacobbrewer1/wolf/cmd/bot/config"
	"github.com/Jacobbrewer1/wolf/cmd/bot/monitoring"
	"github.com/Jacobbrewer1/wolf/pkg/logging"
	"github.com/Jacobbrewer1/wolf/pkg/request"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/exp/slog"
)

// IApp is the interface for the application.
type IApp interface {
	// Log returns the logger.
	Log() *slog.Logger

	// Session returns the discord session.
	Session() *discordgo.Session
}

type App struct {
	// is the logger.
	*slog.Logger

	// r is the router for the application.
	r *mux.Router

	// svr is the server for the application.
	svr *http.Server

	// s is the discord session.
	s *discordgo.Session

	// eventNotifier is the channel for notifying of events.
	eventNotifier chan any
}

// NewApp creates a new instance of App.
func NewApp(l *slog.Logger, r *mux.Router) *App {
	return &App{
		Logger: l,
		r:      r,
	}
}

func (a *App) Run() error {
	// Register bot.
	if err := a.RegisterBot(); err != nil {
		return fmt.Errorf("error registering bot: %w", err)
	}

	a.s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		a.Info(fmt.Sprintf("Logged in as %s#%s", r.User.Username, r.User.Discriminator))
	})

	if err := a.RegisterDiscordHandlers(); err != nil {
		return fmt.Errorf("error registering discord handlers: %w", err)
	}

	// Start event listener.
	go a.eventListener()

	// Open websocket.
	if err := a.s.Open(); err != nil {
		return fmt.Errorf("error opening connection to Discord: %w", err)
	}

	a.Info("Bot is now running.")

	a.generateServer()
	a.setupRoutes()
	a.monitor()

	// Register listerner for shutdown signal.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Process shutdown signal.
	for sig := range c {
		a.Info("Received shutdown signal", slog.String("signal", sig.String()))
		if err := a.ShutdownHook(); err != nil {
			a.Error("Error shutting down application", slog.String(logging.KeyError, err.Error()))
		}
		os.Exit(1)
	}
	return nil
}

func (a *App) ShutdownHook() error {
	// Reset the total number of guilds to 0.
	monitoring.TotalDiscordGuilds.Set(0)

	// Close the connection to Discord.
	if err := a.s.Close(); err != nil {
		return fmt.Errorf("error closing connection to Discord: %w", err)
	}
	return nil
}

func (a *App) RegisterBot() error {
	// Default the number of guilds to 0.
	monitoring.TotalDiscordGuilds.Set(0)

	dg, err := discordgo.New("Bot " + config.BotToken)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	if a.eventNotifier == nil {
		// Create event notifier. This is used to monitor events. It is buffered to prevent blocking.
		a.eventNotifier = make(chan any, 100)
	}

	dg.SetEventNotifier(a.eventNotifier)

	a.s = dg
	return nil
}

func (a *App) monitor() {
	go func() {
		a.Info("Starting monitoring server")
		if err := a.svr.ListenAndServe(); err != nil {
			a.Error("Error starting monitoring server", slog.String(logging.KeyError, err.Error()))
			a.Warn("Monitoring server will not be available")
		}
	}()
}

func (a *App) setupRoutes() {
	// PathMetrics is the path for metrics.
	a.r.HandleFunc(PathMetrics, promhttp.Handler().ServeHTTP).Methods(http.MethodGet)

	// PathHealth is the path for health check.
	a.r.HandleFunc(PathHealth, middlewareHttp(a.healthCheck(), authOptionNone, a)).Methods(http.MethodGet)

	// NotFoundHandler is the handler for 404.
	a.r.NotFoundHandler = middlewareHttp(Controller(request.NotFoundHandler(a.Log())), authOptionNone, a)

	// MethodNotAllowedHandler is the handler for 405.
	a.r.MethodNotAllowedHandler = middlewareHttp(Controller(request.MethodNotAllowedHandler(a.Log())), authOptionNone, a)
}

func (a *App) generateServer() {
	a.svr = &http.Server{
		Addr:    ":" + config.MonitoringPort,
		Handler: a.r,
	}
}

func (a *App) GetGuilds() ([]*discordgo.UserGuild, error) {
	guilds, err := a.s.UserGuilds(0, "", "")
	if err != nil {
		return nil, fmt.Errorf("error getting guilds: %w", err)
	}
	return guilds, nil
}

func (a *App) RegisterDiscordHandlers() error {
	// Bot joined guild.
	a.s.AddHandler(guildJoinedHandler(a))

	// Bot left guild.
	a.s.AddHandler(guildLeaveHandler(a))

	// Slash commands.
	a.s.AddHandler(slashCommandHandler(a, map[string]slashCommandController{}))
	return nil
}

func (a *App) eventListener() {
	for e := range a.eventNotifier {
		switch t := e.(type) {
		case *discordgo.Event:
			if t.Type != "" {
				monitoring.TotalDiscordEvents.WithLabelValues(t.Type).Inc()
			} else {
				// If there is no type, then use the operation name.
				monitoring.TotalDiscordEvents.WithLabelValues(strings.ToUpper(t.Operation.String())).Inc()
			}
		default:
			a.Error("Unknown event type", slog.String("type", fmt.Sprintf("%T", e)))
			monitoring.TotalDiscordEvents.WithLabelValues("UNKNOWN").Inc()
		}
	}
}

func (a *App) Log() *slog.Logger {
	return a.Logger
}

func (a *App) Session() *discordgo.Session {
	return a.s
}
