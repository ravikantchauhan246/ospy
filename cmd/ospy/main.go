package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/ravikantchauhan246/ospy/internal/config"
	"github.com/ravikantchauhan246/ospy/internal/monitor"
	"github.com/ravikantchauhan246/ospy/internal/storage"
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

	// Ensure data directory exists
	if err := os.MkdirAll(filepath.Dir(cfg.Storage.Path), 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Initialize storage
	storage, err := storage.NewSQLiteStorage(cfg.Storage.Path)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer storage.Close()

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

	// Create monitor
	mon := monitor.NewMonitor(workerPool, websites, cfg.Monitoring.Interval, storage)

	// Connect worker pool results to monitor
	go func() {
		for result := range workerPool.Results() {
			mon.GetResults() <- result
		}
	}()

	// Start monitoring
	mon.Start()

	// Cleanup routine
	go func() {
		ticker := time.NewTicker(24 * time.Hour) // Daily cleanup
		defer ticker.Stop()
		
		for range ticker.C {
			if err := storage.Cleanup(cfg.Storage.RetentionDays); err != nil {
				log.Printf("Cleanup failed: %v", err)
			}
		}
	}()

	log.Printf("Ospy started - monitoring %d websites every %v", 
		len(websites), cfg.Monitoring.Interval)
	log.Printf("Data stored in: %s", cfg.Storage.Path)
	log.Printf("Press Ctrl+C to stop")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	mon.Stop()
	workerPool.Close()
	log.Println("Shutdown complete")
}