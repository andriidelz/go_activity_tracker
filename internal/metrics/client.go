package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	EventsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "events_total",
			Help: "Total number of events per action type",
		},
		[]string{"action"},
	)
)

// Register all metrics here.
func Register() {
	prometheus.MustRegister(EventsTotal)
}
