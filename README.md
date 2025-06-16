# Ospy - Website Availability Monitor

A lightweight, concurrent website availability monitor built in Go that provides real-time monitoring, instant notifications, and comprehensive uptime tracking.

## 🌟 Overview

Ospy is designed to monitor multiple websites simultaneously using Go's powerful concurrency features. It tracks response times, availability statistics, and sends instant notifications when sites go down, making it perfect for production monitoring and SLA compliance.

## 🎯 Use Cases

- **Production Website Monitoring** - Keep track of your critical websites 24/7
- **SLA Compliance Tracking** - Monitor and report uptime statistics
- **Performance Monitoring** - Track response times and identify bottlenecks
- **Automated Incident Detection** - Get notified immediately when issues occur
- **Uptime Reporting** - Generate comprehensive availability reports

## ✨ Features

### Core Monitoring
- ✅ **Concurrent Monitoring** - Monitor multiple URLs simultaneously using goroutines
- ✅ **Configurable Intervals** - Set check intervals from seconds to hours
- ✅ **Response Time Tracking** - Monitor and analyze website performance
- ✅ **HTTP Status Monitoring** - Track status codes and identify issues
- ✅ **Custom Timeouts** - Configure timeout settings per website
- ✅ **Retry Logic** - Built-in exponential backoff for reliable monitoring

### Notification System
- 📧 **Email Alerts** - SMTP-based notifications with HTML templates
- 📱 **Telegram Bot** - Instant messaging via Telegram API
- 🔔 **Alert Rules** - Configurable thresholds and escalation policies
- 🚫 **Rate Limiting** - Prevents notification spam

### Data Storage
- 💾 **SQLite Database** - Structured data storage with powerful queries
- 📊 **JSON Logging** - Alternative lightweight storage option
- 📈 **Historical Data** - Long-term uptime and performance tracking
- 📋 **Statistics** - Uptime percentages and SLA calculations

### Interfaces
- 💻 **CLI Tool** - Command-line interface for easy management
- 🌐 **Web Dashboard** - Real-time status monitoring interface
- 📱 **REST API** - Programmatic access to monitoring data

## 🚀 Quick Start

### Prerequisites
- Go 1.24.4 or later
- SQLite3 (for database storage)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/ravikantchauhan246/ospy.git
cd ospy
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the application:
```bash
go build -o ospy cmd/ospy/main.go
```

### Configuration

Create a `configs/config.yaml` file:

```yaml
# Example configuration
check_interval: 30s
timeout: 10s
max_retries: 3

websites:
  - name: "Example Site"
    url: "https://example.com"
    expected_status: 200
    timeout: 5s

notifications:
  email:
    enabled: true
    smtp_host: "smtp.gmail.com"
    smtp_port: 587
    username: "your-email@gmail.com"
    password: "your-app-password"
    
  telegram:
    enabled: true
    bot_token: "your-bot-token"
    chat_id: "your-chat-id"

database:
  type: "sqlite"
  path: "ospy.db"
```

### Running

#### CLI Mode
```bash
# Start monitoring
./ospy monitor

# Check status
./ospy status

# View statistics
./ospy stats
```

#### Web Dashboard
```bash
# Start web server
./ospy serve --port 8080
```

Then open `http://localhost:8080` in your browser.

## 📊 Usage Examples

### Adding Websites to Monitor

```bash
# Add a website via CLI
./ospy add --name "My Website" --url "https://mysite.com" --interval 60s

# Add multiple websites via config file
./ospy config --file configs/websites.yaml
```

### Setting Up Notifications

```bash
# Test email notifications
./ospy test-email --to "admin@example.com"

# Test Telegram notifications
./ospy test-telegram --message "Test notification"
```

### Viewing Reports

```bash
# Generate uptime report
./ospy report --period 30d --format json

# Export historical data
./ospy export --start 2025-01-01 --end 2025-01-31 --format csv
```

## 🏗️ Project Structure

```
ospy/
├── cmd/
│   └── ospy/
│       └── main.go          # Application entry point
├── configs/
│   └── config.yaml          # Configuration files
├── internal/
│   ├── config/
│   │   └── config.go        # Configuration management
│   ├── monitor/
│   │   └── checker.go       # Website monitoring logic
│   ├── notifier/
│   │   ├── email.go         # Email notification service
│   │   └── telegram.go      # Telegram notification service
│   ├── storage/
│   │   └── sqlite.go        # Database operations
│   └── web/
│       └── ...              # Web dashboard components
├── go.mod
├── go.sum
└── README.md
```

## 🔧 Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `check_interval` | How often to check websites | `30s` |
| `timeout` | Request timeout | `10s` |
| `max_retries` | Maximum retry attempts | `3` |
| `concurrent_checks` | Number of concurrent monitors | `10` |
| `alert_threshold` | Failure threshold for alerts | `3` |

## 📈 API Reference

### REST Endpoints

- `GET /api/status` - Get current status of all monitored sites
- `GET /api/stats` - Get uptime statistics
- `POST /api/websites` - Add new website to monitor
- `DELETE /api/websites/{id}` - Remove website from monitoring
- `GET /api/history/{id}` - Get historical data for a website

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🛟 Support

- 📧 Email: [ravikantchauhan246@gmail.com](mailto:ravikantchauhan246@gmail.com)
- 🐛 Issues: [GitHub Issues](https://github.com/ravikantchauhan246/ospy/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/ravikantchauhan246/ospy/discussions)

## 🙏 Acknowledgments

- Built with [Go](https://golang.org/)
- Uses [Gorilla Mux](https://github.com/gorilla/mux) for HTTP routing
- SQLite for lightweight database storage
- Telegram Bot API for instant notifications

---

**Made with ❤️ by [Ravikant Chauhan](https://github.com/ravikantchauhan246)**
