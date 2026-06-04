package plugin

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (d *Datasource) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	return d.resource.CallResource(ctx, req, sender)
}
