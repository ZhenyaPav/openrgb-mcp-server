package openrgb

import (
	"errors"
	"net"
	"strings"
	"sync"
	"syscall"

	"github.com/csutorasa/go-openrgb-sdk"
)

type Client struct {
	c    *openrgb.Client
	host string
	port int
	mu   sync.RWMutex
}

func ConnectClient(host string, port int) (*Client, error) {
	wrapped := &Client{
		host: host,
		port: port,
	}
	if err := wrapped.reconnect(); err != nil {
		return nil, err
	}
	return wrapped, nil
}

func (c *Client) Close() error {
	return c.client().Close()
}

func (c *Client) client() *openrgb.Client {
	c.mu.RLock()
	cl := c.c
	c.mu.RUnlock()
	return cl
}

func (c *Client) reconnect() error {
	cl, err := openrgb.NewClientHostPort(c.host, c.port)
	if err != nil {
		return err
	}
	c.mu.Lock()
	old := c.c
	c.c = cl
	c.mu.Unlock()
	if old != nil {
		_ = old.Close()
	}
	return nil
}

func (c *Client) withRetryErr(fn func(*openrgb.Client) error) error {
	_, err := c.withRetryValue(func(cl *openrgb.Client) (any, error) {
		return nil, fn(cl)
	})
	return err
}

func (c *Client) withRetryValue(fn func(*openrgb.Client) (any, error)) (any, error) {
	cl := c.client()
	if cl == nil {
		if err := c.reconnect(); err != nil {
			return nil, err
		}
		cl = c.client()
	}
	val, err := fn(cl)
	if err == nil {
		return val, nil
	}
	if !isConnError(err) {
		return nil, err
	}
	if err := c.reconnect(); err != nil {
		return nil, err
	}
	return fn(c.client())
}

func isConnError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, net.ErrClosed) {
		return true
	}
	if errors.Is(err, syscall.EPIPE) || errors.Is(err, syscall.ECONNRESET) || errors.Is(err, syscall.ECONNABORTED) {
		return true
	}
	if strings.Contains(strings.ToLower(err.Error()), "broken pipe") {
		return true
	}
	if opErr := new(net.OpError); errors.As(err, &opErr) {
		if errors.Is(opErr.Err, syscall.EPIPE) || errors.Is(opErr.Err, syscall.ECONNRESET) {
			return true
		}
	}
	return false
}
