package plugin

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (d *Datasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	return d.health.CheckHealth(ctx, req)
}
