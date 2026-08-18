package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type config struct {
	serverUserName string
	serverIP       string
	downloadLogs   bool
}

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

func run(cng config) {
	fmt.Print("ping...")
}
