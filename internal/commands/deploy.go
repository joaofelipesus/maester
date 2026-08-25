package commands

import (
	"fmt"
	"maester/internal"
	"strings"
)

func commandMergedWithCd(cfg internal.Config, cmd string) []string {
	mergedCommand := fmt.Sprintf("%s@%s cd %s && %s", cfg.ServerUserName, cfg.ServerIP, cfg.AppPath, cmd)
	return strings.Split(mergedCommand, " ")
}

// TODO: add coverage
func DeployNewVersion(cfg internal.Config, createCommand internal.CommandFactory) error {
	fmt.Println("Start deploy")
	fmt.Println("Stop container")

	// TODO: extract private function that merge and execute command, DeployNewVersion should only orchestrate
	//       the steps
	stopCommand := commandMergedWithCd(cfg, cfg.StopCommand)
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
	buildAndStartCommand := commandMergedWithCd(cfg, cfg.StartCommand)
	cmd = createCommand("ssh", buildAndStartCommand...)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	fmt.Println("Start app [SUCCESS]")

	return nil
}
