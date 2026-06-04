package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/concurrent"
)

const MaxConcurrentQueries = 10

func (d *Datasource) handleSingleQueryData(ctx context.Context, q concurrent.Query) backend.DataResponse {
	ctx, span := d.tracer.StartSpan(ctx, "handleSingleQueryData")
	defer span.End()

	start := time.Now()
	d.logger.Debug("Processing concurrent query",
		"refId", q.DataQuery.RefID,
	)

	res := d.query(ctx, q.PluginContext, q.DataQuery)

	duration := time.Since(start)
	d.metrics.ObserveQueryDuration("single_query", duration.Seconds())
	d.logger.Debug("Concurrent query processed",
		"refId", q.DataQuery.RefID,
		"duration", duration,
		"status", res.Status,
		"hasError", res.Error != nil,
	)

	return res
}

func generateCacheKey(req *backend.QueryDataRequest) string {
	var keyBuilder strings.Builder
	for _, q := range req.Queries {
		keyBuilder.WriteString(fmt.Sprintf("%s:%d:%d:%s;",
			q.RefID,
			q.TimeRange.From.UnixNano(),
			q.TimeRange.To.UnixNano(),
			string(q.JSON)))
	}
	return keyBuilder.String()
}

func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	if len(req.Queries) == 0 {
		return &backend.QueryDataResponse{
			Responses: make(map[string]backend.DataResponse),
		}, nil
	}

	if len(req.Queries) > MaxConcurrentQueries {
		return &backend.QueryDataResponse{
			Responses: map[string]backend.DataResponse{
				req.Queries[0].RefID: {
					Error:  fmt.Errorf("query limit exceeded: %d/%d", len(req.Queries), MaxConcurrentQueries),
					Status: backend.StatusTooManyRequests,
				},
			},
		}, nil
	}

	cacheKey := generateCacheKey(req)
	d.cacheMutex.RLock()
	if cached, exists := d.queryCache[cacheKey]; exists && time.Now().Before(cached.ValidUntil) {
		d.cacheMutex.RUnlock()
		response := backend.NewQueryDataResponse()
		response.Responses[req.Queries[0].RefID] = cached.Response
		return response, nil
	}
	d.cacheMutex.RUnlock()

	response, err := d.mux.QueryData(ctx, req)
	if err != nil {
		return nil, err
	}

	d.cacheMutex.Lock()
	d.queryCache[cacheKey] = &QueryCacheEntry{
		Response:   response.Responses[req.Queries[0].RefID],
		ValidUntil: time.Now().Add(d.cacheTime),
		Updating:   false,
	}
	d.cacheMutex.Unlock()

	return response, nil
}

func (d *Datasource) handleMetricsQueryType(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	ctx, span := d.tracer.StartSpan(ctx, "handleMetricsQueryType")
	defer span.End()

	response := backend.NewQueryDataResponse()

	for _, q := range req.Queries {
		response.Responses[q.RefID] = d.handleSingleQueryData(ctx, concurrent.Query{
			DataQuery:     q,
			PluginContext: req.PluginContext,
		})
	}

	return response, nil
}

func (d *Datasource) handleManualQueryType(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	_, span := d.tracer.StartSpan(ctx, "handleManualQueryType")
	defer span.End()

	response := backend.NewQueryDataResponse()

	for _, q := range req.Queries {
		var qm queryModel
		if err := json.Unmarshal(q.JSON, &qm); err != nil {
			response.Responses[q.RefID] = backend.ErrDataResponse(backend.StatusBadRequest, "failed to parse query")
			continue
		}

		response.Responses[q.RefID] = d.handleManualQuery(qm, q.TimeRange, fmt.Sprintf("manual_%s", q.RefID))
	}

	return response, nil
}

func (d *Datasource) handlePropertyQueryType(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	ctx, span := d.tracer.StartSpan(ctx, "handlePropertyQueryType")
	defer span.End()

	response := backend.NewQueryDataResponse()

	for _, q := range req.Queries {
		var qm queryModel
		if err := json.Unmarshal(q.JSON, &qm); err != nil {
			response.Responses[q.RefID] = backend.ErrDataResponse(backend.StatusBadRequest, "failed to parse query")
			continue
		}

		response.Responses[q.RefID] = d.handlePropertyQuery(
			ctx,
			qm,
			qm.Property,
			qm.FilterProperty,
			fmt.Sprintf("property_%s", q.RefID),
		)
	}

	return response, nil
}

func (d *Datasource) handleQueryFallback(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	d.logger.Warn("Query type not supported", "queries", len(req.Queries))
	return backend.NewQueryDataResponse(), nil
}
