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
- 📧 **Email Alerts** - SMTP-based notifications with configurable templates
- 📱 **Telegram Bot** - Instant messaging via Telegram API
- 🔔 **Alert Rules** - Configurable thresholds and escalation policies
- 🚫 **Rate Limiting** - Prevents notification spam

### Data Storage
- 💾 **SQLite Database** - Structured data storage with powerful queries
-  **Historical Data** - Long-term uptime and performance tracking
- 📋 **Statistics** - Uptime percentages and SLA calculations

### Deployment Options
- 💻 **CLI Tool** - Command-line interface for easy management
- 🐳 **Docker Support** - Containerized deployment with Docker Compose
- 🔧 **System Service** - Linux systemd and Windows service support
- 🌐 **Web Dashboard** - Real-time status monitoring (coming soon)

## 🚀 Quick Start

### Prerequisites
- Go 1.24.4 or later (for building from source)
- Or download pre-built binaries from [releases](https://github.com/ravikantchauhan246/ospy/releases)

### Installation

#### Option 1: Download Pre-built Binary
```bash
# Linux
wget https://github.com/ravikantchauhan246/ospy/releases/latest/download/ospy-linux-amd64.tar.gz
tar -xzf ospy-linux-amd64.tar.gz
cd ospy-linux-amd64

# Windows - Download and extract ospy-windows-amd64.zip
```

#### Option 2: Build from Source
```bash
git clone https://github.com/ravikantchauhan246/ospy.git
cd ospy
go mod tidy
go build -o ospy cmd/ospy/main.go
```

#### Option 3: Docker
```bash
git clone https://github.com/ravikantchauhan246/ospy.git
cd ospy
cp .env.example .env
# Edit .env with your credentials
docker-compose up -d
```

### Configuration

Create or edit `configs/config.yaml`:

```yaml
monitoring:
  interval: 30s        # Check every 30 seconds
  timeout: 10s         # Request timeout
  retries: 3           # Number of retries
  workers: 10          # Concurrent workers

websites:
  - name: "My Website"
    url: "https://example.com"
    method: "GET"
    expected_status: 200
    timeout: 5s

  - name: "API Endpoint"
    url: "https://api.example.com/health"
    method: "GET" 
    expected_status: 200
    headers:
      Authorization: "Bearer token"

notifications:
  email:
    enabled: true
    smtp_host: "smtp.gmail.com"
    smtp_port: 587
    from: "alerts@yourdomain.com"
    to: ["admin@yourdomain.com"]
    
  telegram:
    enabled: false
    chat_id: "your-chat-id"

storage:
  type: "sqlite"
  path: "./data/ospy.db"
  retention_days: 30

web:
  enabled: true
  host: "0.0.0.0"
  port: 8080

logging:
  level: "info"
  file: "./logs/ospy.log"
```

### Environment Variables

Set these for sensitive credentials:
```bash
# Email configuration
export SMTP_USERNAME="your-email@gmail.com"
export SMTP_PASSWORD="your-app-password"

# Telegram configuration (optional)
export TELEGRAM_BOT_TOKEN="123456789:ABCdefGHIjklMNOpqrsTUVwxyz"
```

### Running

```bash
# Start monitoring with example config
./ospy -config example-config.yaml

# Start monitoring with default config (if configs/config.yaml exists)
./ospy

# Check version
./ospy -version

# Show help
./ospy -help
```

## 📊 Usage Examples

### Basic Monitoring
```bash
# Monitor with example config
./ospy -config example-config.yaml

# Monitor with custom config
./ospy -config my-custom-config.yaml

# Get help
./ospy -help

# Check version
./ospy -version
```

### Docker Deployment
```bash
# Quick start with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f ospy

# Stop monitoring
docker-compose down
```

### Linux Service Installation
```bash
# Download release and extract
wget https://github.com/ravikantchauhan246/ospy/releases/latest/download/ospy-linux-amd64.tar.gz
tar -xzf ospy-linux-amd64.tar.gz
cd ospy-linux-amd64

# Install as system service
sudo ./install.sh

# Manage service
sudo systemctl start ospy
sudo systemctl status ospy
sudo systemctl stop ospy
```

## 🏗️ Project Structure

```
ospy/
├── cmd/
│   └── ospy/
│       └── main.go              # Application entry point
├── configs/
│   └── config.yaml              # Configuration template
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── monitor/
│   │   ├── checker.go           # Website checking logic
│   │   ├── monitor.go           # Main monitoring coordinator
│   │   ├── scheduler.go         # Task scheduling
│   │   └── worker.go            # Worker pool management
│   ├── notifier/
│   │   ├── manager.go           # Notification manager
│   │   ├── email.go             # Email notifications
│   │   └── telegram.go          # Telegram notifications
│   ├── storage/
│   │   ├── models.go            # Data models
│   │   └── sqlite.go            # SQLite database operations
│   └── web/
│       └── server.go            # Web server (coming soon)
├── deploy/
│   ├── install.sh               # Linux installation script
│   └── ospy.service             # Systemd service file
├── scripts/
│   ├── build.sh                 # Build script (Linux/Mac)
│   ├── build.ps1                # Build script (Windows)
│   ├── release.ps1              # Release packaging script
│   └── build.config             # Build configuration
├── docker-compose.yml           # Docker Compose configuration
├── Dockerfile                   # Docker image definition
├── Makefile                     # Build automation
├── .env.example                 # Environment variables template
├── DEPLOYMENT.md                # Deployment guide
├── LICENSE                      # MIT License
├── README.md                    # This file
├── go.mod                       # Go module definition
└── go.sum                       # Go module checksums
```

## 🔧 Configuration Options

### Monitoring Settings
| Option | Description | Default | Example |
|--------|-------------|---------|---------|
| `interval` | Check interval | `5m` | `30s`, `1m`, `5m` |
| `timeout` | Request timeout | `30s` | `5s`, `10s`, `30s` |
| `retries` | Retry attempts | `3` | `1`, `3`, `5` |
| `workers` | Concurrent workers | `10` | `5`, `10`, `20` |

### Website Configuration
| Option | Description | Required | Example |
|--------|-------------|----------|---------|
| `name` | Display name | ✅ | `"My Website"` |
| `url` | URL to monitor | ✅ | `"https://example.com"` |
| `method` | HTTP method | ❌ | `"GET"`, `"POST"` |
| `expected_status` | Expected HTTP status | ❌ | `200`, `404` |
| `timeout` | Per-site timeout | ❌ | `5s`, `10s` |
| `headers` | Custom headers | ❌ | `{"Auth": "Bearer token"}` |

## � Deployment

### Docker (Recommended)
```bash
# Using Docker Compose
docker-compose up -d

# Using Docker directly
docker build -t ospy .
docker run -d --name ospy -p 8080:8080 \
  -e SMTP_USERNAME="your-email" \
  -e SMTP_PASSWORD="your-password" \
  ospy
```

### Linux System Service
```bash
# Automated installation
sudo ./deploy/install.sh

# Manual installation
sudo cp ospy /usr/local/bin/
sudo cp configs/config.yaml /etc/ospy/
sudo systemctl enable ospy
sudo systemctl start ospy
```

### Windows
```bash
# Run as application
ospy.exe -config config.yaml

# Install as Windows Service (using NSSM)
nssm install Ospy "C:\path\to\ospy.exe" "-config C:\path\to\config.yaml"
nssm start Ospy
```

## 📈 Monitoring Output

Ospy provides real-time monitoring feedback:

```
2025/06/23 20:08:07 Monitor started
2025/06/23 20:08:07 Starting check of 5 websites
2025/06/23 20:08:07 Ospy started - monitoring 5 websites every 10s
2025/06/23 20:08:07 Data stored in: ./data/ospy.db
2025/06/23 20:08:07 ✅ GitHub (https://github.com) - Status 200 (Time: 237ms)
2025/06/23 20:08:07 ✅ Google (https://google.com) - Status 200 (Time: 729ms)
2025/06/23 20:08:08 ❌ Test Site (https://httpstat.us/500) - HTTP request failed
```

## 🛠️ Building and Development

### Build Scripts
```bash
# Cross-platform build
make build                    # Current platform
make cross-compile           # All platforms

# Windows PowerShell
scripts/build.ps1

# Linux/Mac
./scripts/build.sh
```

### Creating Releases
```bash
# Create packaged release
scripts/release.ps1 -Version "1.0.0"

# This creates:
# - releases/v1.0.0/ospy-v1.0.0-linux-amd64.tar.gz
# - releases/v1.0.0/ospy-v1.0.0-windows-amd64.zip
# - releases/v1.0.0/checksums.txt
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Setup
```bash
git clone https://github.com/ravikantchauhan246/ospy.git
cd ospy
go mod tidy
go run cmd/ospy/main.go -config configs/config.yaml
```

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🛟 Support

- 📧 Email: [ravikantchauhan246@gmail.com](mailto:ravikantchauhan246@gmail.com)
- 🐛 Issues: [GitHub Issues](https://github.com/ravikantchauhan246/ospy/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/ravikantchauhan246/ospy/discussions)
- 📖 Documentation: [Deployment Guide](DEPLOYMENT.md)

## 🙏 Acknowledgments

- Built with [Go](https://golang.org/)
- Uses [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) for pure Go SQLite
- Uses [gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3) for configuration
- Inspired by monitoring tools like Uptime Robot and StatusCake

## 🔮 Roadmap

- [ ] Web dashboard with real-time charts
- [ ] Prometheus metrics export
- [ ] Slack and Discord notifications
- [ ] Advanced alerting rules
- [ ] Multi-region monitoring
- [ ] REST API for management
- [ ] Custom check scripts

---

**Made with ❤️ by [Ravikant Chauhan](https://github.com/ravikantchauhan246)**
