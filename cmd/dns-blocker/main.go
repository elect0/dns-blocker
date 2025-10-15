package main

import (
	"fmt"
	"github.com/elect0/dns-blocker/internal/config"
	"github.com/elect0/dns-blocker/internal/logging"
)

func main() {

	config, err := config.LoadConfig()
	if err != nil {
		return
	}

	logLevel, err := logging.string

	fmt.Printf("Configuration loaded successfully:\n%+v\n", config)
}
