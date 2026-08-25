package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"maester/internal/logs"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

// LATER: add the service external URL to ping and check if its alive while deploying it
type config struct {
	serverUserName string
	serverIP       string
	appPath        string
	stopCommand    string
	startCommand   string
	downloadLogs   bool
	deploy         bool
}

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

	cfg := config{
		serverUserName: yamlConfigs["serverUserName"],
		serverIP:       yamlConfigs["serverIP"],
		appPath:        yamlConfigs["appPath"],
		stopCommand:    yamlConfigs["stopCommand"],
		startCommand:   yamlConfigs["startCommand"],
		downloadLogs:   *downloadLogs,
		deploy:         *deploy,
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

	run(cfg, realCommand, outputFile)
}

// validates if any required tag is missing
func validateRequiredConfigs(cfg config) error {
	fmt.Println(cfg)

	if cfg.serverUserName == "" {
		return errors.New("user is required")
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
func run(cfg config, createCommand commandFactory, outputFile io.Writer) {
	if err := PingServer(cfg, createCommand); err != nil {
		fmt.Printf("Failed to ping server, check if its up, and in the same network")
		os.Exit(1)
	}

	if err := CheckSSHAvailable(cfg, createCommand); err != nil {
		fmt.Printf("Failed to check SSH server")
		os.Exit(1)
	}

	if cfg.downloadLogs {
		if err := DownloadLogs(cfg, createCommand, outputFile); err != nil {
			fmt.Printf("Failed to download logs")
			os.Exit(1)
		}
	}

	if cfg.deploy {
		if err := DeployNewVersion(cfg, createCommand); err != nil {
			fmt.Printf("Failed to deploy\n")
			fmt.Println(err)
			os.Exit(1)
		}
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

// TODO: move functions to a module with the commands implementations
// TODO: add coverage
func DownloadLogs(cfg config, createCommand commandFactory, output io.Writer) error {
	fmt.Println("Start fetching logs")

	serverUserAndIP := fmt.Sprintf("%s@%s", cfg.serverUserName, cfg.serverIP)
	dockerCommand := fmt.Sprintf("%s cd %s && docker compose -f docker-compose.yml -f docker-compose.prod.yml logs --since 2h app", serverUserAndIP, cfg.appPath)
	splittedDockerCommand := strings.Split(dockerCommand, " ")
	cmd := createCommand("ssh", splittedDockerCommand...)

	commandOutput, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Fetching logs [SUCCESS]")
	fmt.Println("Writing to a file")

	var buffer bytes.Buffer
	buffer.Write(commandOutput)
	logs.RemoveUpCalls(&buffer, output)
	fmt.Println("Writing to a file [SUCCESS]")

	return nil
}

func commandMergedWithCd(cfg config, cmd string) []string {
	mergedCommand := fmt.Sprintf("%s@%s cd %s && %s", cfg.serverUserName, cfg.serverIP, cfg.appPath, cmd)
	return strings.Split(mergedCommand, " ")
}

// TODO: move functions to a module with the commands implementations
// TODO: add coverage
func DeployNewVersion(cfg config, createCommand commandFactory) error {
	fmt.Println("Start deploy")
	fmt.Println("Stop container")

	// TODO: extract private function that merge and execute command, DeployNewVersion should only orchestrate
	//       the steps
	stopCommand := commandMergedWithCd(cfg, cfg.stopCommand)
	cmd := createCommand("ssh", stopCommand...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	fmt.Println("Stop container [SUCCESS]")

	gitPullCommand := commandMergedWithCd(cfg, "git pull")
	cmd = createCommand("ssh", gitPullCommand...)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	fmt.Println("Update repository [SUCCESS]")

	fmt.Println("Start app")
	buildAndStartCommand := commandMergedWithCd(cfg, cfg.startCommand)
	cmd = createCommand("ssh", buildAndStartCommand...)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	fmt.Println("Start app [SUCCESS]")

	return nil
}
