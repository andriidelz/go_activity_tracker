package metrics

import "github.com/prometheus/client_golang/prometheus"

func Register() {
	prometheus.MustRegister(
		EventsTotal,

		CPUUsage,
		MemoryUsage,
		Goroutines,
		DBOpenConnections,
		DBInUseConnections,
		DBIdleConnections,
	)
}
