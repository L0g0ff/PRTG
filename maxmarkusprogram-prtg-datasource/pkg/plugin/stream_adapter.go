package plugin

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (d *Datasource) SubscribeStream(ctx context.Context, req *backend.SubscribeStreamRequest) (*backend.SubscribeStreamResponse, error) {
	return d.stream.SubscribeStream(ctx, req)
}

func (d *Datasource) PublishStream(ctx context.Context, req *backend.PublishStreamRequest) (*backend.PublishStreamResponse, error) {
	return d.stream.PublishStream(ctx, req)
}

func (d *Datasource) RunStream(ctx context.Context, req *backend.RunStreamRequest, sender *backend.StreamSender) error {
	return d.stream.RunStream(ctx, req, sender)
}
