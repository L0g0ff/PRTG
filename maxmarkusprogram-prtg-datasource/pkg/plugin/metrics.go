package plugin

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

/* =================================== METRICS DEFINITIONS ===================================== */

var (
	// Active connections metric - now used with logging
	activeConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "grafana_plugin",
			Name:      "prtg_active_connections",
			Help:      "Current number of active connections.",
		},
	)
)

/* =================================== METRICS STRUCT ======================================== */
type Metrics struct {
	apiRequests   *prometheus.CounterVec
	apiLatency    *prometheus.HistogramVec
	queryDuration *prometheus.HistogramVec
	cacheHits     *prometheus.CounterVec
	errorCounter  *prometheus.CounterVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	if reg == nil {
		reg = prometheus.DefaultRegisterer
	}

	m := &Metrics{
		apiRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "prtg_api_requests_total",
				Help: "Total number of API requests made to PRTG",
			},
			[]string{"endpoint"},
		),
		apiLatency: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "prtg_api_request_duration_seconds",
				Help: "Duration of API requests to PRTG",
			},
			[]string{"endpoint"},
		),
		queryDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "prtg_query_duration_seconds",
				Help: "Duration of PRTG queries",
			},
			[]string{"query_type"},
		),
		cacheHits: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "prtg_cache_hits_total",
				Help: "Total number of cache hits",
			},
			[]string{"type"},
		),
		errorCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "prtg_errors_total",
				Help: "Total number of errors",
			},
			[]string{"type"},
		),
	}

	m.apiRequests = registerCounterVec(reg, m.apiRequests)
	m.apiLatency = registerHistogramVec(reg, m.apiLatency)
	m.queryDuration = registerHistogramVec(reg, m.queryDuration)
	m.cacheHits = registerCounterVec(reg, m.cacheHits)
	m.errorCounter = registerCounterVec(reg, m.errorCounter)

	return m
}

func registerCounterVec(reg prometheus.Registerer, collector *prometheus.CounterVec) *prometheus.CounterVec {
	if err := reg.Register(collector); err != nil {
		existing := alreadyRegisteredCollector(err)
		if counter, ok := existing.(*prometheus.CounterVec); ok {
			return counter
		}
	}

	return collector
}

func registerHistogramVec(reg prometheus.Registerer, collector *prometheus.HistogramVec) *prometheus.HistogramVec {
	if err := reg.Register(collector); err != nil {
		existing := alreadyRegisteredCollector(err)
		if histogram, ok := existing.(*prometheus.HistogramVec); ok {
			return histogram
		}
	}

	return collector
}

func alreadyRegisteredCollector(err error) prometheus.Collector {
	var alreadyRegistered prometheus.AlreadyRegisteredError
	if errors.As(err, &alreadyRegistered) {
		return alreadyRegistered.ExistingCollector
	}

	return nil
}

func (m *Metrics) IncAPIRequest(endpoint string) {
	m.apiRequests.WithLabelValues(endpoint).Inc()
}

func (m *Metrics) ObserveAPILatency(endpoint string, duration float64) {
	m.apiLatency.WithLabelValues(endpoint).Observe(duration)
}

func (m *Metrics) ObserveQueryDuration(queryType string, duration float64) {
	m.queryDuration.WithLabelValues(queryType).Observe(duration)
}

func (m *Metrics) IncCacheHit(type_ string) {
	m.cacheHits.WithLabelValues(type_).Inc()
}

func (m *Metrics) IncError(type_ string) {
	m.errorCounter.WithLabelValues(type_).Inc()
}

// Add this method to the Metrics struct
func (m *Metrics) UpdateActiveConnections(count float64, logger PrtgLogger) {
	activeConnections.Set(count)
	logger.Debug("Updated active connections metric",
		"count", count,
		"metric", "prtg_active_connections",
	)
}
