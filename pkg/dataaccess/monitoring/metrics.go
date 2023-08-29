package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// MongoLatency is the duration of Mongo queries.
	MongoLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "dataaccess_mongo_latency",
			Help: "Duration of Mongo queries",
		},
		[]string{"dal", "query", "database", "collection"},
	)

	// MongoTotalRequests is the total number of Mongo requests.
	MongoTotalRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dataaccess_mongo_total_requests",
			Help: "Total number of Mongo requests",
		},
		[]string{"dal", "query", "database", "collection"},
	)
)
