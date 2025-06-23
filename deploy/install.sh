#!/bin/bash

# Installation script for Ospy

set -e

# Configuration
INSTALL_DIR="/opt/ospy"
SERVICE_FILE="/etc/systemd/system/ospy.service"
USER="ospy"
GROUP="ospy"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

echo_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

echo_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo_error "This script must be run as root"
   exit 1
fi

echo_info "Installing Ospy Website Monitor..."

# Create user and group
if ! id "$USER" &>/dev/null; then
    echo_info "Creating user: $USER"
    useradd --system --home-dir $INSTALL_DIR --shell /bin/false $USER
else
    echo_info "User $USER already exists"
fi

# Create installation directory
echo_info "Creating installation directory: $INSTALL_DIR"
mkdir -p $INSTALL_DIR/{data,logs,configs}
chown -R $USER:$GROUP $INSTALL_DIR

# Copy binary
if [[ -f "ospy" ]]; then
    echo_info "Installing binary to $INSTALL_DIR"
    cp ospy $INSTALL_DIR/
    chmod +x $INSTALL_DIR/ospy
    chown $USER:$GROUP $INSTALL_DIR/ospy
else
    echo_error "Binary 'ospy' not found in current directory"
    exit 1
fi

# Copy configuration
if [[ -f "config.yaml" ]]; then
    echo_info "Installing configuration"
    cp config.yaml $INSTALL_DIR/configs/
    chown $USER:$GROUP $INSTALL_DIR/configs/config.yaml
elif [[ -f "configs/config.yaml" ]]; then
    cp configs/config.yaml $INSTALL_DIR/configs/
    chown $USER:$GROUP $INSTALL_DIR/configs/config.yaml
else
    echo_warn "No configuration file found, creating default"
    cat > $INSTALL_DIR/configs/config.yaml << EOF
monitoring:
  interval: 5m
  timeout: 30s
  retries: 3
  workers: 10

websites:
  - name: "Example"
    url: "https://example.com"
    expected_status: 200

notifications:
  email:
    enabled: false
  telegram:
    enabled: false

storage:
  type: "sqlite"
  path: "./data/ospy.db"

web:
  enabled: true
  host: "0.0.0.0"
  port: 8080

logging:
  level: "info"
  file: "./logs/ospy.log"
EOF
    chown $USER:$GROUP $INSTALL_DIR/configs/config.yaml
fi

# Install systemd service
if [[ -f "deploy/ospy.service" ]]; then
    echo_info "Installing systemd service"
    cp deploy/ospy.service $SERVICE_FILE
    systemctl daemon-reload
    systemctl enable ospy
else
    echo_warn "Service file not found, creating default"
    cat > $SERVICE_FILE << EOF
[Unit]
Description=Ospy Website Monitor
After=network.target

[Service]
Type=simple
User=$USER
Group=$GROUP
WorkingDirectory=$INSTALL_DIR
ExecStart=$INSTALL_DIR/ospy -config $INSTALL_DIR/configs/config.yaml
Restart=always
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF
    systemctl daemon-reload
    systemctl enable ospy
fi

echo_info "Installation completed successfully!"
echo_info "Configuration file: $INSTALL_DIR/configs/config.yaml"
echo_info "Log file: $INSTALL_DIR/logs/ospy.log"
echo_info "Data directory: $INSTALL_DIR/data"
echo ""
echo_info "To start Ospy:"
echo "  sudo systemctl start ospy"
echo ""
echo_info "To check status:"
echo "  sudo systemctl status ospy"
echo ""
echo_info "To view logs:"
echo "  sudo journalctl -u ospy -f"
echo ""
echo_warn "Don't forget to configure your monitoring settings in:"
echo "  $INSTALL_DIR/configs/config.yaml"
