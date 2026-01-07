package openrgb

import "github.com/csutorasa/go-openrgb-sdk"

func (c *Client) ListProfiles() ([]Profile, error) {
	val, err := c.withRetryValue(func(cl *openrgb.Client) (any, error) {
		err := cl.RequestProtocolVersion()
		if err != nil {
			return nil, err
		}
		rplRsp, err := cl.RequestProfileList()

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
	})
	if err != nil {
		return nil, err
	}
	return val.([]Profile), nil
}

// TODO: Fix that it doesn't return an error if the profile doesn't exist
func (c *Client) SetProfile(name string) error {
	return c.withRetryErr(func(cl *openrgb.Client) error {
		err := cl.RequestProtocolVersion()
		if err != nil {
			return err
		}
		rsplReq := &openrgb.RequestLoadProfileRequest{
			ProfileName: name,
		}
		return cl.RequestLoadProfile(rsplReq)
	})
}
