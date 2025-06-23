package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/ravikantchauhan246/ospy/internal/storage"
)

// Server provides a web interface
type Server struct {
	storage storage.Storage
	port    int
}

// NewServer creates a new web server
func NewServer(storage storage.Storage, port int) *Server {
	return &Server{
		storage: storage,
		port:    port,
	}
}

// Start starts the web server
func (s *Server) Start() error {
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/api/stats", s.handleStats)
	http.HandleFunc("/api/logs", s.handleLogs)
	
	fmt.Printf("Web dashboard available at: http://localhost:%d\n", s.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}

// handleIndex serves the main dashboard page
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	stats, err := s.storage.GetAllStats(24 * time.Hour)
	if err != nil {
		http.Error(w, "Failed to get stats", http.StatusInternalServerError)
		return
	}

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Ospy Dashboard</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; }
        .header { background: #2c3e50; color: white; padding: 20px; border-radius: 8px; margin-bottom: 20px; }
        .stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; }
        .stat-card { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .stat-card h3 { margin-top: 0; color: #2c3e50; }
        .status-up { color: #27ae60; font-weight: bold; }
        .status-down { color: #e74c3c; font-weight: bold; }
        .metric { display: flex; justify-content: space-between; margin: 10px 0; }
        .metric-label { color: #7f8c8d; }
        .metric-value { font-weight: bold; }
        .footer { text-align: center; margin-top: 40px; color: #7f8c8d; }
    </style>
    <script>
        function refreshData() {
            location.reload();
        }
        setInterval(refreshData, 60000); // Refresh every minute
    </script>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>📊 Ospy Website Monitor</h1>
            <p>Real-time website monitoring dashboard</p>
        </div>
        
        {{if .}}
        <div class="stats-grid">
            {{range .}}
            <div class="stat-card">
                <h3>{{.WebsiteName}}</h3>
                <div class="metric">
                    <span class="metric-label">Status:</span>
                    <span class="metric-value {{if eq .LastStatus "UP"}}status-up{{else}}status-down{{end}}">
                        {{if eq .LastStatus "UP"}}🟢 UP{{else}}🔴 DOWN{{end}}
                    </span>
                </div>
                <div class="metric">
                    <span class="metric-label">URL:</span>
                    <span class="metric-value">{{.URL}}</span>
                </div>
                <div class="metric">
                    <span class="metric-label">Uptime (24h):</span>
                    <span class="metric-value">{{printf "%.2f%%" .UptimePercent}}</span>
                </div>
                <div class="metric">
                    <span class="metric-label">Avg Response:</span>
                    <span class="metric-value">{{printf "%.0fms" .AvgResponseTime}}</span>
                </div>
                <div class="metric">
                    <span class="metric-label">Total Checks:</span>
                    <span class="metric-value">{{.TotalChecks}}</span>
                </div>
                <div class="metric">
                    <span class="metric-label">Last Check:</span>
                    <span class="metric-value">{{.LastCheck.Format "15:04:05"}}</span>
                </div>
            </div>
            {{end}}
        </div>
        {{else}}
        <div class="stat-card">
            <h3>No Data Available</h3>
            <p>No monitoring data found. Check if monitoring is running.</p>
        </div>
        {{end}}
        
        <div class="footer">
            <p>Last updated: {{.Now.Format "2006-01-02 15:04:05"}} | Auto-refresh: 60s</p>
        </div>
    </div>
</body>
</html>`

	t, err := template.New("dashboard").Parse(tmpl)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Stats []storage.WebsiteStats
		Now   time.Time
	}{
		Stats: stats,
		Now:   time.Now(),
	}

	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, data.Stats)
}

// handleStats serves statistics as JSON
func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	durationStr := r.URL.Query().Get("duration")
	duration := 24 * time.Hour // default
	
	if durationStr != "" {
		if hours, err := strconv.Atoi(durationStr); err == nil {
			duration = time.Duration(hours) * time.Hour
		}
	}

	stats, err := s.storage.GetAllStats(duration)
	if err != nil {
		http.Error(w, "Failed to get stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// handleLogs serves recent logs as JSON
func (s *Server) handleLogs(w http.ResponseWriter, r *http.Request) {
	websiteName := r.URL.Query().Get("website")
	limitStr := r.URL.Query().Get("limit")
	
	limit := 50 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if websiteName == "" {
		http.Error(w, "website parameter required", http.StatusBadRequest)
		return
	}

	logs, err := s.storage.GetLogs(websiteName, limit)
	if err != nil {
		http.Error(w, "Failed to get logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}