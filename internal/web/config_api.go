package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ravikantchauhan246/ospy/internal/config"
)

// Website represents a website configuration for the API
type Website struct {
	ID             int               `json:"id"`
	Name           string            `json:"name"`
	URL            string            `json:"url"`
	Method         string            `json:"method"`
	ExpectedStatus int               `json:"expected_status"`
	Headers        map[string]string `json:"headers,omitempty"`
	Enabled        bool              `json:"enabled"`
}

// ConfigAPI handles configuration management endpoints
type ConfigAPI struct {
	configPath string
	config     *config.Config
}

// NewConfigAPI creates a new configuration API handler
func NewConfigAPI(configPath string, cfg *config.Config) *ConfigAPI {
	return &ConfigAPI{
		configPath: configPath,
		config:     cfg,
	}
}

// setupConfigRoutes sets up the configuration management routes
func (s *Server) setupConfigRoutes(configAPI *ConfigAPI) {
	// Enable CORS for all routes
	http.HandleFunc("/api/config/websites", corsMiddleware(configAPI.handleWebsites))
	http.HandleFunc("/api/config/websites/", corsMiddleware(configAPI.handleWebsiteByID))
	http.HandleFunc("/api/config/settings", corsMiddleware(configAPI.handleSettings))
}

// corsMiddleware adds CORS headers
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next(w, r)
	}
}

// handleWebsites handles GET /api/config/websites and POST /api/config/websites
func (api *ConfigAPI) handleWebsites(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		api.getWebsites(w, r)
	case "POST":
		api.addWebsite(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleWebsiteByID handles PUT /api/config/websites/{id} and DELETE /api/config/websites/{id}
func (api *ConfigAPI) handleWebsiteByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	idStr := r.URL.Path[len("/api/config/websites/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid website ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "PUT":
		api.updateWebsite(w, r, id)
	case "DELETE":
		api.deleteWebsite(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getWebsites returns all configured websites
func (api *ConfigAPI) getWebsites(w http.ResponseWriter, r *http.Request) {
	websites := make([]Website, len(api.config.Websites))
	for i, site := range api.config.Websites {
		websites[i] = Website{
			ID:             i,
			Name:           site.Name,
			URL:            site.URL,
			Method:         site.Method,
			ExpectedStatus: site.ExpectedStatus,
			Headers:        site.Headers,
			Enabled:        true, // All websites in config are enabled by default
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(websites)
}

// addWebsite adds a new website to the configuration
func (api *ConfigAPI) addWebsite(w http.ResponseWriter, r *http.Request) {
	var website Website
	if err := json.NewDecoder(r.Body).Decode(&website); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if website.Name == "" || website.URL == "" {
		http.Error(w, "Name and URL are required", http.StatusBadRequest)
		return
	}

	// Set defaults
	if website.Method == "" {
		website.Method = "GET"
	}
	if website.ExpectedStatus == 0 {
		website.ExpectedStatus = 200
	}
	// Add to config
	newSite := config.WebsiteConfig{
		Name:           website.Name,
		URL:            website.URL,
		Method:         website.Method,
		ExpectedStatus: website.ExpectedStatus,
		Headers:        website.Headers,
	}

	api.config.Websites = append(api.config.Websites, newSite)
	
	// Save configuration
	if err := api.saveConfig(); err != nil {
		http.Error(w, "Failed to save configuration", http.StatusInternalServerError)
		return
	}

	// Return the new website with ID
	website.ID = len(api.config.Websites) - 1
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(website)
}

// updateWebsite updates an existing website
func (api *ConfigAPI) updateWebsite(w http.ResponseWriter, r *http.Request, id int) {
	if id < 0 || id >= len(api.config.Websites) {
		http.Error(w, "Website not found", http.StatusNotFound)
		return
	}

	var website Website
	if err := json.NewDecoder(r.Body).Decode(&website); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Update the website
	api.config.Websites[id].Name = website.Name
	api.config.Websites[id].URL = website.URL
	api.config.Websites[id].Method = website.Method
	api.config.Websites[id].ExpectedStatus = website.ExpectedStatus
	api.config.Websites[id].Headers = website.Headers

	// Save configuration
	if err := api.saveConfig(); err != nil {
		http.Error(w, "Failed to save configuration", http.StatusInternalServerError)
		return
	}

	website.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(website)
}

// deleteWebsite removes a website from the configuration
func (api *ConfigAPI) deleteWebsite(w http.ResponseWriter, r *http.Request, id int) {
	if id < 0 || id >= len(api.config.Websites) {
		http.Error(w, "Website not found", http.StatusNotFound)
		return
	}

	// Remove the website
	api.config.Websites = append(api.config.Websites[:id], api.config.Websites[id+1:]...)

	// Save configuration
	if err := api.saveConfig(); err != nil {
		http.Error(w, "Failed to save configuration", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleSettings returns monitoring settings
func (api *ConfigAPI) handleSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	settings := map[string]interface{}{
		"monitoring": map[string]interface{}{
			"interval": api.config.Monitoring.Interval.String(),
			"timeout":  api.config.Monitoring.Timeout.String(),
			"retries":  api.config.Monitoring.Retries,
			"workers":  api.config.Monitoring.Workers,
		},
		"notifications": map[string]interface{}{
			"email_enabled":    api.config.Notifications.Email.Enabled,
			"telegram_enabled": api.config.Notifications.Telegram.Enabled,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

// saveConfig saves the current configuration to file
func (api *ConfigAPI) saveConfig() error {
	return config.SaveConfig(api.config, api.configPath)
}
