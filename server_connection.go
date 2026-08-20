package main

import "fmt"

// TODO: move functions to internal dir

func PingServer(cfg config, createComand commandFactory) error {
	fmt.Printf("Start ping server on address %s\n", cfg.serverIP)

	cmd := createComand("ping", "-c", "1", cfg.serverIP)

	if _, err := cmd.CombinedOutput(); err != nil {
		return err
	}

	fmt.Println("Server ping [SUCCESS]")

	return nil
}

// The command nc (netcat) is used for scan ports
func CheckSSHAvailable(cfg config, createCommand commandFactory) error {
	fmt.Printf("Start cherck SSH port on %s:22", cfg.serverIP)

	cmd := createCommand("nc", "-zv", cfg.serverIP, "22")

	if _, err := cmd.CombinedOutput(); err != nil {
		return err
	}

	fmt.Println("SHH check [SUCCESS]")

	return nil
}
