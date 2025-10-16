package main

import (
	"fmt"
	"log"
	"os"

	"github.com/elect0/dns-blocker/internal/config"
	"github.com/elect0/dns-blocker/internal/dns"
	logging "github.com/elect0/dns-blocker/internal/logging"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}

	logLevel, err := logging.StringToLevel(config.Logging.Level)
	if err != nil {
		log.Fatalf("invalid log level: %v", err)
	}

	logger := logging.New(os.Stdout, logLevel)

	logger.Info("logger intialized successfully", nil)

	logger.Info("configuration loaded successfully", map[string]string{
		"listen_address":       config.ListenAddress,
		"upstream_server":      config.UpstreamServer,
		"custom_records_count": fmt.Sprintf("%d", len(config.CustomRecords)),
	})

	blocklist, err := dns.LoadBlocklist(config.BlocklistPath)
	if err != nil {
	logger.Fatal("failed to load blocklist", err, nil)
	}

	logger.Info("blocklist loaded successfully", map[string]string{
		"domains_loaded": fmt.Sprintf("%d", len(blocklist)),
	})

}
