package plugin

import (
	"fmt"
	"strconv"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func (d *Datasource) handleManualQuery(qm queryModel, timeRange backend.TimeRange, frameBaseName string) backend.DataResponse {
	d.logger.Debug("Processing manual query",
		"method", qm.ManualMethod,
		"objectId", qm.ManualObjectId,
		"timeRange", fmt.Sprintf("%v to %v", timeRange.From, timeRange.To),
	)

	if qm.ManualMethod == "" {
		d.logger.Error("Manual method is required")
		d.metrics.IncError("missing_manual_method")
		return backend.ErrDataResponse(backend.StatusBadRequest, "manual method is required")
	}

	response, err := d.api.ExecuteManualMethod(qm.ManualMethod, qm.ManualObjectId)
	if err != nil {
		d.logger.Error("Manual query failed",
			"error", err,
			"method", qm.ManualMethod,
		)
		d.metrics.IncError("manual_query_failed")
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("API request failed: %v", err))
	}

	keys := make([]string, len(response.KeyValues))
	values := make([]string, len(response.KeyValues))

	for i, kv := range response.KeyValues {
		keys[i] = kv.Key
		switch v := kv.Value.(type) {
		case string:
			values[i] = v
		case float64:
			values[i] = strconv.FormatFloat(v, 'f', -1, 64)
		case bool:
			values[i] = strconv.FormatBool(v)
		case nil:
			values[i] = "null"
		default:
			values[i] = fmt.Sprintf("%v", v)
		}
	}

	frame := data.NewFrame(frameBaseName,
		data.NewField("Key", nil, keys).SetConfig(&data.FieldConfig{
			DisplayName: "Property",
		}),
		data.NewField("Value", nil, values).SetConfig(&data.FieldConfig{
			DisplayName: "Value",
		}),
	).SetMeta(&data.FrameMeta{
		Type:   data.FrameTypeTimeSeriesWide,
		Custom: response.Manuel,
	})

	return backend.DataResponse{
		Frames: []*data.Frame{frame},
	}
}
