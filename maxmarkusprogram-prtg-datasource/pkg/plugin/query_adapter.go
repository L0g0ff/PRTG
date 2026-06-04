package plugin

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return d.query.QueryData(ctx, req)
}
