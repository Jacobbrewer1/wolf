package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// TotalDiscordEvents is the total number of events.
	TotalDiscordEvents = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_total_discord_events", AppName),
			Help: "Total number of events",
		},
		[]string{"event"},
	)

	// HttpTotalRequests is the total number of http requests.
	HttpTotalRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_http_total_requests", AppName),
			Help: "Total number of http requests",
		},
		[]string{"path", "method", "status_code"},
	)

	// HttpRequestDuration is the duration of the http request.
	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: fmt.Sprintf("%s_http_request_duration", AppName),
			Help: "Duration of the http request",
		},
		[]string{"path", "method", "status_code"},
	)

	TotalDiscordGuilds = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("%s_total_discord_guilds", AppName),
			Help: "Total number of discord guilds",
		},
	)

	DiscordCommandDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: fmt.Sprintf("%s_discord_command_duration", AppName),
			Help: "Duration of the discord command",
		},
		[]string{"command"},
	)
)
