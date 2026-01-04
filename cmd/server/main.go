package main

import (
	"context"
	"fmt"
	"time"

	"github.com/theankitbhardwaj/openrgb-mcp-server/internal/app"
	"github.com/theankitbhardwaj/openrgb-mcp-server/internal/mcp"
	"github.com/theankitbhardwaj/openrgb-mcp-server/internal/openrgb"
	"github.com/theankitbhardwaj/openrgb-mcp-server/pkg/util"
)

func main() {
	fmt.Println("OpenRGB MCP server")

	cfg, err := util.LoadConfig("config/config.yaml")

	if err != nil {
		fmt.Println(err)
	}

	client, err := openrgb.ConnectClient(cfg.OpenRGB.Host, cfg.OpenRGB.Port) // client, err

	if err != nil {
		fmt.Printf("Failed to connect to OpenRGB server: %v\n", err)
		return
	}

	defer client.Close()

	svc := app.NewService(client, time.Duration(cfg.OpenRGB.TimeoutMS)*time.Millisecond)

	mcpServer := mcp.NewServer(cfg.Server.Name, cfg.Server.Version)

	mcp.RegisterTools(mcpServer, svc)

	if err := mcp.RunStdio(context.Background(), mcpServer); err != nil {
		fmt.Printf("server runtime error: %v", err)
	}
}
