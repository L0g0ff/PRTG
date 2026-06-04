package plugin

import (
	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/observability"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics = observability.Metrics
type PrtgLogger = observability.PrtgLogger

var Logger = observability.Logger

func NewLogger() PrtgLogger {
	return observability.NewLogger()
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	return observability.NewMetrics(reg)
}
