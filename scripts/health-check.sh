#!/bin/bash
# Health Check Script for Ospy Configuration
# This script validates your Ospy configuration and environment

set -e

CONFIG_FILE="${1:-configs/config.yaml}"
OSPY_BINARY="${2:-./ospy}"

echo "🔍 Ospy Configuration Health Check"
echo "=================================="
echo

# Check if config file exists
echo "📋 Checking configuration file..."
if [ ! -f "$CONFIG_FILE" ]; then
    echo "❌ Configuration file not found: $CONFIG_FILE"
    echo "   Create a config file or specify the path as the first argument"
    exit 1
fi
echo "✅ Configuration file found: $CONFIG_FILE"
echo

# Check if Ospy binary exists
echo "🔧 Checking Ospy binary..."
if [ ! -f "$OSPY_BINARY" ] && ! command -v ospy &> /dev/null; then
    echo "❌ Ospy binary not found: $OSPY_BINARY"
    echo "   Build Ospy with: go build -o ospy ./cmd/ospy"
    exit 1
fi
echo "✅ Ospy binary found"
echo

# Check if data directory will be created
echo "📁 Checking data directory..."
DATA_DIR=$(dirname "$(grep -A 5 "storage:" "$CONFIG_FILE" | grep "path:" | cut -d'"' -f2 | head -1)")
if [ "$DATA_DIR" != "." ] && [ ! -d "$DATA_DIR" ]; then
    echo "⚠️  Data directory will be created: $DATA_DIR"
    mkdir -p "$DATA_DIR" 2>/dev/null || echo "❌ Cannot create data directory: $DATA_DIR"
else
    echo "✅ Data directory accessible"
fi
echo

# Check if log directory will be created
echo "📝 Checking log directory..."
if grep -q "file:" "$CONFIG_FILE"; then
    LOG_FILE=$(grep "file:" "$CONFIG_FILE" | cut -d'"' -f2 | head -1)
    LOG_DIR=$(dirname "$LOG_FILE")
    if [ "$LOG_DIR" != "." ] && [ ! -d "$LOG_DIR" ]; then
        echo "⚠️  Log directory will be created: $LOG_DIR"
        mkdir -p "$LOG_DIR" 2>/dev/null || echo "❌ Cannot create log directory: $LOG_DIR"
    else
        echo "✅ Log directory accessible"
    fi
else
    echo "ℹ️  No log file configured"
fi
echo

# Check environment variables for notifications
echo "🔔 Checking notification configuration..."
if grep -q "enabled: true" "$CONFIG_FILE"; then
    if grep -A 10 "email:" "$CONFIG_FILE" | grep -q "enabled: true"; then
        echo "📧 Email notifications enabled"
        if [ -z "$SMTP_USERNAME" ]; then
            echo "⚠️  SMTP_USERNAME environment variable not set"
        else
            echo "✅ SMTP_USERNAME environment variable set"
        fi
        if [ -z "$SMTP_PASSWORD" ]; then
            echo "⚠️  SMTP_PASSWORD environment variable not set"
        else
            echo "✅ SMTP_PASSWORD environment variable set"
        fi
    fi
    
    if grep -A 10 "telegram:" "$CONFIG_FILE" | grep -q "enabled: true"; then
        echo "📱 Telegram notifications enabled"
        if [ -z "$TELEGRAM_BOT_TOKEN" ]; then
            echo "⚠️  TELEGRAM_BOT_TOKEN environment variable not set"
        else
            echo "✅ TELEGRAM_BOT_TOKEN environment variable set"
        fi
    fi
else
    echo "ℹ️  No notifications enabled"
fi
echo

# Test configuration syntax by running Ospy version command
echo "🧪 Testing Ospy binary..."
if [ -f "$OSPY_BINARY" ]; then
    VERSION_OUTPUT=$("$OSPY_BINARY" -version 2>&1) || {
        echo "❌ Failed to run Ospy binary"
        exit 1
    }
    echo "✅ Ospy version: $(echo "$VERSION_OUTPUT" | head -1)"
else
    VERSION_OUTPUT=$(ospy -version 2>&1) || {
        echo "❌ Failed to run Ospy command"
        exit 1
    }
    echo "✅ Ospy version: $(echo "$VERSION_OUTPUT" | head -1)"
fi
echo

# Check websites accessibility (basic connectivity test)
echo "🌐 Testing website connectivity..."
WEBSITES=$(grep -A 20 "websites:" "$CONFIG_FILE" | grep "url:" | cut -d'"' -f2)
for url in $WEBSITES; do
    if command -v curl &> /dev/null; then
        if curl -s --head "$url" --connect-timeout 10 > /dev/null; then
            echo "✅ $url - accessible"
        else
            echo "⚠️  $url - not accessible (may be temporary)"
        fi
    elif command -v wget &> /dev/null; then
        if wget --spider --timeout=10 "$url" > /dev/null 2>&1; then
            echo "✅ $url - accessible"
        else
            echo "⚠️  $url - not accessible (may be temporary)"
        fi
    else
        echo "ℹ️  $url - cannot test (curl/wget not available)"
    fi
done
echo

echo "🎉 Health check completed!"
echo
echo "💡 Tips:"
echo "   - Run 'ospy -help' to see all available options"
echo "   - Use 'ospy -config $CONFIG_FILE' to start monitoring"
echo "   - Check the examples/ directory for more configuration templates"
echo "   - Monitor the log files for any runtime issues"
