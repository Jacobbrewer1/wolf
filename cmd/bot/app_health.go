package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Jacobbrewer1/wolf/pkg/dataaccess"
	dbMonitoring "github.com/Jacobbrewer1/wolf/pkg/dataaccess/monitoring"
	"github.com/alexliesenfeld/health"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
)

func (a *App) healthCheck() Controller {
	checker := health.NewChecker(
		// Set a TTL of 1 second for the results of the checks.
		health.WithCacheDuration(1*time.Second),

		// Set a timeout of 2 seconds for the checks.
		health.WithTimeout(2*time.Second),

		// Monitor the health of the database (MongoDB).
		health.WithCheck(health.Check{
			Name: "MongoDB",
			Check: func(ctx context.Context) error {
				// Create a new timer to measure the latency of the check.
				t := prometheus.NewTimer(dbMonitoring.MongoLatency.WithLabelValues("health_check", "ping", "-", "-"))
				defer t.ObserveDuration()
				dbMonitoring.MongoTotalRequests.WithLabelValues("health_check", "ping", "-", "-").Inc()

				if err := dataaccess.MongoDB.Ping(ctx, nil); err != nil {
					return fmt.Errorf("failed to ping MongoDB: %w", err)
				}
				return nil
			},
			Timeout:            2 * time.Second,
			MaxTimeInError:     0,
			MaxContiguousFails: 0,
			StatusListener: func(ctx context.Context, name string, state health.CheckState) {
				a.Log().Info("MongoDB health check status changed",
					slog.String("name", name),
					slog.String("state", string(state.Status)),
				)
			},
			Interceptors:         nil,
			DisablePanicRecovery: false,
		}),

		// Monitor the health of the Discord API.
		health.WithPeriodicCheck(15*time.Second, 5*time.Second, health.Check{
			Name: "Discord_API",
			Check: func(ctx context.Context) error {
				if _, err := a.Session().GatewayBot(); err != nil {
					return fmt.Errorf("failed to ping Discord API: %w", err)
				}
				return nil
			},
			Timeout:            3 * time.Second,
			MaxTimeInError:     0,
			MaxContiguousFails: 0,
			StatusListener: func(ctx context.Context, name string, state health.CheckState) {
				a.Log().Info("Discord API health check status changed",
					slog.String("name", name),
					slog.String("state", string(state.Status)),
				)
			},
			Interceptors:         nil,
			DisablePanicRecovery: false,
		}),
	)

	return Controller(health.NewHandler(checker))
}
