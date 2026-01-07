package openrgb

import "github.com/csutorasa/go-openrgb-sdk"

func (c *Client) SetDeviceColor(device DeviceInfo, r, g, b int) error {
	return c.withRetryErr(func(cl *openrgb.Client) error {
		col := openrgb.Color{R: uint8(r), G: uint8(g), B: uint8(b)}

		ulreq := openrgb.RGBControllerUpdateLedsRequest{
			LedColor: make([]openrgb.Color, device.LEDCount),
		}
		for i := range ulreq.LedColor {
			ulreq.LedColor[i] = col
		}

		return cl.RGBControllerUpdateLeds(uint32(device.ID), &ulreq)
	})
}

func (c *Client) SetAllDeviceColor(r, g, b int) error {
	devices, err := c.ListDeviceInfos()

	if err != nil {
		return err
	}

	col := openrgb.Color{R: uint8(r), G: uint8(g), B: uint8(b)}

	for _, device := range devices {
		if device.LEDCount == 0 {
			continue
		}
		ulreq := openrgb.RGBControllerUpdateLedsRequest{
			LedColor: make([]openrgb.Color, device.LEDCount),
		}
		for i := range ulreq.LedColor {
			ulreq.LedColor[i] = col
		}
		if err := c.withRetryErr(func(cl *openrgb.Client) error {
			return cl.RGBControllerUpdateLeds(uint32(device.ID), &ulreq)
		}); err != nil {
			return err
		}
	}

	return nil
}
