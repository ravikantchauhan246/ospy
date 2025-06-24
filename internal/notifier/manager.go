package notifier

import (
	"log"
	"sync"
	"time"

	"github.com/ravikantchauhan246/ospy/internal/storage"
)

// Notifier interface for all notification types
type Notifier interface {
	IsEnabled() bool
	SendDownAlert(websiteName, url, message string) error
	SendUpAlert(websiteName, url string, downtime time.Duration) error
	SendSummaryReport(stats []storage.WebsiteStats) error
}

// Manager manages all notification services
type Manager struct {
	notifiers    []Notifier
	websiteState map[string]WebsiteState
	mutex        sync.RWMutex
}

// WebsiteState tracks the state of a website
type WebsiteState struct {
	IsUp      bool
	LastUp    time.Time
	LastDown  time.Time
	LastAlert time.Time
}

// NewManager creates a new notification manager
func NewManager(notifiers []Notifier) *Manager {
	return &Manager{
		notifiers:    notifiers,
		websiteState: make(map[string]WebsiteState),
	}
}

// HandleResult processes a check result and sends notifications if needed
func (m *Manager) HandleResult(result CheckResult) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	websiteName := result.WebsiteName
	currentState, exists := m.websiteState[websiteName]	// Initialize state if doesn't exist
	if !exists {
		currentState = WebsiteState{
			IsUp:     result.IsUp,
			LastUp:   time.Now(),
			LastDown: time.Now(),
		}
		m.websiteState[websiteName] = currentState
		log.Printf("🆕 Initialized state for %s: IsUp=%v", websiteName, result.IsUp)
		return // Don't send notifications on first check
	}

	log.Printf("🔄 State check for %s: Current=%v, New=%v", websiteName, currentState.IsUp, result.IsUp)

	// Check for state changes
	if !currentState.IsUp && result.IsUp {
		// Website came back up
		downtime := time.Since(currentState.LastDown)
		m.sendUpAlert(result.WebsiteName, result.URL, downtime)
		
		currentState.IsUp = true
		currentState.LastUp = time.Now()
		currentState.LastAlert = time.Now()
		
	} else if currentState.IsUp && !result.IsUp {
		// Website went down
		m.sendDownAlert(result.WebsiteName, result.URL, result.Message)		
		currentState.IsUp = false
		currentState.LastDown = time.Now()
		currentState.LastAlert = time.Now()
	} else if !currentState.IsUp && !result.IsUp {
		// Website is still down - for testing, send alert on second consecutive failure
		if currentState.LastAlert.IsZero() {
			log.Printf("🔔 Website still down, sending delayed alert for %s", websiteName)
			m.sendDownAlert(result.WebsiteName, result.URL, result.Message)
			currentState.LastAlert = time.Now()
		}
	}

	m.websiteState[websiteName] = currentState
}

// sendDownAlert sends down alerts to all enabled notifiers
func (m *Manager) sendDownAlert(websiteName, url, message string) {
	log.Printf("📧 Sending down alert for %s", websiteName)
	
	for _, notifier := range m.notifiers {
		if notifier.IsEnabled() {
			if err := notifier.SendDownAlert(websiteName, url, message); err != nil {
				log.Printf("Failed to send down alert: %v", err)
			}
		}
	}
}

// sendUpAlert sends up alerts to all enabled notifiers
func (m *Manager) sendUpAlert(websiteName, url string, downtime time.Duration) {
	log.Printf("📧 Sending up alert for %s (downtime: %v)", websiteName, downtime)
	
	for _, notifier := range m.notifiers {
		if notifier.IsEnabled() {
			if err := notifier.SendUpAlert(websiteName, url, downtime); err != nil {
				log.Printf("Failed to send up alert: %v", err)
			}
		}
	}
}

// SendSummaryReport sends summary reports to all enabled notifiers
func (m *Manager) SendSummaryReport(stats []storage.WebsiteStats) {
	log.Printf("📧 Sending summary report for %d websites", len(stats))
	
	for _, notifier := range m.notifiers {
		if notifier.IsEnabled() {
			if err := notifier.SendSummaryReport(stats); err != nil {
				log.Printf("Failed to send summary report: %v", err)
			}
		}
	}
}

// CheckResult represents a monitoring result (same as monitor.CheckResult)
type CheckResult struct {
	WebsiteName  string
	URL          string
	Status       int
	ResponseTime time.Duration
	Error        error
	Timestamp    time.Time
	IsUp         bool
	Message      string
}