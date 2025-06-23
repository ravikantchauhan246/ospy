package storage

import (
	"time"
)

// MonitorLog represents a monitoring log entry
type MonitorLog struct {
	ID           int64     `json:"id"`
	WebsiteName  string    `json:"website_name"`
	URL          string    `json:"url"`
	Status       int       `json:"status"`
	ResponseTime int64     `json:"response_time"` // microseconds
	IsUp         bool      `json:"is_up"`
	Error        string    `json:"error"`
	Message      string    `json:"message"`
	Timestamp    time.Time `json:"timestamp"`
}

// WebsiteStats represents statistics for a website
type WebsiteStats struct {
	WebsiteName     string    `json:"website_name"`
	URL             string    `json:"url"`
	TotalChecks     int       `json:"total_checks"`
	SuccessfulChecks int      `json:"successful_checks"`
	UptimePercent   float64   `json:"uptime_percent"`
	AvgResponseTime float64   `json:"avg_response_time"`
	LastCheck       time.Time `json:"last_check"`
	LastStatus      string    `json:"last_status"`
}

// Storage interface defines storage operations
type Storage interface {
	SaveLog(log MonitorLog) error
	GetLogs(websiteName string, limit int) ([]MonitorLog, error)
	GetStats(websiteName string, duration time.Duration) (WebsiteStats, error)
	GetAllStats(duration time.Duration) ([]WebsiteStats, error)
	Cleanup(retentionDays int) error
	Close() error
}