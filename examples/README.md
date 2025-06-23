# Configuration Examples

This directory contains example configurations for different use cases and environments.

## Available Examples

### `minimal-config.yaml`
A bare-bones configuration perfect for getting started quickly. Includes basic monitoring without notifications.

**Use case:** Quick testing, personal projects, or simple monitoring setups.

### `development-config.yaml`
Optimized for development environments with verbose logging and local service monitoring.

**Features:**
- Quick check intervals for fast feedback
- Debug-level logging
- Local development server monitoring
- Notifications typically disabled

**Use case:** Development and testing environments.

### `production-config.yaml`
Enterprise-ready configuration optimized for production monitoring.

**Features:**
- Longer intervals for stability
- Multiple notification channels
- Extended data retention
- Comprehensive monitoring coverage
- Production-level logging

**Use case:** Production environments, SLA monitoring, enterprise deployments.

## Using These Examples

1. Copy the example that best fits your needs:
   ```bash
   cp examples/production-config.yaml configs/config.yaml
   ```

2. Edit the configuration file to match your specific requirements:
   - Update website URLs
   - Configure notification settings
   - Set appropriate file paths
   - Adjust monitoring intervals

3. Set up environment variables for sensitive data:
   ```bash
   export SMTP_USERNAME="your-email@company.com"
   export SMTP_PASSWORD="your-app-password"
   export TELEGRAM_BOT_TOKEN="your-bot-token"
   ```

4. Run Ospy with your configuration:
   ```bash
   ./ospy -config configs/config.yaml
   ```

## Configuration Tips

- **Monitoring Intervals**: Balance between responsiveness and system load
- **Timeouts**: Set timeouts longer than your slowest expected response time
- **Workers**: Use more workers for monitoring many websites simultaneously
- **Retention**: Adjust retention days based on your storage capacity and compliance needs
- **Logging**: Use debug level for development, info for general use, warn for production

## Security Best Practices

- Never commit sensitive credentials to version control
- Use environment variables for passwords and tokens
- Restrict file permissions on configuration files containing sensitive data
- Use HTTPS URLs for all monitored websites when possible
- Consider using dedicated service accounts for SMTP and Telegram notifications
