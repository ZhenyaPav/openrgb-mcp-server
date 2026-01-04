package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/theankitbhardwaj/openrgb-mcp-server/internal/app"
	"github.com/theankitbhardwaj/openrgb-mcp-server/internal/mcp"
	"github.com/theankitbhardwaj/openrgb-mcp-server/internal/openrgb"
	"github.com/theankitbhardwaj/openrgb-mcp-server/pkg/util"
)

func main() {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags)

	if wd, err := os.Getwd(); err == nil {
		log.Printf("cwd=%s", wd)
	} else {
		log.Printf("cwd unavailable: %v", err)
	}

	cfgPath, err := resolveConfigPath()
	if err != nil {
		log.Printf("Failed to resolve config path: %v", err)
		os.Exit(1)
	}
	cfg, err := util.LoadConfig(cfgPath)
	if err != nil {
		log.Printf("Failed to load config (%s): %v", cfgPath, err)
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

func resolveConfigPath() (string, error) {
	if envPath := os.Getenv("OPENRGB_MCP_CONFIG"); envPath != "" {
		return envPath, nil
	}
	cwdPath := "config/config.yaml"
	if _, err := os.Stat(cwdPath); err == nil {
		return cwdPath, nil
	}
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	exeDir := filepath.Dir(exe)
	exePath := filepath.Join(exeDir, "config", "config.yaml")
	if _, err := os.Stat(exePath); err == nil {
		return exePath, nil
	}
	exeParentPath := filepath.Join(filepath.Dir(exeDir), "config", "config.yaml")
	if _, err := os.Stat(exeParentPath); err == nil {
		return exeParentPath, nil
	}
	return cwdPath, fmt.Errorf("config not found at %s, %s, or %s", cwdPath, exePath, exeParentPath)
}
