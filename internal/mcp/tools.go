package mcp

import (
	"context"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theankitbhardwaj/openrgb-mcp-server/internal/app"
	"github.com/theankitbhardwaj/openrgb-mcp-server/internal/openrgb"
)

func RegisterTools(server *mcp.Server, svc *app.Service) {
	listDevices := &mcp.Tool{
		Name:        "list_devices",
		Description: "List all devices connected to OpenRGB",
	}
	setDeviceColor := &mcp.Tool{
		Name:        "set_device_color",
		Description: "Set the color of a specific device",
	}
	setAllColor := &mcp.Tool{
		Name:        "set_all_color",
		Description: "Set the color of all devices",
	}
	listProfiles := &mcp.Tool{
		Name:        "list_profiles",
		Description: "List all saved profiles on OpenRGB",
	}
	setProfile := &mcp.Tool{
		Name:        "set_profile",
		Description: "Set the saved profile",
	}
	mcp.AddTool(server, listDevices, handleListDevices(svc))
	mcp.AddTool(server, setDeviceColor, handleSetDeviceColor(svc))
	mcp.AddTool(server, setAllColor, handleSetAllColor(svc))
	mcp.AddTool(server, listProfiles, handleListProfiles(svc))
	mcp.AddTool(server, setProfile, handleSetProfile(svc))
}

func handleListDevices(svc *app.Service) mcp.ToolHandlerFor[ListDevicesParams, []openrgb.DeviceInfo] {
	return func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[ListDevicesParams]) (*mcp.CallToolResultFor[[]openrgb.DeviceInfo], error) {
		// Call the service to list devices
		devices, err := svc.ListDevices(ctx)
		if err != nil {
			log.Printf("list_devices failed: %v", err)
			return &mcp.CallToolResultFor[[]openrgb.DeviceInfo]{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error listing devices: %v", err)}},
				IsError: true,
			}, nil
		}

		return &mcp.CallToolResultFor[[]openrgb.DeviceInfo]{
			Content:           []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("List of devices: %v", devices)}},
			StructuredContent: devices,
		}, nil
	}
}

func handleSetDeviceColor(svc *app.Service) mcp.ToolHandlerFor[SetDeviceColorParams, string] {
	return func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[SetDeviceColorParams]) (*mcp.CallToolResultFor[string], error) {
		err := svc.SetDeviceColor(ctx, params.Arguments.DeviceID, params.Arguments.R, params.Arguments.G, params.Arguments.B)
		if err != nil {
			log.Printf("set_device_color failed: %v", err)
			return &mcp.CallToolResultFor[string]{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error setting device color: %v", err)}},
				IsError: true,
			}, nil
		}

		return &mcp.CallToolResultFor[string]{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Successfully set color for device: %v", params.Arguments.DeviceID)}},
		}, nil
	}
}

func handleSetAllColor(svc *app.Service) mcp.ToolHandlerFor[SetAllColorParams, string] {
	return func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[SetAllColorParams]) (*mcp.CallToolResultFor[string], error) {
		err := svc.SetAllDevicesColor(ctx, params.Arguments.R, params.Arguments.G, params.Arguments.B)
		if err != nil {
			log.Printf("set_all_color failed: %v", err)
			return &mcp.CallToolResultFor[string]{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error setting all device colors: %v", err)}},
				IsError: true,
			}, nil
		}

		return &mcp.CallToolResultFor[string]{
			Content: []mcp.Content{&mcp.TextContent{Text: "Successfully set color for devices"}},
		}, nil
	}
}

func handleListProfiles(svc *app.Service) mcp.ToolHandlerFor[ListProfilesParams, []openrgb.Profile] {
	return func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[ListProfilesParams]) (*mcp.CallToolResultFor[[]openrgb.Profile], error) {
		profiles, err := svc.ListProfiles(ctx)
		if err != nil {
			log.Printf("list_profiles failed: %v", err)
			return &mcp.CallToolResultFor[[]openrgb.Profile]{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error listing profiles: %v", err)}},
				IsError: true,
			}, nil
		}

		return &mcp.CallToolResultFor[[]openrgb.Profile]{
			Content:           []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("List of devices: %v", profiles)}},
			StructuredContent: profiles,
		}, nil
	}
}

func handleSetProfile(svc *app.Service) mcp.ToolHandlerFor[SetProfileParams, string] {
	return func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[SetProfileParams]) (*mcp.CallToolResultFor[string], error) {
		err := svc.SetProfile(ctx, params.Arguments.ProfileName)
		if err != nil {
			log.Printf("set_profile failed: %v", err)
			return &mcp.CallToolResultFor[string]{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error setting profile: %v", err)}},
				IsError: true,
			}, nil
		}

		return &mcp.CallToolResultFor[string]{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Successfully set profile: %v", params.Arguments.ProfileName)}},
		}, nil
	}
}
