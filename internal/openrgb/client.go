package openrgb

import (
	"context"

	"github.com/csutorasa/go-openrgb-sdk"
)

type Client struct {
	c *openrgb.Client
}

func ConnectClient(host string, port int) (*Client, error) {
	cl, err := openrgb.NewClientHostPort(host, port)

	if err != nil {
		return nil, err
	}

	wrapped := &Client{cl}
	if err := wrapped.c.Initialize("openrgb-mcp-server"); err != nil {
		return nil, err
	}
	return wrapped, nil
}

func (c *Client) Close() error {
	return c.c.Close()
}

func (c *Client) RequestProtocolVersionCtx(ctx context.Context) error {
	return c.c.RequestProtocolVersionCtx(ctx)
}
