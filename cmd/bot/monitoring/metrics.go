package monitoring

import (
	"fmt"

	"github.com/Jacobbrewer1/wolf/cmd/bot/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// TotalDiscordEvents is the total number of events.
	TotalDiscordEvents = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_total_discord_events", config.AppName),
			Help: "Total number of events",
		},
		[]string{"event"},
	)

	// HttpTotalRequests is the total number of http requests.
	HttpTotalRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_http_total_requests", config.AppName),
			Help: "Total number of http requests",
		},
		[]string{"path", "method", "status_code"},
	)

	// HttpRequestDuration is the duration of the http request.
	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: fmt.Sprintf("%s_http_request_duration", config.AppName),
			Help: "Duration of the http request",
		},
		[]string{"path", "method", "status_code"},
	)

	TotalDiscordGuilds = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("%s_total_discord_guilds", config.AppName),
			Help: "Total number of discord guilds",
		},
	)
)
