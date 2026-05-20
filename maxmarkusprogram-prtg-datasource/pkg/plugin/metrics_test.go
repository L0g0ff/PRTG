package plugin

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestNewMetricsReusesRegisteredCollectors(t *testing.T) {
	registry := prometheus.NewRegistry()

	first := NewMetrics(registry)
	second := NewMetrics(registry)

	if first.apiRequests != second.apiRequests {
		t.Fatal("expected apiRequests collector to be reused")
	}
	if first.apiLatency != second.apiLatency {
		t.Fatal("expected apiLatency collector to be reused")
	}
	if first.queryDuration != second.queryDuration {
		t.Fatal("expected queryDuration collector to be reused")
	}
	if first.cacheHits != second.cacheHits {
		t.Fatal("expected cacheHits collector to be reused")
	}
	if first.errorCounter != second.errorCounter {
		t.Fatal("expected errorCounter collector to be reused")
	}
}
