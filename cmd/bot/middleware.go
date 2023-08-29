package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/Jacobbrewer1/discordgo"
	"github.com/Jacobbrewer1/wolf/cmd/bot/monitoring"
	"github.com/Jacobbrewer1/wolf/pkg/logging"
	"github.com/Jacobbrewer1/wolf/pkg/request"
	"github.com/gorilla/mux"
	"golang.org/x/exp/slog"
)

// slashCommandController is the handler for slash commands.
type slashCommandController func(a IApp, cmd string) (slashProcessor, error)

// slashProcessor is the processor for slash commands.
type slashProcessor func(a IApp, i *discordgo.InteractionCreate) error

// authOption is an option for the auth middleware. It indicates the type of authentication required.
type authOption int

const (
	// authOptionNone indicates that no authentication is required.
	authOptionNone authOption = iota
)

type Controller func(w http.ResponseWriter, r *http.Request)

func middlewareHttp(handler Controller, authRequired authOption, a IApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().UTC()
		cw := request.NewClientWriter(w)

		// Recover from any panics that occur in the handler.
		defer func() {
			if rec := recover(); rec != nil {
				a.Log().Error("Panic in handler",
					slog.String(logging.KeyError, rec.(error).Error()),
					slog.String("stack", string(debug.Stack())),
				)
				cw.WriteHeader(http.StatusInternalServerError)
				if err := json.NewEncoder(cw).Encode(request.NewMessage(request.ErrInternalServer.Error())); err != nil {
					a.Log().Error("Error encoding response", slog.String(logging.KeyError, err.Error()))
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
				a.Log().Error("Error getting path template", slog.String(logging.KeyError, err.Error()))
				path = r.URL.Path // If the route does not define a path, use the URL path.
			}
		} else {
			path = r.URL.Path // If the route is nil, use the URL path.
		}

		defer func() {
			// Run the deferred function after the request has been handled, as the status code will not be available until then.
			monitoring.HttpTotalRequests.WithLabelValues(path, r.Method, fmt.Sprintf("%d", cw.StatusCode())).Inc()
			monitoring.HttpRequestDuration.WithLabelValues(path, r.Method, fmt.Sprintf("%d", cw.StatusCode())).Observe(time.Since(now).Seconds())
		}()

		handler(cw, r)
	}
}

// slashCommandHandler is the handler for slash commands.
func slashCommandHandler(a IApp, controllers map[string]slashCommandController) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(_ *discordgo.Session, i *discordgo.InteractionCreate) {
		a.Log().Debug("Handling interaction " + i.ApplicationCommandData().Name)
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		if controller, ok := controllers[i.ApplicationCommandData().Name]; ok {
			processor, err := controller(a, i.Interaction.Data.(discordgo.ApplicationCommandInteractionData).Options[0].Name)
			if err != nil {
				a.Log().Error(fmt.Sprintf("Error getting processor for command %s", i.ApplicationCommandData().Name),
					slog.String(logging.KeyError, err.Error()))

				if err := respondSlashError(a, i); err != nil {
					a.Log().Error("Error responding to interaction", slog.String(logging.KeyError, err.Error()))
					return
				}
				return
			}

			if err := processor(a, i); err != nil {
				a.Log().Error(fmt.Sprintf("Error processing command %s", i.ApplicationCommandData().Name),
					slog.String(logging.KeyError, err.Error()))

				if err := respondSlashError(a, i); err != nil {
					a.Log().Error("Error responding to interaction", slog.String(logging.KeyError, err.Error()))
					return
				}
				return
			}
		} else {
			a.Log().Error(fmt.Sprintf("No controller found for command %s", i.ApplicationCommandData().Name),
				slog.String("command", i.ApplicationCommandData().Name))

			if err := respondSlashError(a, i); err != nil {
				a.Log().Error("Error responding to interaction", slog.String(logging.KeyError, err.Error()))
				return
			}
		}
	}
}
