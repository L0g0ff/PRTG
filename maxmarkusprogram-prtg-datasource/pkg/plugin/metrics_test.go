package plugin

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestNewMetricsReusesAlreadyRegisteredCollectors(t *testing.T) {
	registry := prometheus.NewRegistry()

	first := NewMetrics(registry)
	second := NewMetrics(registry)

	if first.apiRequests != second.apiRequests {
		t.Fatal("expected api request counter to be reused")
	}
	if first.apiLatency != second.apiLatency {
		t.Fatal("expected api latency histogram to be reused")
	}
	if first.queryDuration != second.queryDuration {
		t.Fatal("expected query duration histogram to be reused")
	}
	if first.cacheHits != second.cacheHits {
		t.Fatal("expected cache hit counter to be reused")
	}
	if first.errorCounter != second.errorCounter {
		t.Fatal("expected error counter to be reused")
	}

	first.IncAPIRequest("getstatus.htm")
	second.ObserveAPILatency("getstatus.htm", 0.1)
	first.ObserveQueryDuration("health_check", 0.1)
	second.IncCacheHit("query")
	first.IncError("health_check")

	if _, err := registry.Gather(); err != nil {
		t.Fatalf("gather metrics: %v", err)
	}
}
