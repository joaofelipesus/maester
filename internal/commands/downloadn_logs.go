package commands

import (
	"bytes"
	"fmt"
	"io"
	"maester/internal"
	"maester/internal/logs"
	"strings"
)

// TODO: move functions to a module with the commands implementations
// TODO: add coverage
func DownloadLogs(cfg internal.Config, createCommand internal.CommandFactory, output io.Writer) error {
	fmt.Println("Start fetching logs")

	serverUserAndIP := fmt.Sprintf("%s@%s", cfg.ServerUserName, cfg.ServerIP)
	dockerCommand := fmt.Sprintf("%s cd %s && docker compose -f docker-compose.yml -f docker-compose.prod.yml logs --since 2h app", serverUserAndIP, cfg.AppPath)
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
