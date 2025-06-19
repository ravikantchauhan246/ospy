package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Monitoring    MonitoringConfig   `yaml:"monitoring"`
	Websites      []WebsiteConfig    `yaml:"websites"`
	Notifications NotificationConfig `yaml:"notifications"`
	Storage       StorageConfig      `yaml:"storage"`
	Web           WebConfig          `yaml:"web"`
	Logging       LoggingConfig      `yaml:"logging"`
}

// MonitoringConfig contains monitoring settings
type MonitoringConfig struct {
	Interval time.Duration `yaml:"interval"`
	Timeout  time.Duration `yaml:"timeout"`
	Retries  int           `yaml:"retries"`
	Workers  int           `yaml:"workers"`
}

// WebsiteConfig represents a website to monitor
type WebsiteConfig struct {
	Name           string            `yaml:"name"`
	URL            string            `yaml:"url"`
	Method         string            `yaml:"method"`
	Headers        map[string]string `yaml:"headers"`
	ExpectedStatus int               `yaml:"expected_status"`
	CheckContent   string            `yaml:"check_content"`
	Timeout        time.Duration     `yaml:"timeout"`
}

// NotificationConfig contains notification settings
type NotificationConfig struct {
	Email    EmailConfig    `yaml:"email"`
	Telegram TelegramConfig `yaml:"telegram"`
}

// EmailConfig contains email notification settings
type EmailConfig struct {
	Enabled  bool     `yaml:"enabled"`
	SMTPHost string   `yaml:"smtp_host"`
	SMTPPort int      `yaml:"smtp_port"`
	Username string   // Loaded from environment variable
	Password string   // Loaded from environment variable
	From     string   `yaml:"from"`
	To       []string `yaml:"to"`
}

// TelegramConfig contains Telegram notification settings
type TelegramConfig struct {
	Enabled  bool   `yaml:"enabled"`
	BotToken string // Loaded from environment variable
	ChatID   string `yaml:"chat_id"`
}

// StorageConfig contains storage settings
type StorageConfig struct {
	Type          string `yaml:"type"`
	Path          string `yaml:"path"`
	RetentionDays int    `yaml:"retention_days"`
}

// WebConfig contains web server settings
type WebConfig struct {
	Enabled bool   `yaml:"enabled"`
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Level      string `yaml:"level"`
	File       string `yaml:"file"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
}

// Load reads configuration from file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Set defaults

	config.Notifications.Email.Username = os.Getenv("SMTP_USERNAME")
	config.Notifications.Email.Password = os.Getenv("SMTP_PASSWORD")
	config.Notifications.Telegram.BotToken = os.Getenv("TELEGRAM_BOT_TOKEN")

	if config.Monitoring.Interval == 0 {
		config.Monitoring.Interval = 5 * time.Minute
	}
	if config.Monitoring.Timeout == 0 {
		config.Monitoring.Timeout = 30 * time.Second
	}
	if config.Monitoring.Workers == 0 {
		config.Monitoring.Workers = 10
	}
	if config.Monitoring.Retries == 0 {
		config.Monitoring.Retries = 3
	}

	return &config, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if len(c.Websites) == 0 {
		return fmt.Errorf("no websites configured")
	}

	for i, website := range c.Websites {
		if website.URL == "" {
			return fmt.Errorf("website %d: URL is required", i)
		}
		if website.Name == "" {
			return fmt.Errorf("website %d: Name is required", i)
		}
	}

	return nil
}
