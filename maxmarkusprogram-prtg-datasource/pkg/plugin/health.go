package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (d *Datasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	d.ClearAllCaches()

	_, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	d.logger.Debug("Starting health check")

	status, err := d.api.GetStatusList()
	if err != nil {
		d.logger.Error("PRTG health check failed", "error", err)
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: fmt.Sprintf("PRTG API error: %s", err.Error()),
		}, nil
	}

	config, err := models.LoadPluginSettings(*req.PluginContext.DataSourceInstanceSettings)
	timezone := ""
	if err == nil {
		timezone = config.Timezone
	}

	details := map[string]interface{}{
		"version":      status.Version,
		"totalSensors": status.TotalSens,
	}
	if timezone != "" {
		details["timezone"] = timezone
	}

	message := fmt.Sprintf("Data source is working. PRTG Version: %s", status.Version)
	if timezone != "" {
		message = fmt.Sprintf("Data source is working. PRTG Version: %s | Timezone: %s", status.Version, timezone)
	}

	detailsJSON, _ := json.Marshal(details)

	return &backend.CheckHealthResult{
		Status:      backend.HealthStatusOk,
		Message:     message,
		JSONDetails: detailsJSON,
	}, nil
}
