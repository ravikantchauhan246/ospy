package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/ravikantchauhan246/ospy/internal/config"
	"github.com/ravikantchauhan246/ospy/internal/monitor"
	"github.com/ravikantchauhan246/ospy/internal/notifier"
	"github.com/ravikantchauhan246/ospy/internal/storage"
	"github.com/ravikantchauhan246/ospy/internal/web"
)

// Build-time variables
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func showHelp() {
	fmt.Printf("Ospy - Website Availability Monitor %s\n\n", Version)
	fmt.Println("USAGE:")
	fmt.Printf("  %s [OPTIONS]\n\n", os.Args[0])
	fmt.Println("OPTIONS:")
	fmt.Println("  -config string")
	fmt.Println("        Path to configuration file (default \"configs/config.yaml\")")
	fmt.Println("  -version")
	fmt.Println("        Show version information")
	fmt.Println("  -help")
	fmt.Println("        Show this help information")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Printf("  %s -config my-config.yaml    # Run with custom config\n", os.Args[0])
	fmt.Printf("  %s -version                  # Show version\n", os.Args[0])
	fmt.Println()
	fmt.Println("CONFIGURATION:")
	fmt.Println("  Create a YAML config file with your websites to monitor.")
	fmt.Println("  See https://github.com/ravikantchauhan246/ospy for examples.")
	fmt.Println()
	fmt.Println("ENVIRONMENT VARIABLES:")
	fmt.Println("  SMTP_USERNAME     - Email username for notifications")
	fmt.Println("  SMTP_PASSWORD     - Email password for notifications")
	fmt.Println("  TELEGRAM_BOT_TOKEN - Telegram bot token for notifications")
}

func main() {
	configPath := flag.String("config", "configs/config.yaml", "Path to configuration file")
	version := flag.Bool("version", false, "Show version information")
	help := flag.Bool("help", false, "Show help information")
	flag.Parse()

	// Show help information if requested
	if *help {
		showHelp()
		return
	}

	// Show version information if requested
	if *version {
		fmt.Printf("Ospy %s\n", Version)
		fmt.Printf("Build time: %s\n", BuildTime)
		fmt.Printf("Git commit: %s\n", GitCommit)
		return
	}

	// Check if config file exists
	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		fmt.Printf("❌ Configuration file not found: %s\n\n", *configPath)
		fmt.Println("To get started:")
		fmt.Println("1. Create a config file:")
		fmt.Printf("   mkdir -p %s\n", filepath.Dir(*configPath))
		fmt.Printf("   # Copy example config to %s\n\n", *configPath)
		fmt.Println("2. Run with config:")
		fmt.Printf("   %s -config %s\n\n", os.Args[0], *configPath)
		fmt.Println("3. Or run with -help for more information")
		os.Exit(1)
	}

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

	// Initialize notifiers
	var notifiers []notifier.Notifier

	// Email notifier
	if cfg.Notifications.Email.Enabled {
		emailNotifier := notifier.NewEmailNotifier(
			cfg.Notifications.Email.SMTPHost,
			cfg.Notifications.Email.SMTPPort,
			cfg.Notifications.Email.Username,
			cfg.Notifications.Email.Password,
			cfg.Notifications.Email.From,
			cfg.Notifications.Email.To,
		)
		notifiers = append(notifiers, emailNotifier)
		log.Printf("Email notifications enabled: %s", cfg.Notifications.Email.SMTPHost)
	}

	// Telegram notifier
	if cfg.Notifications.Telegram.Enabled {
		telegramNotifier := notifier.NewTelegramNotifier(
			cfg.Notifications.Telegram.BotToken,
			cfg.Notifications.Telegram.ChatID,
		)
		notifiers = append(notifiers, telegramNotifier)
		log.Println("Telegram notifications enabled")
	}

	// Create notification manager
	notifManager := notifier.NewManager(notifiers)

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

	// Connect worker pool results to monitor and notifications
	go func() {
		for result := range workerPool.Results() {
			// Send to monitor for logging
			mon.GetResults() <- result

			// Send to notification manager
			notifResult := notifier.CheckResult{
				WebsiteName:  result.WebsiteName,
				URL:          result.URL,
				Status:       result.Status,
				ResponseTime: result.ResponseTime,
				Error:        result.Error,
				Timestamp:    result.Timestamp,
				IsUp:         result.IsUp,
				Message:      result.Message,
			}
			notifManager.HandleResult(notifResult)
		}
	}()
	// Start monitoring
	mon.Start()

	// Start web server if enabled
	if cfg.Web.Enabled {
		webServer := web.NewServer(storage, cfg.Web.Port)
		go func() {
			log.Printf("Starting web dashboard on http://%s:%d", cfg.Web.Host, cfg.Web.Port)
			if err := webServer.Start(); err != nil {
				log.Printf("Web server error: %v", err)
			}
		}()
	}

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

	// Weekly summary report
	if len(notifiers) > 0 {
		go func() {
			// Send first report after 1 hour, then weekly
			time.Sleep(1 * time.Hour)

			ticker := time.NewTicker(7 * 24 * time.Hour) // Weekly
			defer ticker.Stop()

			for {
				stats, err := mon.GetAllStats(7 * 24 * time.Hour) // Last 7 days
				if err != nil {
					log.Printf("Failed to get stats for summary: %v", err)
				} else if len(stats) > 0 {
					notifManager.SendSummaryReport(stats)
				}

				<-ticker.C
			}
		}()
	}
	log.Printf("Ospy started - monitoring %d websites every %v",
		len(websites), cfg.Monitoring.Interval)
	log.Printf("Data stored in: %s", cfg.Storage.Path)
	log.Printf("Notifications: %d providers enabled", len(notifiers))
	if cfg.Web.Enabled {
		log.Printf("Web dashboard: http://%s:%d", cfg.Web.Host, cfg.Web.Port)
	}
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
