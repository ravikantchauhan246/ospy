package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Monitoring    MonitoringConfig    `yaml:"monitoring"`
	Websites      []WebsiteConfig     `yaml:"websites"`
	Storage       StorageConfig       `yaml:"storage"`
	Notifications NotificationConfig  `yaml:"notifications"`
}

// MonitoringConfig holds monitoring settings
type MonitoringConfig struct {
	Interval time.Duration `yaml:"interval"`
	Timeout  time.Duration `yaml:"timeout"`
	Workers  int           `yaml:"workers"`
}

// WebsiteConfig holds individual website configuration
type WebsiteConfig struct {
	Name           string            `yaml:"name"`
	URL            string            `yaml:"url"`
	Method         string            `yaml:"method"`
	Headers        map[string]string `yaml:"headers"`
	ExpectedStatus int               `yaml:"expected_status"`
	CheckContent   string            `yaml:"check_content"`
	Timeout        time.Duration     `yaml:"timeout"`
}

// StorageConfig holds storage settings
type StorageConfig struct {
	Path          string `yaml:"path"`
	RetentionDays int    `yaml:"retention_days"`
}

// NotificationConfig holds notification settings
type NotificationConfig struct {
	Email    EmailConfig    `yaml:"email"`
	Telegram TelegramConfig `yaml:"telegram"`
}

// EmailConfig holds email notification settings
type EmailConfig struct {
	Enabled  bool     `yaml:"enabled"`
	SMTPHost string   `yaml:"smtp_host"`
	SMTPPort int      `yaml:"smtp_port"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	From     string   `yaml:"from"`
	To       []string `yaml:"to"`
}

// TelegramConfig holds Telegram notification settings
type TelegramConfig struct {
	Enabled  bool   `yaml:"enabled"`
	BotToken string `yaml:"bot_token"`
	ChatID   string `yaml:"chat_id"`
}

// Load loads configuration from a file
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
	if config.Monitoring.Interval == 0 {
		config.Monitoring.Interval = 5 * time.Minute
	}
	if config.Monitoring.Timeout == 0 {
		config.Monitoring.Timeout = 30 * time.Second
	}
	if config.Monitoring.Workers == 0 {
		config.Monitoring.Workers = 5
	}
	if config.Storage.Path == "" {
		config.Storage.Path = "data/ospy.db"
	}
	if config.Storage.RetentionDays == 0 {
		config.Storage.RetentionDays = 30
	}

	// Set website defaults
	for i := range config.Websites {
		if config.Websites[i].Method == "" {
			config.Websites[i].Method = "GET"
		}
		if config.Websites[i].ExpectedStatus == 0 {
			config.Websites[i].ExpectedStatus = 200
		}
		if config.Websites[i].Timeout == 0 {
			config.Websites[i].Timeout = config.Monitoring.Timeout
		}
	}

	return &config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if len(c.Websites) == 0 {
		return fmt.Errorf("no websites configured")
	}

	for i, website := range c.Websites {
		if website.Name == "" {
			return fmt.Errorf("website %d: name is required", i)
		}
		if website.URL == "" {
			return fmt.Errorf("website %d: URL is required", i)
		}
	}

	// Validate email config if enabled
	if c.Notifications.Email.Enabled {
		if c.Notifications.Email.SMTPHost == "" {
			return fmt.Errorf("email notifications enabled but smtp_host not configured")
		}
		if c.Notifications.Email.Username == "" {
			return fmt.Errorf("email notifications enabled but username not configured")
		}
		if c.Notifications.Email.Password == "" {
			return fmt.Errorf("email notifications enabled but password not configured")
		}
		if len(c.Notifications.Email.To) == 0 {
			return fmt.Errorf("email notifications enabled but no recipients configured")
		}
	}

	// Validate Telegram config if enabled
	if c.Notifications.Telegram.Enabled {
		if c.Notifications.Telegram.BotToken == "" {
			return fmt.Errorf("telegram notifications enabled but bot_token not configured")
		}
		if c.Notifications.Telegram.ChatID == "" {
			return fmt.Errorf("telegram notifications enabled but chat_id not configured")
		}
	}

	return nil
}