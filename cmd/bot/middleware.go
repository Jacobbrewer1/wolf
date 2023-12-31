package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/Jacobbrewer1/discordgo"
	"github.com/Jacobbrewer1/wolf/pkg/logging"
	"github.com/Jacobbrewer1/wolf/pkg/messages"
	"github.com/Jacobbrewer1/wolf/pkg/request"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

// commandController is the handler for slash commands.
type commandController func(a IApp, i *discordgo.InteractionCreate) (commandProcessor, error)

// commandProcessor is the processor for slash commands.
type commandProcessor func(a IApp, i *discordgo.InteractionCreate) error

type Controller func(w http.ResponseWriter, r *http.Request)

func middlewareHttp(handler Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().UTC()
		cw := request.NewClientWriter(w)

		// Recover from any panics that occur in the handler.
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("Panic in handler",
					slog.String(logging.KeyError, rec.(error).Error()),
					slog.String("stack", string(debug.Stack())),
				)
				cw.WriteHeader(http.StatusInternalServerError)
				if err := json.NewEncoder(cw).Encode(request.NewMessage(messages.ErrInternalServerError)); err != nil {
					slog.Error("Error encoding response", slog.String(logging.KeyError, err.Error()))
				}
			}
		}()

		var path string
		route := mux.CurrentRoute(r)
		if route != nil { // The route may be nil if the request is not routed.
			var err error
			path, err = route.GetPathTemplate()
			if err != nil {
				// An error here is only returned if the route does not define a path.
				slog.Error("Error getting path template", slog.String(logging.KeyError, err.Error()))
				path = r.URL.Path // If the route does not define a path, use the URL path.
			}
		} else {
			path = r.URL.Path // If the route is nil, use the URL path.
		}

		defer func() {
			// Run the deferred function after the request has been handled, as the status code will not be available until then.
			HttpTotalRequests.WithLabelValues(path, r.Method, fmt.Sprintf("%d", cw.StatusCode())).Inc()
			HttpRequestDuration.WithLabelValues(path, r.Method, fmt.Sprintf("%d", cw.StatusCode())).Observe(time.Since(now).Seconds())
		}()

		handler(cw, r)
	}
}

// interactionHandler is the handler for interactions.
func interactionHandler(
	a IApp,
	slashControllers map[string]commandController,
	buttonControllers map[string]commandProcessor,
) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// Process the latency for the interaction.
		t := prometheus.NewTimer(DiscordCommandDuration.WithLabelValues(i.Type.String()))
		defer t.ObserveDuration()

		switch i.Type {
		// Slash commands.
		case discordgo.InteractionApplicationCommand:
			slashCommandHandler(a, slashControllers)(s, i)
		// Button interactions.
		case discordgo.InteractionMessageComponent:
			buttonHandler(a, buttonControllers)(s, i)
		// Unknown interaction type.
		default:
			slog.Error(fmt.Sprintf("Unknown interaction type %d", i.Type),
				slog.Int("type", int(i.Type)))

			if err := respondEphemeralError(a, i); err != nil {
				slog.Error("Error responding to interaction", slog.String(logging.KeyError, err.Error()))
				return
			}
		}
	}
}

// slashCommandHandler is the handler for slash commands.
func slashCommandHandler(a IApp, controllers map[string]commandController) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(_ *discordgo.Session, i *discordgo.InteractionCreate) {
		slog.Debug("Handling interaction " + i.ApplicationCommandData().Name)
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		if controller, ok := controllers[i.ApplicationCommandData().Name]; ok {
			processor, err := controller(a, i)
			if err != nil {
				slog.Error(fmt.Sprintf("Error getting processor for command %s", i.ApplicationCommandData().Name),
					slog.String(logging.KeyError, err.Error()))

				if err := respondEphemeralError(a, i); err != nil {
					slog.Error("Error responding to interaction", slog.String(logging.KeyError, err.Error()))
					return
				}
				return
			}

			if err := processor(a, i); err != nil {
				slog.Error(fmt.Sprintf("Error processing command %s", i.ApplicationCommandData().Name),
					slog.String(logging.KeyError, err.Error()))

				if err := respondEphemeralError(a, i); err != nil {
					slog.Error("Error responding to interaction", slog.String(logging.KeyError, err.Error()))
					return
				}
				return
			}
		} else {
			slog.Error(fmt.Sprintf("No controller found for command %s", i.ApplicationCommandData().Name),
				slog.String("command", i.ApplicationCommandData().Name))

			if err := respondEphemeralError(a, i); err != nil {
				slog.Error("Error responding to interaction", slog.String(logging.KeyError, err.Error()))
				return
			}
		}
	}
}

// buttonHandler is the handler for button interactions.
func buttonHandler(a IApp, controllers map[string]commandProcessor) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(_ *discordgo.Session, i *discordgo.InteractionCreate) {
		slog.Debug("Handling interaction " + i.MessageComponentData().CustomID)
		if i.Type != discordgo.InteractionMessageComponent {
			return
		}

		if processor, ok := controllers[i.MessageComponentData().CustomID]; ok {
			if err := processor(a, i); err != nil {
				slog.Error(fmt.Sprintf("Error processing command %s", i.MessageComponentData().CustomID),
					slog.String(logging.KeyError, err.Error()))

				if err := respondEphemeralError(a, i); err != nil {
					slog.Error("Error responding to interaction", slog.String(logging.KeyError, err.Error()))
					return
				}
				return
			}
		} else {
			slog.Error(fmt.Sprintf("No controller found for command %s", i.MessageComponentData().CustomID),
				slog.String("command", i.MessageComponentData().CustomID))

			if err := respondEphemeralError(a, i); err != nil {
				slog.Error("Error responding to interaction", slog.String(logging.KeyError, err.Error()))
				return
			}
		}
	}
}
