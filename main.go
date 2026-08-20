package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// LATER: add the service external URL to ping and check if its alive while deploying it
type config struct {
	serverUserName string
	serverIP       string
	appPath        string
	downloadLogs   bool
}

// NOTE: the local network IP address is used because the deployment is done using SSH commands which is
// disabled on ngrok.
func main() {
	serverUserName := flag.String("user", "", "The user that will be used to access the server")
	serverIP := flag.String("ip", "", "The server IP")
	appPath := flag.String("app-path", "", "The absolute path of the app in the server")
	downloadLogs := flag.Bool("doload-logs", false, "Download logs snapshot")
	flag.Parse()

	cfg := config{
		serverUserName: *serverUserName,
		serverIP:       *serverIP,
		appPath:        *appPath,
		downloadLogs:   *downloadLogs,
	}

	if err := validateRequiredConfigs(cfg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	run(cfg)
}

// validates if any required tag is missing
func validateRequiredConfigs(cfg config) error {
	if cfg.serverUserName == "" {
		return errors.New("user tag is required")
	}

	if cfg.serverIP == "" {
		return errors.New("IP address is required")
	}

	if cfg.appPath == "" {
		return errors.New("App path is required")
	}

	return nil
}

// 1. ping server
// 2. check if SSH is available (TODO)
// 3. run command
func run(cfg config) {
	if err := PingServer(cfg, realCommand); err != nil {
		fmt.Printf("Failed to ping server, check if its up, and in the same network")
		os.Exit(1)
	}

	if err := CheckSSHAvailable(cfg, realCommand); err != nil {
		fmt.Printf("Failed to check SSH server")
		os.Exit(1)
	}

	if err := DownloadLogs(cfg, realCommand); err != nil {
		fmt.Printf("Failed to download logs")
		os.Exit(1)
	}
}

// interface extracted to enable cover function pingServer and other functions that uses
// direct command calls.
type command interface {
	CombinedOutput() ([]byte, error)
}

type commandFactory func(name string, args ...string) command

func realCommand(name string, args ...string) command {
	return exec.Command(name, args...)
}

// TODO: move functions to a module with the commands imple mentations
func DownloadLogs(cfg config, createCommand commandFactory) error {
	fmt.Println("Start fetching logs")

	serverUserAndIP := fmt.Sprintf("%s@%s", cfg.serverUserName, cfg.serverIP)
	dockerCommand := fmt.Sprintf("%s cd %s && docker compose -f docker-compose.yml -f docker-compose.prod.yml logs --since 2h app", serverUserAndIP, cfg.appPath)
	splittedDockerCommand := strings.Split(dockerCommand, " ")
	cmd := createCommand("ssh", splittedDockerCommand...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Fetching logs [SUCCESS]")
	fmt.Println("Writing to a file")

	file, err := os.Create("logs.txt")
	if err != nil {
		fmt.Println("Error while creating the log file")
		return err
	}
	defer file.Close()

	file.WriteString(string(output))
	fmt.Println("Writing to a file [SUCCESS]")

	return nil
}
