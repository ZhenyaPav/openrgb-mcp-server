package openrgb

import (
	"context"

	"github.com/csutorasa/go-openrgb-sdk"
)

func (c *Client) SetDeviceColor(device DeviceInfo, r, g, b int) error {
	return c.SetDeviceColorCtx(context.Background(), device, r, g, b)
}

func (c *Client) SetDeviceColorCtx(ctx context.Context, device DeviceInfo, r, g, b int) error {
	col := openrgb.Color{R: uint8(r), G: uint8(g), B: uint8(b)}

	ulreq := openrgb.RGBControllerUpdateLedsRequest{
		LedColor: make([]openrgb.Color, device.LEDCount),
	}
	for i := range ulreq.LedColor {
		ulreq.LedColor[i] = col
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return c.c.RGBControllerUpdateLeds(uint32(device.ID), &ulreq)
}

func (c *Client) SetAllDeviceColor(r, g, b int) error {
	return c.SetAllDeviceColorCtx(context.Background(), r, g, b)
}

func (c *Client) SetAllDeviceColorCtx(ctx context.Context, r, g, b int) error {
	devices, err := c.ListDeviceInfosCtx(ctx)

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
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err := c.c.RGBControllerUpdateLeds(uint32(device.ID), &ulreq); err != nil {
			return err
		}
	}

	return nil
}
