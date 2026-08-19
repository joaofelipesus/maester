package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

// LATER: add the service external URL to ping and check if its alive while deploying it
type config struct {
	serverUserName string
	serverIP       string
	downloadLogs   bool
}

// NOTE: the local network IP address is used because the deployment is done using SSH commands which is
// disabled on ngrok.
func main() {
	serverUserName := flag.String("user", "", "The user that will be used to access the server")
	serverIP := flag.String("ip", "", "Ther server IP")
	downloadLogs := flag.Bool("doload-logs", false, "Download logs snapshot")
	flag.Parse()

	cfg := config{
		serverUserName: *serverUserName,
		serverIP:       *serverIP,
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
		return errors.New("IP address ir required")
	}

	return nil
}

// 1. ping server
// 2. check if SSH is available (TODO)
// 3. run command
func run(cfg config) {
	if err := pingServer(cfg, realCommand); err != nil {
		fmt.Printf("Failed to ping server, check if its up, and in the same network")
		os.Exit(1)
	}

	if err := checkSSHAvailable(cfg, realCommand); err != nil {
		fmt.Printf("Failed to ping server, check if its up, and in the same network")
		os.Exit(1)
	}
}

// TODO: move to a external module

// interface extracted to enable cover function pingServer and other functions that uses
// direct command calls.
type command interface {
	CombinedOutput() ([]byte, error)
}

type commandFactory func(name string, args ...string) command

func realCommand(name string, args ...string) command {
	return exec.Command(name, args...)
}

func pingServer(cfg config, createComand commandFactory) error {
	fmt.Printf("Start ping server on address %s\n", cfg.serverIP)

	cmd := createComand("ping", "-c", "1", cfg.serverIP)

	if _, err := cmd.CombinedOutput(); err != nil {
		return err
	}

	fmt.Println("Server ping [SUCCESS]")

	return nil
}

// The command nc (netcat) is used for scan ports
func checkSSHAvailable(cfg config, createCommand commandFactory) error {
	fmt.Printf("Start cherck SSH port on %s:22", cfg.serverIP)

	cmd := createCommand("nc", "-zv", cfg.serverIP, "22")

	if _, err := cmd.CombinedOutput(); err != nil {
		return err
	}

	fmt.Println("SHH check [SUCCESS]")

	return nil
}
