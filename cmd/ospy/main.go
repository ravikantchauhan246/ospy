package main

import (
	
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	

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

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Create checker and worker pool
	checker := monitor.NewChecker(cfg.Monitoring.Timeout)
	workerPool := monitor.NewWorkerPool(cfg.Monitoring.Workers, checker)

	// Convert config websites to monitor websites
	websites := make([]monitor.Website, len(cfg.Websites))
	for i, w := range cfg.Websites {
		websites[i] = monitor.Website{
			Name:           w.Name,
			URL:            w.URL,
			Method:         w.Method,
			Headers:        w.Headers,
			ExpectedStatus: w.ExpectedStatus,
			CheckContent:   w.CheckContent,
			Timeout:        w.Timeout,
		}
	}

	// Start worker pool
	workerPool.Start()

	// Create and start scheduler
	scheduler := monitor.NewScheduler(workerPool, websites, cfg.Monitoring.Interval)
	scheduler.Start()

	// Handle results
	go func() {
		for result := range workerPool.Results() {
			status := "✅"
			if !result.IsUp {
				status = "❌"
			}
			
			log.Printf("%s %s (%s) - %s (Time: %v)", 
				status, result.WebsiteName, result.URL, result.Message, result.ResponseTime)
			
			if result.Error != nil {
				log.Printf("   Error: %v", result.Error)
			}
		}
	}()

	log.Printf("Ospy started - monitoring %d websites every %v", 
		len(websites), cfg.Monitoring.Interval)
	log.Printf("Press Ctrl+C to stop")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	scheduler.Stop()
	workerPool.Close()
	log.Println("Shutdown complete")
}