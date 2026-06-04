package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

const (
	DefaultBufferSize = 500

	MaxStreamsPerPanel = 5

	MinUpdateInterval = 1000 * time.Millisecond
	MaxUpdateInterval = 60000 * time.Millisecond
	DefaultInterval   = 5000 * time.Millisecond

	DefaultCacheTime = 5 * time.Second

	InitialBufferCapacity = 32
)

func (d *Datasource) SubscribeStream(ctx context.Context, req *backend.SubscribeStreamRequest) (*backend.SubscribeStreamResponse, error) {
	d.logger.Debug("Subscribe to stream", "path", req.Path)

	if !strings.HasPrefix(req.Path, "prtg-stream/") {
		return &backend.SubscribeStreamResponse{Status: backend.SubscribeStreamStatusNotFound}, nil
	}

	var query queryModel
	if err := json.Unmarshal(req.Data, &query); err != nil {
		d.logger.Error("Invalid subscription data", "error", err)
		return &backend.SubscribeStreamResponse{Status: backend.SubscribeStreamStatusPermissionDenied}, nil
	}

	if query.SensorId == "" || (len(query.ChannelArray) == 0 && query.Channel == "") {
		d.logger.Error("Missing required fields", "sensorId", query.SensorId)
		return &backend.SubscribeStreamResponse{Status: backend.SubscribeStreamStatusPermissionDenied}, nil
	}

	panelId := fmt.Sprintf("%v", query.PanelID)
	if panelStreams := d.getStreamsByPanel(panelId); len(panelStreams) >= MaxStreamsPerPanel {
		d.logger.Warn("Maximum streams reached for panel", "panelId", panelId)
		return &backend.SubscribeStreamResponse{Status: backend.SubscribeStreamStatusPermissionDenied}, nil
	}

	if timeRangeInfo, err := extractTimeRangeInfo(req.Data); err == nil {
		d.logger.Debug("Stream time range",
			"windowSize", time.Duration(timeRangeInfo.To-timeRangeInfo.From)*time.Millisecond)
	}

	return &backend.SubscribeStreamResponse{Status: backend.SubscribeStreamStatusOK}, nil
}

func extractTimeRangeInfo(data []byte) (struct{ From, To int64 }, error) {
	var result struct {
		TimeRange struct {
			From int64 `json:"from"`
			To   int64 `json:"to"`
		} `json:"timeRange"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return struct{ From, To int64 }{0, 0}, err
	}

	return struct{ From, To int64 }{
		From: result.TimeRange.From,
		To:   result.TimeRange.To,
	}, nil
}

func (d *Datasource) PublishStream(ctx context.Context, req *backend.PublishStreamRequest) (*backend.PublishStreamResponse, error) {
	return &backend.PublishStreamResponse{Status: backend.PublishStreamStatusPermissionDenied}, nil
}

func (d *Datasource) RunStream(ctx context.Context, req *backend.RunStreamRequest, sender *backend.StreamSender) error {
	var query queryModel
	if err := json.Unmarshal(req.Data, &query); err != nil {
		return fmt.Errorf("failed to parse stream data: %w", err)
	}

	if query.SensorId == "" {
		return fmt.Errorf("missing required field: sensorId")
	}

	channels := getChannels(query)
	if len(channels) == 0 {
		return fmt.Errorf("missing required field: channel or channelArray")
	}

	interval := getBoundedInterval(query.StreamInterval)
	streamID := generateStreamID(query, channels)
	timeRangeFrom, timeRangeTo := getTimeRange(query)
	cacheDuration := getCacheDuration(query)
	bufferSize := getBufferSize(query)

	stream := createStream(query, channels, streamID, timeRangeFrom, timeRangeTo,
		interval, cacheDuration, bufferSize)

	d.logger.Info("Stream starting",
		"streamID", streamID,
		"intervalMs", interval.Milliseconds())

	if existingStream := d.getExistingStream(streamID); existingStream != nil {
		d.updateExistingStream(existingStream, timeRangeFrom, timeRangeTo, cacheDuration)
		return nil
	}

	d.registerNewStream(stream, streamID)

	return d.runStreamLoop(ctx, stream, query, sender, timeRangeFrom, timeRangeTo)
}

func getChannels(query queryModel) []string {
	if len(query.ChannelArray) > 0 {
		return query.ChannelArray
	}
	if query.Channel != "" {
		return []string{query.Channel}
	}
	return nil
}

func getBoundedInterval(requestedInterval int64) time.Duration {
	if requestedInterval <= 0 {
		return DefaultInterval
	}

	interval := time.Duration(requestedInterval) * time.Millisecond
	if interval < MinUpdateInterval {
		return MinUpdateInterval
	}
	if interval > MaxUpdateInterval {
		return MaxUpdateInterval
	}
	return interval
}

func generateStreamID(query queryModel, channels []string) string {
	channelKey := strings.Join(channels, "_")
	return fmt.Sprintf("%v_%s_%s_%s",
		query.PanelID,
		query.RefID,
		query.SensorId,
		channelKey)
}

func getTimeRange(query queryModel) (time.Time, time.Time) {
	now := time.Now()
	if query.From > 0 && query.To > 0 {
		return time.Unix(0, query.From*int64(time.Millisecond)),
			time.Unix(0, query.To*int64(time.Millisecond))
	}
	return now.Add(-30 * time.Minute), now
}

func getCacheDuration(query queryModel) time.Duration {
	if query.StreamInterval > 0 {
		cacheDuration := time.Duration(query.StreamInterval/2) * time.Millisecond
		if cacheDuration < time.Second {
			return time.Second
		}
		return cacheDuration
	}
	return DefaultCacheTime
}

func getBufferSize(_ queryModel) int64 {
	return DefaultBufferSize
}

func createStream(query queryModel, channels []string, streamID string,
	fromTime, toTime time.Time, interval, cacheDuration time.Duration, bufferSize int64) *activeStream {

	stream := &activeStream{
		sensorId:          query.SensorId,
		channelArray:      channels,
		interval:          interval,
		group:             query.Group,
		device:            query.Device,
		sensor:            query.Sensor,
		includeGroupName:  query.IncludeGroupName,
		includeDeviceName: query.IncludeDeviceName,
		includeSensorName: query.IncludeSensorName,
		fromTime:          fromTime,
		toTime:            toTime,
		lastUpdate:        time.Now().Add(-cacheDuration),
		cacheTime:         cacheDuration,
		isActive:          true,
		refID:             query.RefID,
		streamID:          streamID,
		panelId:           fmt.Sprintf("%v", query.PanelID),
		queryId:           query.RefID,
		multiChannelKey:   strings.Join(channels, "_"),
		channelStates:     make(map[string]*channelState, len(channels)),
		updateChan:        make(chan struct{}, 1),
		updateMode:        query.UpdateMode,
		bufferSize:        bufferSize,
		status: &streamStatus{
			active:    true,
			updating:  false,
			lastError: nil,
		},
		lastDataTimestamp: time.Now().UnixMilli(),
	}

	if stream.updateMode == "" {
		stream.updateMode = "append"
	}

	for _, channelName := range channels {
		stream.channelStates[channelName] = &channelState{
			lastValue: 0,
			isActive:  true,
			buffer: &dataBuffer{
				times:  make([]time.Time, 0, InitialBufferCapacity),
				values: make([]float64, 0, InitialBufferCapacity),
				size:   bufferSize,
			},
		}
	}

	return stream
}
