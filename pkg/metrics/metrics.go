package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// PrometheusService struct containing prometheus counter(s)
type PrometheusService struct {
	TotalRequests *prometheus.CounterVec
}

var totalRequests = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "go_app",
	Name:      "server_rps",
	Help:      "number of total requests",
},
	[]string{"method", "path"})

// NewPrometheusService returns a new service with registered prometheus counters
func NewPrometheusService() *PrometheusService {
	return &PrometheusService{
		TotalRequests: totalRequests,
	}
}
