package server

import (
	"fmt"
	"maester/internal"
)

// TODO: move functions to internal dir

func PingServer(cfg internal.Config, createComand internal.CommandFactory) error {
	fmt.Printf("Start ping server on address %s\n", cfg.ServerIP)

	cmd := createComand("ping", "-c", "1", cfg.ServerIP)

	if _, err := cmd.CombinedOutput(); err != nil {
		return err
	}

	fmt.Println("Server ping [SUCCESS]")

	return nil
}

// The command nc (netcat) is used for scan ports
func CheckSSHAvailable(cfg internal.Config, createCommand internal.CommandFactory) error {
	fmt.Printf("Start check SSH port on %s:22", cfg.ServerIP)

	cmd := createCommand("nc", "-zv", cfg.ServerIP, "22")

	if _, err := cmd.CombinedOutput(); err != nil {
		return err
	}

	fmt.Println("SHH check [SUCCESS]")

	return nil
}
