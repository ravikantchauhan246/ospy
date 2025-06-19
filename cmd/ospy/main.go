package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/ravikantchauhan246/ospy/internal/config"
	"github.com/ravikantchauhan246/ospy/internal/monitor"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	fmt.Printf("Loaded configuration:\n")
	fmt.Printf("- Monitoring interval: %v\n", cfg.Monitoring.Interval)
	fmt.Printf("- Monitoring timeout: %v\n", cfg.Monitoring.Timeout)
	fmt.Printf("- Workers: %d\n", cfg.Monitoring.Workers)
	fmt.Printf("- Websites: %d\n", len(cfg.Websites))

	// Test checker with configured websites
	checker := monitor.NewChecker(cfg.Monitoring.Timeout)
	
	for _, website := range cfg.Websites {
		result := checker.Check(context.Background(), website.URL)
		
		status := "✅"
		if !result.IsUp {
			status = "❌"
		}
		
		fmt.Printf("%s %s (%s) - Status: %d, Time: %v\n", 
			status, website.Name, website.URL, result.Status, result.ResponseTime)
	}
}