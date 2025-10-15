package main

import (
	"fmt"
	"github.com/elect0/dns-blocker/internal/config"
	"log"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error while loading config: %v", err)
	}

	fmt.Printf("Configuration loaded successfully:\n%+v\n", config)
}
