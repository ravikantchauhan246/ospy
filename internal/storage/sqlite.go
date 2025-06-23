package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// SQLiteStorage implements Storage interface using SQLite
type SQLiteStorage struct {
	db *sql.DB
}

// NewSQLiteStorage creates a new SQLite storage
func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	storage := &SQLiteStorage{db: db}

	if err := storage.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return storage, nil
}

// createTables creates the necessary database tables
func (s *SQLiteStorage) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS monitor_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		website_name TEXT NOT NULL,
		url TEXT NOT NULL,
		status INTEGER,
		response_time INTEGER,
		is_up BOOLEAN,
		error TEXT,
		message TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_website_timestamp ON monitor_logs(website_name, timestamp);
	CREATE INDEX IF NOT EXISTS idx_timestamp ON monitor_logs(timestamp);
	`

	_, err := s.db.Exec(query)
	return err
}

// SaveLog saves a monitoring log entry
func (s *SQLiteStorage) SaveLog(log MonitorLog) error {
	query := `
	INSERT INTO monitor_logs (website_name, url, status, response_time, is_up, error, message, timestamp)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query,
		log.WebsiteName,
		log.URL,
		log.Status,
		log.ResponseTime,
		log.IsUp,
		log.Error,
		log.Message,
		log.Timestamp)

	return err
}

// GetLogs retrieves recent logs for a website
func (s *SQLiteStorage) GetLogs(websiteName string, limit int) ([]MonitorLog, error) {
	query := `
	SELECT id, website_name, url, status, response_time, is_up, error, message, timestamp
	FROM monitor_logs
	WHERE website_name = ?
	ORDER BY timestamp DESC
	LIMIT ?`

	rows, err := s.db.Query(query, websiteName, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []MonitorLog
	for rows.Next() {
		var log MonitorLog
		var errorStr sql.NullString
		err := rows.Scan(
			&log.ID,
			&log.WebsiteName,
			&log.URL,
			&log.Status,
			&log.ResponseTime,
			&log.IsUp,
			&errorStr,
			&log.Message,
			&log.Timestamp,
		)
		if err != nil {
			return nil, err
		}

		if errorStr.Valid {
			log.Error = errorStr.String
		}

		logs = append(logs, log)
	}

	return logs, nil
}

// GetStats calculates statistics for a website
func (s *SQLiteStorage) GetStats(websiteName string, duration time.Duration) (WebsiteStats, error) {
	since := time.Now().Add(-duration)

	query := `
	SELECT 
		COUNT(*) as total_checks,
		SUM(CASE WHEN is_up = 1 THEN 1 ELSE 0 END) as successful_checks,
		AVG(response_time) as avg_response_time,
		MAX(timestamp) as last_check
	FROM monitor_logs
	WHERE website_name = ? AND timestamp >= ?`

	var stats WebsiteStats
	var avgResponseTime sql.NullFloat64
	var lastCheck sql.NullTime

	err := s.db.QueryRow(query, websiteName, since).Scan(
		&stats.TotalChecks,
		&stats.SuccessfulChecks,
		&avgResponseTime,
		&lastCheck,
	)

	if err != nil {
		return stats, err
	}

	stats.WebsiteName = websiteName

	// Get URL from latest record
	urlQuery := `SELECT url FROM monitor_logs WHERE website_name = ? ORDER BY timestamp DESC LIMIT 1`
	s.db.QueryRow(urlQuery, websiteName).Scan(&stats.URL)

	if stats.TotalChecks > 0 {
		stats.UptimePercent = float64(stats.SuccessfulChecks) / float64(stats.TotalChecks) * 100
	}

	if avgResponseTime.Valid {
		stats.AvgResponseTime = float64(avgResponseTime.Float64) / 1000 // Convert to milliseconds
	}

	if lastCheck.Valid {
		stats.LastCheck = lastCheck.Time

		// Get last status
		statusQuery := `SELECT is_up FROM monitor_logs WHERE website_name = ? ORDER BY timestamp DESC LIMIT 1`
		var isUp bool
		if err := s.db.QueryRow(statusQuery, websiteName).Scan(&isUp); err == nil {
			if isUp {
				stats.LastStatus = "UP"
			} else {
				stats.LastStatus = "DOWN"
			}
		}
	}

	return stats, nil
}

// GetAllStats retrieves statistics for all websites
func (s *SQLiteStorage) GetAllStats(duration time.Duration) ([]WebsiteStats, error) {
	since := time.Now().Add(-duration)

	// Get all website names
	websiteQuery := `SELECT DISTINCT website_name FROM monitor_logs WHERE timestamp >= ?`
	rows, err := s.db.Query(websiteQuery, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allStats []WebsiteStats
	for rows.Next() {
		var websiteName string
		if err := rows.Scan(&websiteName); err != nil {
			continue
		}

		stats, err := s.GetStats(websiteName, duration)
		if err != nil {
			continue
		}

		allStats = append(allStats, stats)
	}

	return allStats, nil
}

// Cleanup removes old log entries
func (s *SQLiteStorage) Cleanup(retentionDays int) error {
	cutoff := time.Now().AddDate(0, 0, -retentionDays)

	query := `DELETE FROM monitor_logs WHERE timestamp < ?`
	result, err := s.db.Exec(query, cutoff)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("Cleaned up %d old log entries\n", rowsAffected)

	return nil
}

// Close closes the database connection
func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}