package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"maester/internal"
	"maester/internal/commands"
	"maester/internal/server"
	"os"

	"gopkg.in/yaml.v3"
)

// TODO: add README.md
// NOTE: the local network IP address is used because the deployment is done using SSH commands which is
// disabled on ngrok.
func main() {
	downloadLogs := flag.Bool("download-logs", false, "Download logs snapshot")
	deploy := flag.Bool("deploy", false, "Deploy a new version")
	flag.Parse()

	// TODO: document
	configsFile, err := os.ReadFile("configs.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var yamlConfigs map[string]string

	err = yaml.Unmarshal(configsFile, &yamlConfigs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Print(yamlConfigs)

	cfg := internal.Config{
		ServerUserName: yamlConfigs["serverUserName"],
		ServerIP:       yamlConfigs["serverIP"],
		AppPath:        yamlConfigs["appPath"],
		StopCommand:    yamlConfigs["stopCommand"],
		StartCommand:   yamlConfigs["startCommand"],
		DownloadLogs:   *downloadLogs,
		Deploy:         *deploy,
	}

	if err := validateRequiredConfigs(cfg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// TODO: add conditional when support new commands
	outputFile, err := os.Create("logs/logs.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer outputFile.Close()

	run(cfg, internal.RealCommand, outputFile)
}

// validates if any required tag is missing
func validateRequiredConfigs(cfg internal.Config) error {
	fmt.Println(cfg)

	if cfg.ServerUserName == "" {
		return errors.New("user is required")
	}

	if cfg.ServerIP == "" {
		return errors.New("IP address is required")
	}

	if cfg.AppPath == "" {
		return errors.New("App path is required")
	}

	return nil
}

// 1. ping server
// 2. check if SSH is available (TODO)
// 3. run command
func run(cfg internal.Config, createCommand internal.CommandFactory, outputFile io.Writer) {
	if err := server.PingServer(cfg, createCommand); err != nil {
		fmt.Printf("Failed to ping server, check if its up, and in the same network")
		os.Exit(1)
	}

	if err := server.CheckSSHAvailable(cfg, createCommand); err != nil {
		fmt.Printf("Failed to check SSH server")
		os.Exit(1)
	}

	if cfg.DownloadLogs {
		if err := commands.DownloadLogs(cfg, createCommand, outputFile); err != nil {
			fmt.Printf("Failed to download logs")
			os.Exit(1)
		}
	}

	if cfg.Deploy {
		if err := commands.DeployNewVersion(cfg, createCommand); err != nil {
			fmt.Printf("Failed to deploy\n")
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
