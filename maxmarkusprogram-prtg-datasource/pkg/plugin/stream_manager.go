package plugin

import (
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (d *Datasource) getExistingStream(streamID string) *activeStream {
	d.streamManager.mu.RLock()
	existingStream, exists := d.streamManager.streams[streamID]
	d.streamManager.mu.RUnlock()

	if !exists {
		return nil
	}
	return existingStream
}

func (d *Datasource) updateExistingStream(stream *activeStream, from time.Time, to time.Time, cacheDuration time.Duration) {
	d.streamManager.mu.Lock()
	stream.timeRange = &backend.TimeRange{From: from, To: to}
	stream.isActive = true
	stream.lastUpdate = time.Now().Add(-cacheDuration)
	d.streamManager.mu.Unlock()

	select {
	case stream.updateChan <- struct{}{}:
	default:
	}
}

func (d *Datasource) registerNewStream(stream *activeStream, streamID string) {
	d.streamManager.mu.Lock()
	d.streamManager.streams[streamID] = stream
	d.streamManager.mu.Unlock()

	d.trackStream(stream.panelId, streamID, stream)
}

func (d *Datasource) cleanupStream(stream *activeStream) {
	d.streamManager.mu.Lock()
	delete(d.streamManager.streams, stream.streamID)

	if panelStreams, exists := d.streamManager.activeStreams[stream.panelId]; exists {
		delete(panelStreams, stream.streamID)
		if len(panelStreams) == 0 {
			delete(d.streamManager.activeStreams, stream.panelId)
		}
	}
	d.streamManager.mu.Unlock()

	d.logger.Debug("Stream closed", "streamID", stream.streamID)
}

func (d *Datasource) trackStream(panelId string, streamId string, stream *activeStream) {
	d.streamManager.mu.Lock()
	defer d.streamManager.mu.Unlock()

	if _, exists := d.streamManager.activeStreams[panelId]; !exists {
		d.streamManager.activeStreams[panelId] = make(map[string]*activeStream)
	}

	d.streamManager.streams[streamId] = stream
	d.streamManager.activeStreams[panelId][streamId] = stream

	d.logger.Debug("Stream tracked",
		"panelId", panelId,
		"streamId", streamId,
		"totalStreams", len(d.streamManager.streams),
		"panelStreams", len(d.streamManager.activeStreams[panelId]))
}

func (d *Datasource) getStreamsByPanel(panelId string) []*activeStream {
	d.streamManager.mu.RLock()
	defer d.streamManager.mu.RUnlock()

	result := make([]*activeStream, 0, 5)
	if panelStreams, exists := d.streamManager.activeStreams[panelId]; exists {
		for _, stream := range panelStreams {
			result = append(result, stream)
		}
	}
	return result
}
