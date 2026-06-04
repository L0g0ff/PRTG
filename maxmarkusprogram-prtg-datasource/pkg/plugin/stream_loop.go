package plugin

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func (d *Datasource) runStreamLoop(
	ctx context.Context,
	stream *activeStream,
	query queryModel,
	sender *backend.StreamSender,
	timeRangeFrom, timeRangeTo time.Time,
) error {
	defer d.cleanupStream(stream)

	timeRange := backend.TimeRange{
		From: timeRangeFrom,
		To:   timeRangeTo,
	}

	jitter := time.Duration(rand.Int63n(250)) * time.Millisecond
	ticker := time.NewTicker(stream.interval + jitter)
	defer ticker.Stop()

	if err := d.updateStreamWithMetricsQuery(ctx, stream, sender, query, timeRange); err != nil {
		d.logger.Error("Initial stream update failed", "error", err)
	}

	return d.streamUpdateLoop(ctx, stream, sender, query, timeRange, ticker)
}

func (d *Datasource) streamUpdateLoop(
	ctx context.Context,
	stream *activeStream,
	sender *backend.StreamSender,
	query queryModel,
	timeRange backend.TimeRange,
	ticker *time.Ticker,
) error {
	for {
		select {
		case <-ctx.Done():
			stream.isActive = false
			return nil

		case <-ticker.C:
			if stream.updateMode == "sliding" {
				now := time.Now()
				windowSize := timeRange.To.Sub(timeRange.From)
				timeRange.From = now.Add(-windowSize)
				timeRange.To = now
			}

			updateCtx, cancel := context.WithTimeout(ctx, stream.interval/2)
			err := d.updateStreamWithMetricsQuery(updateCtx, stream, sender, query, timeRange)
			cancel()

			d.handleStreamError(stream, err)

		case <-stream.updateChan:
			if err := d.updateStreamWithMetricsQuery(ctx, stream, sender, query, timeRange); err != nil {
				d.logger.Error("Manual update failed", "error", err)
			}
		}
	}
}

func (d *Datasource) handleStreamError(stream *activeStream, err error) {
	if err != nil {
		stream.errorCount++
		if stream.errorCount <= 3 || stream.errorCount%10 == 0 {
			d.logger.Error("Stream update failed",
				"error", err,
				"count", stream.errorCount,
				"streamID", stream.streamID)
		}
	} else {
		stream.errorCount = 0
	}
}

func (d *Datasource) updateStreamWithMetricsQuery(
	ctx context.Context,
	stream *activeStream,
	sender *backend.StreamSender,
	query queryModel,
	timeRange backend.TimeRange,
) error {
	if time.Since(stream.lastUpdate) < stream.cacheTime {
		return nil
	}

	stream.status.updating = true
	defer func() { stream.status.updating = false }()

	stream.lastUpdate = time.Now()

	baseFrameName := fmt.Sprintf("stream_%s", stream.refID)
	response := d.handleMetricsQuery(ctx, query, timeRange, baseFrameName)

	if response.Error != nil {
		stream.status.lastError = response.Error
		return response.Error
	}

	if len(response.Frames) == 0 {
		return nil
	}

	return d.processResponseFrames(stream, sender, response, timeRange)
}

func (d *Datasource) processResponseFrames(
	stream *activeStream,
	sender *backend.StreamSender,
	response backend.DataResponse,
	timeRange backend.TimeRange,
) error {
	for _, frame := range response.Frames {
		if len(frame.Fields) < 2 {
			continue
		}

		channelName := extractChannelName(frame)
		if channelName == "" {
			continue
		}

		channelState, exists := stream.channelStates[channelName]
		if !exists {
			continue
		}

		times, values := extractFrameData(frame)
		if len(times) == 0 {
			continue
		}

		updateChannelBuffer(stream, channelState, times, values)

		streamFrame := createStreamingFrame(stream, channelName, channelState, timeRange.From, timeRange.To)

		if err := sender.SendFrame(streamFrame, data.IncludeAll); err != nil {
			d.logger.Error("Failed to send frame", "error", err, "channel", channelName)
			continue
		}
	}

	stream.lastDataTimestamp = time.Now().UnixMilli()
	return nil
}
