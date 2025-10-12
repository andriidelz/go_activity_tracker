package metrics

import (
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v4/cpu"
	"gorm.io/gorm"
)

var (
	CPUUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "app_cpu_usage_percent",
		Help: "Current CPU usage percentage",
	})

	MemoryUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "app_memory_usage_bytes",
		Help: "Current memory usage in bytes",
	})

	Goroutines = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "app_goroutines_total",
		Help: "Number of active goroutines",
	})

	DBOpenConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_open_connections",
		Help: "Number of open database connections",
	})

	DBInUseConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_in_use_connections",
		Help: "Number of connections currently in use",
	})

	DBIdleConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_idle_connections",
		Help: "Number of idle connections in the pool",
	})
)

func CollectSystemMetrics(db *gorm.DB) {
	go func() {
		for {
			updateMetrics(db)
			time.Sleep(5 * time.Second)
		}
	}()
}

func updateMetrics(db *gorm.DB) {
	cpuPercent, _ := cpu.Percent(0, false)
	if len(cpuPercent) > 0 {
		CPUUsage.Set(cpuPercent[0])
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	MemoryUsage.Set(float64(m.Alloc))
	Goroutines.Set(float64(runtime.NumGoroutine()))

	sqlDB, err := db.DB()
	if err == nil {
		stats := sqlDB.Stats()
		DBOpenConnections.Set(float64(stats.OpenConnections))
		DBInUseConnections.Set(float64(stats.InUse))
		DBIdleConnections.Set(float64(stats.Idle))
	}
}
