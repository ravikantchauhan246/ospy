# Health Check Script for Ospy Configuration (PowerShell)
# This script validates your Ospy configuration and environment

param(
    [string]$ConfigFile = "configs\config.yaml",
    [string]$OspyBinary = ".\ospy.exe"
)

Write-Host "🔍 Ospy Configuration Health Check" -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan
Write-Host

# Check if config file exists
Write-Host "📋 Checking configuration file..." -ForegroundColor Yellow
if (-not (Test-Path $ConfigFile)) {
    Write-Host "❌ Configuration file not found: $ConfigFile" -ForegroundColor Red
    Write-Host "   Create a config file or specify the path as the first argument" -ForegroundColor Red
    exit 1
}
Write-Host "✅ Configuration file found: $ConfigFile" -ForegroundColor Green
Write-Host

# Check if Ospy binary exists
Write-Host "🔧 Checking Ospy binary..." -ForegroundColor Yellow
if (-not (Test-Path $OspyBinary)) {
    Write-Host "❌ Ospy binary not found: $OspyBinary" -ForegroundColor Red
    Write-Host "   Build Ospy with: go build -o ospy.exe ./cmd/ospy" -ForegroundColor Red
    exit 1
}
Write-Host "✅ Ospy binary found" -ForegroundColor Green
Write-Host

# Check if data directory will be created
Write-Host "📁 Checking data directory..." -ForegroundColor Yellow
$configContent = Get-Content $ConfigFile -Raw
if ($configContent -match 'path:\s*"([^"]+)"') {
    $dataPath = $Matches[1]
    $dataDir = Split-Path $dataPath -Parent
    if ($dataDir -and $dataDir -ne "." -and -not (Test-Path $dataDir)) {
        Write-Host "⚠️  Data directory will be created: $dataDir" -ForegroundColor Yellow
        try {
            New-Item -ItemType Directory -Path $dataDir -Force | Out-Null
            Write-Host "✅ Data directory created" -ForegroundColor Green
        } catch {
            Write-Host "❌ Cannot create data directory: $dataDir" -ForegroundColor Red
        }
    } else {
        Write-Host "✅ Data directory accessible" -ForegroundColor Green
    }
}
Write-Host

# Check if log directory will be created
Write-Host "📝 Checking log directory..." -ForegroundColor Yellow
if ($configContent -match 'file:\s*"([^"]+)"') {
    $logFile = $Matches[1]
    $logDir = Split-Path $logFile -Parent
    if ($logDir -and $logDir -ne "." -and -not (Test-Path $logDir)) {
        Write-Host "⚠️  Log directory will be created: $logDir" -ForegroundColor Yellow
        try {
            New-Item -ItemType Directory -Path $logDir -Force | Out-Null
            Write-Host "✅ Log directory created" -ForegroundColor Green
        } catch {
            Write-Host "❌ Cannot create log directory: $logDir" -ForegroundColor Red
        }
    } else {
        Write-Host "✅ Log directory accessible" -ForegroundColor Green
    }
} else {
    Write-Host "ℹ️  No log file configured" -ForegroundColor Cyan
}
Write-Host

# Check environment variables for notifications
Write-Host "🔔 Checking notification configuration..." -ForegroundColor Yellow
if ($configContent -match 'enabled:\s*true') {
    if ($configContent -match '(?s)email:.*?enabled:\s*true') {
        Write-Host "📧 Email notifications enabled" -ForegroundColor Cyan
        if (-not $env:SMTP_USERNAME) {
            Write-Host "⚠️  SMTP_USERNAME environment variable not set" -ForegroundColor Yellow
        } else {
            Write-Host "✅ SMTP_USERNAME environment variable set" -ForegroundColor Green
        }
        if (-not $env:SMTP_PASSWORD) {
            Write-Host "⚠️  SMTP_PASSWORD environment variable not set" -ForegroundColor Yellow
        } else {
            Write-Host "✅ SMTP_PASSWORD environment variable set" -ForegroundColor Green
        }
    }
    
    if ($configContent -match '(?s)telegram:.*?enabled:\s*true') {
        Write-Host "📱 Telegram notifications enabled" -ForegroundColor Cyan
        if (-not $env:TELEGRAM_BOT_TOKEN) {
            Write-Host "⚠️  TELEGRAM_BOT_TOKEN environment variable not set" -ForegroundColor Yellow
        } else {
            Write-Host "✅ TELEGRAM_BOT_TOKEN environment variable set" -ForegroundColor Green
        }
    }
} else {
    Write-Host "ℹ️  No notifications enabled" -ForegroundColor Cyan
}
Write-Host

# Test configuration syntax by running Ospy version command
Write-Host "🧪 Testing Ospy binary..." -ForegroundColor Yellow
try {
    $versionOutput = & $OspyBinary -version 2>&1
    Write-Host "✅ Ospy version: $($versionOutput[0])" -ForegroundColor Green
} catch {
    Write-Host "❌ Failed to run Ospy binary" -ForegroundColor Red
    exit 1
}
Write-Host

# Check websites accessibility (basic connectivity test)
Write-Host "🌐 Testing website connectivity..." -ForegroundColor Yellow
$urls = $configContent | Select-String 'url:\s*"([^"]+)"' -AllMatches | ForEach-Object { $_.Matches } | ForEach-Object { $_.Groups[1].Value }
foreach ($url in $urls) {
    try {
        $response = Invoke-WebRequest -Uri $url -Method Head -TimeoutSec 10 -UseBasicParsing -ErrorAction Stop
        Write-Host "✅ $url - accessible" -ForegroundColor Green
    } catch {
        Write-Host "⚠️  $url - not accessible (may be temporary)" -ForegroundColor Yellow
    }
}
Write-Host

Write-Host "🎉 Health check completed!" -ForegroundColor Green
Write-Host
Write-Host "💡 Tips:" -ForegroundColor Cyan
Write-Host "   - Run 'ospy.exe -help' to see all available options" -ForegroundColor White
Write-Host "   - Use 'ospy.exe -config $ConfigFile' to start monitoring" -ForegroundColor White
Write-Host "   - Check the examples\ directory for more configuration templates" -ForegroundColor White
Write-Host "   - Monitor the log files for any runtime issues" -ForegroundColor White
