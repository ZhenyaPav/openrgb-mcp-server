package openrgb

import (
	"context"

	"github.com/csutorasa/go-openrgb-sdk"
)

func (c *Client) ListProfiles() ([]Profile, error) {
	return c.ListProfilesCtx(context.Background())
}

func (c *Client) ListProfilesCtx(ctx context.Context) ([]Profile, error) {
	err := c.c.RequestProtocolVersionCtx(ctx)
	if err != nil {
		return nil, err
	}
	rplRsp, err := c.c.RequestProfileListCtx(ctx)

	if err != nil {
		return nil, err
	}
	profiles := make([]Profile, len(rplRsp.Names))
	for i, name := range rplRsp.Names {
		profiles[i] = Profile{
			Name: name,
		}
	}
	return profiles, nil
}

// TODO: Fix that it doesn't return an error if the profile doesn't exist
func (c *Client) SetProfile(name string) error {
	return c.SetProfileCtx(context.Background(), name)
}

func (c *Client) SetProfileCtx(ctx context.Context, name string) error {
	err := c.c.RequestProtocolVersionCtx(ctx)
	if err != nil {
		return err
	}
	rsplReq := &openrgb.RequestLoadProfileRequest{
		ProfileName: name,
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return c.c.RequestLoadProfile(rsplReq)
}
