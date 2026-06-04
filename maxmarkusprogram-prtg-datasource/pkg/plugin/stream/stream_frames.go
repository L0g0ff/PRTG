package stream

import (
	"fmt"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func extractChannelName(frame *data.Frame) string {
	if frame.Meta != nil && frame.Meta.Custom != nil {
		if metaMap, ok := frame.Meta.Custom.(map[string]interface{}); ok {
			if ch, exists := metaMap["channel"]; exists {
				return fmt.Sprint(ch)
			}
		}
	}

	parts := strings.Split(frame.Name, "_")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}

	return ""
}

func extractFrameData(frame *data.Frame) ([]time.Time, []float64) {
	if len(frame.Fields) < 2 || frame.Fields[0].Len() == 0 {
		return nil, nil
	}

	timeField := frame.Fields[0]
	valueField := frame.Fields[1]

	times := make([]time.Time, timeField.Len())
	values := make([]float64, valueField.Len())

	for i := 0; i < timeField.Len(); i++ {
		if t, ok := timeField.At(i).(time.Time); ok {
			times[i] = t
		}
		if v, ok := valueField.At(i).(float64); ok {
			values[i] = v
		}
	}

	return times, values
}

func createStreamingFrame(stream *activeStream, channelName string, state *channelState, from, to time.Time) *data.Frame {
	displayName := buildDisplayName(stream, channelName)

	frameName := fmt.Sprintf("stream_%s_%s", stream.sensorId, channelName)
	frame := data.NewFrame(frameName,
		data.NewField("Time", nil, state.buffer.times),
		data.NewField("Value", nil, state.buffer.values).SetConfig(&data.FieldConfig{
			DisplayName: displayName,
		}),
	)

	now := time.Now().UnixMilli()
	streamingStatus := map[string]interface{}{
		"active":      true,
		"lastUpdate":  now,
		"lastValue":   state.lastValue,
		"dataPoints":  len(state.buffer.times),
		"streamId":    stream.streamID,
		"sensorId":    stream.sensorId,
		"channelName": channelName,
		"isLive":      true,
		"state":       "streaming",
	}

	frame.Meta = &data.FrameMeta{
		Type: data.FrameTypeTimeSeriesMulti,
		Custom: map[string]interface{}{
			"from":           from.UnixMilli(),
			"to":             to.UnixMilli(),
			"channel":        channelName,
			"updating":       true,
			"streaming":      true,
			"live":           true,
			"streaming_rate": stream.interval.Milliseconds(),
			"isActive":       true,
			"stable":         true,
			"timezone":       "UTC",
			"state":          "streaming",
			"streamStatus":   streamingStatus,
		},
	}

	return frame
}

func buildDisplayName(stream *activeStream, channelName string) string {
	displayName := channelName

	if stream.includeGroupName && stream.group != "" {
		displayName = fmt.Sprintf("%s - %s", stream.group, displayName)
	}
	if stream.includeDeviceName && stream.device != "" {
		displayName = fmt.Sprintf("%s - %s", stream.device, displayName)
	}
	if stream.includeSensorName && stream.sensor != "" {
		displayName = fmt.Sprintf("%s - %s", stream.sensor, displayName)
	}

	return displayName
}
