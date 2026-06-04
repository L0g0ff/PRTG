package plugin

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	_ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
	_ backend.CallResourceHandler   = (*Datasource)(nil)
	_ backend.StreamHandler         = (*Datasource)(nil)
)

func NewDatasource(ctx context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	config, err := models.LoadPluginSettings(settings)
	if err != nil {
		return nil, err
	}

	SetDefaultTimezone(config.Timezone)

	var cacheTime time.Duration = 60 * time.Second
	if config.CacheTime > 0 {
		cacheTime = config.CacheTime * time.Second
	}

	baseURL := fmt.Sprintf("https://%s", config.Path)
	logger := NewLogger()
	tracer := NewTracer(logger)
	metrics := NewMetrics(prometheus.DefaultRegisterer)

	ds := &Datasource{
		baseURL:    baseURL,
		api:        NewApi(baseURL, config.Secrets.ApiKey, cacheTime, 10*time.Second),
		logger:     logger,
		tracer:     tracer,
		metrics:    metrics,
		queryCache: make(map[string]*QueryCacheEntry),
		cacheMutex: sync.RWMutex{},
		cacheTime:  cacheTime,
		streamManager: &streamManager{
			streams:          make(map[string]*activeStream),
			activeStreams:    make(map[string]map[string]*activeStream),
			defaultCacheTime: cacheTime,
		},
	}

	queryTypeMux := datasource.NewQueryTypeMux()
	queryTypeMux.HandleFunc("metrics", ds.handleMetricsQueryType)
	queryTypeMux.HandleFunc("manual", ds.handleManualQueryType)
	queryTypeMux.HandleFunc("text", ds.handlePropertyQueryType)
	queryTypeMux.HandleFunc("raw", ds.handlePropertyQueryType)
	queryTypeMux.HandleFunc("", ds.handleQueryFallback)

	ds.mux = queryTypeMux
	return ds, nil
}

func (d *Datasource) Dispose() {
	d.cacheMutex.Lock()
	d.queryCache = make(map[string]*QueryCacheEntry)
	d.cacheMutex.Unlock()

	if apiImpl, ok := d.api.(*Api); ok {
		apiImpl.ClearCache()
	}
}

func (d *Datasource) ClearAllCaches() {
	d.cacheMutex.Lock()
	d.queryCache = make(map[string]*QueryCacheEntry)
	d.cacheMutex.Unlock()

	if apiImpl, ok := d.api.(*Api); ok {
		apiImpl.ClearCache()
	}

	d.logger.Debug("All caches cleared")
}
