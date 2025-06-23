package monitor

import (
	// "fmt"
	"log"
	"time"

	"github.com/ravikantchauhan246/ospy/internal/storage"
)

// Monitor combines checking, scheduling, and storage
type Monitor struct {
	scheduler *Scheduler
	storage   storage.Storage
	resultsCh chan CheckResult
}

// NewMonitor creates a new monitoring service
func NewMonitor(workerPool *WorkerPool, websites []Website, interval time.Duration, storage storage.Storage) *Monitor {
	scheduler := NewScheduler(workerPool, websites, interval)
	
	return &Monitor{
		scheduler: scheduler,
		storage:   storage,
		resultsCh: make(chan CheckResult, 100),
	}
}

// Start begins monitoring
func (m *Monitor) Start() {
	// Start processing results
	go m.processResults()
	
	// Start scheduler
	m.scheduler.Start()
	
	log.Println("Monitor started")
}

// processResults handles incoming check results
func (m *Monitor) processResults() {
	for result := range m.resultsCh {
		// Log result
		status := "✅"
		if !result.IsUp {
			status = "❌"
		}
		
		log.Printf("%s %s (%s) - %s (Time: %v)", 
			status, result.WebsiteName, result.URL, result.Message, result.ResponseTime)
		
		if result.Error != nil {
			log.Printf("   Error: %v", result.Error)
		}

		// Save to storage
		logEntry := storage.MonitorLog{
			WebsiteName:  result.WebsiteName,
			URL:          result.URL,
			Status:       result.Status,
			ResponseTime: result.ResponseTime.Microseconds(),
			IsUp:         result.IsUp,
			Message:      result.Message,
			Timestamp:    result.Timestamp,
		}
		
		if result.Error != nil {
			logEntry.Error = result.Error.Error()
		}

		if err := m.storage.SaveLog(logEntry); err != nil {
			log.Printf("Failed to save log: %v", err)
		}
	}
}

// GetResults returns the results channel
func (m *Monitor) GetResults() chan<- CheckResult {
	return m.resultsCh
}

// Stop stops monitoring
func (m *Monitor) Stop() {
	m.scheduler.Stop()
	close(m.resultsCh)
	log.Println("Monitor stopped")
}

// GetStats retrieves statistics for a website
func (m *Monitor) GetStats(websiteName string, duration time.Duration) (storage.WebsiteStats, error) {
	return m.storage.GetStats(websiteName, duration)
}

// GetAllStats retrieves statistics for all websites
func (m *Monitor) GetAllStats(duration time.Duration) ([]storage.WebsiteStats, error) {
	return m.storage.GetAllStats(duration)
}