package app

import (
	"context"
	"fmt"
	"time"

	"github.com/theankitbhardwaj/openrgb-mcp-server/internal/openrgb"
)

type Service struct {
	rgbClient *openrgb.Client
	timeout   time.Duration
}

func NewService(rgbClient *openrgb.Client, timeout time.Duration) *Service {
	return &Service{
		rgbClient: rgbClient,
		timeout:   timeout,
	}
}

func (s *Service) ListDevices(ctx context.Context) ([]openrgb.DeviceInfo, error) {
	ctx, cancel := s.withTimeout(ctx)
	if cancel != nil {
		defer cancel()
	}
	return s.rgbClient.ListDeviceInfosCtx(ctx)
}

func (s *Service) SetDeviceColor(ctx context.Context, deviceID int, r, g, b int) error {
	if err := validateRGB(r, g, b); err != nil {
		return err
	}

	ctx, cancel := s.withTimeout(ctx)
	if cancel != nil {
		defer cancel()
	}
	deviceInfo, err := s.rgbClient.GetDeviceInfoCtx(ctx, deviceID)
	if err != nil {
		return err
	}
	return s.rgbClient.SetDeviceColorCtx(ctx, *deviceInfo, r, g, b)
}

func (s *Service) SetAllDevicesColor(ctx context.Context, r, g, b int) error {
	if err := validateRGB(r, g, b); err != nil {
		return err
	}
	ctx, cancel := s.withTimeout(ctx)
	if cancel != nil {
		defer cancel()
	}
	return s.rgbClient.SetAllDeviceColorCtx(ctx, r, g, b)
}

func (s *Service) ListProfiles(ctx context.Context) ([]openrgb.Profile, error) {
	ctx, cancel := s.withTimeout(ctx)
	if cancel != nil {
		defer cancel()
	}
	profiles, err := s.rgbClient.ListProfilesCtx(ctx)
	if err != nil {
		return nil, err
	}
	return profiles, nil
}

func (s *Service) SetProfile(ctx context.Context, profileName string) error {
	if profileName == "" {
		return fmt.Errorf("profile name cannot be empty")
	}
	ctx, cancel := s.withTimeout(ctx)
	if cancel != nil {
		defer cancel()
	}
	return s.rgbClient.SetProfileCtx(ctx, profileName)
}

func (s *Service) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if s.timeout <= 0 {
		return ctx, nil
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, s.timeout)
	return timeoutCtx, cancel
}
