package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func (d *Datasource) query(ctx context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	var qm struct {
		CacheTime int64 `json:"cacheTime"`
	}
	if err := json.Unmarshal(query.JSON, &qm); err != nil {
		qm.CacheTime = 6000
	}

	cacheKey := QueryCacheKey{
		RefID:     query.RefID,
		QueryType: query.QueryType,
		TimeRange: fmt.Sprintf("%v-%v", query.TimeRange.From.Unix(), query.TimeRange.To.Unix()),
	}

	d.cacheMutex.RLock()
	if cached, exists := d.queryCache[cacheKey.String()]; exists &&
		time.Now().Before(cached.ValidUntil) {
		d.cacheMutex.RUnlock()
		return cached.Response
	}
	d.cacheMutex.RUnlock()

	response := d.executeQuery(ctx, pCtx, query)

	if response.Error == nil {
		d.cacheMutex.Lock()
		d.queryCache[cacheKey.String()] = &QueryCacheEntry{
			Response:   response,
			ValidUntil: time.Now().Add(time.Duration(qm.CacheTime) * time.Millisecond),
		}
		d.cacheMutex.Unlock()
	}

	return response
}

func (d *Datasource) executeQuery(ctx context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	ctx, span := d.tracer.StartSpan(ctx, "query")
	defer span.End()
	backend.Logger.Info("PluginContext", "pCtx", pCtx)

	start := time.Now()
	d.logger.Info("Starting query execution",
		"queryType", query.QueryType,
		"refID", query.RefID,
		"timeRange", fmt.Sprintf("%v to %v", query.TimeRange.From, query.TimeRange.To),
	)

	var qm queryModel
	if err := json.Unmarshal(query.JSON, &qm); err != nil {
		d.logger.Error("Query parsing failed",
			"error", err,
			"raw_query", string(query.JSON),
		)
		d.metrics.IncError("query_parse_error")
		recordError(span, err, "Failed to parse query")
		return backend.ErrDataResponse(backend.StatusBadRequest, "failed to parse query")
	}

	cacheKey := QueryCacheKey{
		RefID:      query.RefID,
		QueryType:  query.QueryType,
		SensorID:   qm.SensorId,
		Channel:    strings.Join(qm.ChannelArray, ","),
		TimeRange:  fmt.Sprintf("%v-%v", query.TimeRange.From.Unix(), query.TimeRange.To.Unix()),
		Property:   qm.Property,
		Parameters: fmt.Sprintf("%s_%s_%s", qm.Group, qm.Device, qm.Sensor),
	}

	cacheTime := d.api.GetCacheTime()
	timeRange := query.TimeRange.To.Sub(query.TimeRange.From)
	var cacheDuration time.Duration

	switch {
	case timeRange <= time.Hour:
		cacheDuration = 6 * time.Second
	case timeRange <= 24*time.Hour:
		cacheDuration = 30 * time.Second
	default:
		cacheDuration = cacheTime
	}

	cacheKeyStr := cacheKey.String()

	d.cacheMutex.RLock()
	if entry, exists := d.queryCache[cacheKeyStr]; exists && time.Now().Before(entry.ValidUntil) {
		d.cacheMutex.RUnlock()
		return entry.Response
	}
	d.cacheMutex.RUnlock()

	addQueryAttributes(span, qm)

	defer func() {
		duration := time.Since(start).Seconds()
		d.metrics.ObserveQueryDuration(qm.QueryType, duration)
		d.logger.Info("Query completed",
			"duration", duration,
			"queryType", qm.QueryType,
			"refID", query.RefID,
		)
	}()

	var response backend.DataResponse
	switch qm.QueryType {
	case "metrics":
		if qm.Channel == "" && len(qm.ChannelArray) == 0 {
			d.logger.Error("Channel selection required for metrics query")
			d.metrics.IncError("missing_channel")
			return backend.ErrDataResponse(backend.StatusBadRequest, "channel selection required")
		}
		response = d.handleMetricsQuery(ctx, qm, query.TimeRange, fmt.Sprintf("metrics_%s", query.RefID))

		if response.Error == nil {
			d.cacheMutex.Lock()
			d.queryCache[cacheKey.String()] = &QueryCacheEntry{
				Response:   response,
				ValidUntil: time.Now().Add(25 * time.Second),
				Updating:   false,
			}
			d.cacheMutex.Unlock()
		}

	case "manual":
		d.logger.Debug("Executing manual query",
			"method", qm.ManualMethod,
			"objectId", qm.ManualObjectId,
		)
		response = d.handleManualQuery(qm, query.TimeRange, fmt.Sprintf("manual_%s", query.RefID))

	case "text", "raw":
		response = d.handlePropertyQuery(ctx, qm, qm.Property, qm.FilterProperty, fmt.Sprintf("property_%s", query.RefID))

	default:
		d.logger.Warn("Unknown query type",
			"type", qm.QueryType,
			"refID", query.RefID,
		)
		d.metrics.IncError("unknown_query_type")
		return backend.DataResponse{
			Frames: []*data.Frame{
				data.NewFrame(fmt.Sprintf("unknown_%s", query.RefID)),
			},
		}
	}

	if response.Error == nil {
		d.cacheMutex.Lock()
		d.queryCache[cacheKeyStr] = &QueryCacheEntry{
			Response:   response,
			ValidUntil: time.Now().Add(cacheDuration),
			Updating:   false,
		}
		d.cacheMutex.Unlock()

		d.logger.Debug("Cached response",
			"key", cacheKeyStr,
			"duration", cacheDuration,
		)
	}

	if response.Error != nil {
		d.logger.Error("Query execution failed",
			"error", response.Error,
			"queryType", qm.QueryType,
			"refID", query.RefID,
		)
		d.metrics.IncError("query_execution")
		recordError(span, response.Error, "Query execution failed")
	}

	return response
}
