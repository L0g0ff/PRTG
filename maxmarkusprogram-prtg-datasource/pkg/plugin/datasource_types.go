package plugin

import (
	"sync"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type Datasource struct {
	baseURL       string
	api           PRTGAPI
	logger        PrtgLogger
	tracer        *Tracer
	metrics       *Metrics
	mux           backend.QueryDataHandler
	queryCache    map[string]*QueryCacheEntry
	cacheMutex    sync.RWMutex
	cacheTime     time.Duration
	streamManager *streamManager
}
