# Deployment Guide

This guide covers different deployment methods for Ospy.

## Quick Start

### 1. Local Development
```bash
# Clone and build
git clone https://github.com/ravikantchauhan246/ospy.git
cd ospy
make build

# Run
./dist/ospy -config configs/config.yaml
```

### 2. Pre-built Binaries
Download from [releases](https://github.com/ravikantchauhan246/ospy/releases):
```bash
# Linux
wget https://github.com/ravikantchauhan246/ospy/releases/latest/download/ospy-linux-amd64.tar.gz
tar -xzf ospy-linux-amd64.tar.gz
./ospy -config config.yaml

# Windows
# Download ospy-windows-amd64.zip and extract
ospy.exe -config config.yaml
```

## Production Deployment

### Docker (Recommended)

1. **Using Docker Compose** (easiest):
```bash
# Clone repository
git clone https://github.com/ravikantchauhan246/ospy.git
cd ospy

# Configure environment
cp .env.example .env
# Edit .env with your credentials

# Start service
docker-compose up -d

# View logs
docker-compose logs -f
```

2. **Using Docker directly**:
```bash
# Build image
docker build -t ospy .

# Run container
docker run -d \
  --name ospy \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/configs:/app/configs \
  -e SMTP_USERNAME=your-email@example.com \
  -e SMTP_PASSWORD=your-password \
  ospy
```

### Linux System Service

1. **Automated installation**:
```bash
# Download and extract release
wget https://github.com/ravikantchauhan246/ospy/releases/latest/download/ospy-linux-amd64.tar.gz
tar -xzf ospy-linux-amd64.tar.gz
cd ospy-linux-amd64

# Run installer (requires root)
sudo ./install.sh

# Start service
sudo systemctl start ospy
sudo systemctl status ospy
```

2. **Manual installation**:
```bash
# Create user and directories
sudo useradd --system --home-dir /opt/ospy ospy
sudo mkdir -p /opt/ospy/{data,logs,configs}

# Copy binary and config
sudo cp ospy /opt/ospy/
sudo cp config.yaml /opt/ospy/configs/
sudo chown -R ospy:ospy /opt/ospy

# Install service
sudo cp deploy/ospy.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable ospy
sudo systemctl start ospy
```

### Windows Service

Use [NSSM](https://nssm.cc/) to run as a Windows service:
```cmd
# Install NSSM
# Download from https://nssm.cc/download

# Install service
nssm install Ospy "C:\path\to\ospy.exe"
nssm set Ospy Parameters "-config C:\path\to\config.yaml"
nssm set Ospy AppDirectory "C:\path\to\ospy"

# Start service
nssm start Ospy
```

## Configuration

### Environment Variables
Required for notifications:
```bash
export SMTP_USERNAME="your-email@example.com"
export SMTP_PASSWORD="your-app-password"
export TELEGRAM_BOT_TOKEN="123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
```

### Configuration File
Edit `configs/config.yaml`:
```yaml
monitoring:
  interval: 5m
  timeout: 30s

websites:
  - name: "My Website"
    url: "https://example.com"
    expected_status: 200

web:
  enabled: true
  port: 8080
```

## Monitoring

### Health Checks
- HTTP: `GET http://localhost:8080/health`
- CLI: `./ospy -version`

### Logs
- Docker: `docker-compose logs -f`
- Systemd: `journalctl -u ospy -f`
- File: Check `logs/ospy.log`

### Metrics
Access web dashboard at: `http://localhost:8080`

## Scaling

### Multiple Instances
Run multiple instances monitoring different websites:
```bash
# Instance 1 - Frontend services
./ospy -config configs/frontend.yaml

# Instance 2 - Backend services  
./ospy -config configs/backend.yaml
```

### Load Balancing
Use nginx to load balance multiple instances:
```nginx
upstream ospy {
    server 127.0.0.1:8080;
    server 127.0.0.1:8081;
}

server {
    listen 80;
    location / {
        proxy_pass http://ospy;
    }
}
```

## Security

### Firewall
Only expose necessary ports:
```bash
# Allow web interface
sudo ufw allow 8080/tcp

# Allow SSH (if needed)
sudo ufw allow 22/tcp
```

### SSL/TLS
Use nginx or traefik for SSL termination:
```yaml
# docker-compose.yml with traefik
version: '3.8'
services:
  ospy:
    labels:
      - traefik.enable=true
      - traefik.http.routers.ospy.rule=Host(`monitor.example.com`)
      - traefik.http.routers.ospy.tls.certresolver=letsencrypt
```

## Troubleshooting

### Common Issues

1. **Permission denied**:
```bash
sudo chown -R ospy:ospy /opt/ospy
sudo chmod +x /opt/ospy/ospy
```

2. **Port already in use**:
```bash
# Check what's using the port
sudo netstat -tulpn | grep :8080
# Kill the process or change port in config
```

3. **Database locked**:
```bash
# Stop service and check file permissions
sudo systemctl stop ospy
sudo chown ospy:ospy /opt/ospy/data/ospy.db
sudo systemctl start ospy
```

4. **Email not sending**:
- Check SMTP credentials
- Verify app passwords for Gmail
- Check firewall blocking SMTP ports

## Backup

### Data Backup
```bash
# Backup database
cp /opt/ospy/data/ospy.db /backup/ospy-$(date +%Y%m%d).db

# Backup configuration
cp /opt/ospy/configs/config.yaml /backup/
```

### Automated Backup
Add to crontab:
```bash
# Daily backup at 2 AM
0 2 * * * cp /opt/ospy/data/ospy.db /backup/ospy-$(date +\%Y\%m\%d).db
```
