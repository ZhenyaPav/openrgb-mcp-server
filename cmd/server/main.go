package main

import (
	"context"
	"log"
	"os"

	"github.com/theankitbhardwaj/openrgb-mcp-server/internal/app"
	"github.com/theankitbhardwaj/openrgb-mcp-server/internal/mcp"
	"github.com/theankitbhardwaj/openrgb-mcp-server/internal/openrgb"
	"github.com/theankitbhardwaj/openrgb-mcp-server/pkg/util"
)

func main() {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags)

	cfg, err := util.LoadConfig("config/config.yaml")
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		os.Exit(1)
	}

	client, err := openrgb.ConnectClient(cfg.OpenRGB.Host, cfg.OpenRGB.Port)
	if err != nil {
		log.Printf("Failed to connect to OpenRGB server: %v", err)
		os.Exit(1)
	}

	defer client.Close()

	svc := app.NewService(client)

	mcpServer := mcp.NewServer(cfg.Server.Name, cfg.Server.Version)

	mcp.RegisterTools(mcpServer, svc)

	if err := mcp.RunStdio(context.Background(), mcpServer); err != nil {
		log.Printf("Server runtime error: %v", err)
	}
}
