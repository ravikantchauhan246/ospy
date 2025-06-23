# Build script for Ospy (PowerShell version)

param(
    [string]$Version = "dev"
)

$ErrorActionPreference = "Stop"

$BUILD_TIME = (Get-Date).ToUniversalTime().ToString("yyyy-MM-dd_HH:mm:ss")
try {
    $GIT_COMMIT = (git rev-parse --short HEAD 2>$null)
} catch {
    $GIT_COMMIT = "unknown"
}

$LDFLAGS = "-X main.Version=$Version -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT"

Write-Host "Building Ospy $Version..." -ForegroundColor Green
Write-Host "Build time: $BUILD_TIME" -ForegroundColor Yellow
Write-Host "Git commit: $GIT_COMMIT" -ForegroundColor Yellow

# Create dist directory
if (!(Test-Path "dist")) {
    New-Item -ItemType Directory -Name "dist" | Out-Null
}

# Build for current platform (Windows)
Write-Host "Building for current platform..." -ForegroundColor Cyan
go build -ldflags $LDFLAGS -o dist/ospy.exe ./cmd/ospy

# Build for multiple platforms
Write-Host "Cross-compiling..." -ForegroundColor Cyan

# Linux AMD64
Write-Host "Building for Linux AMD64..."
$env:GOOS = "linux"; $env:GOARCH = "amd64"
go build -ldflags $LDFLAGS -o dist/ospy-linux-amd64 ./cmd/ospy

# Linux ARM64
Write-Host "Building for Linux ARM64..."
$env:GOOS = "linux"; $env:GOARCH = "arm64"
go build -ldflags $LDFLAGS -o dist/ospy-linux-arm64 ./cmd/ospy

# macOS AMD64
Write-Host "Building for macOS AMD64..."
$env:GOOS = "darwin"; $env:GOARCH = "amd64"
go build -ldflags $LDFLAGS -o dist/ospy-darwin-amd64 ./cmd/ospy

# macOS ARM64 (Apple Silicon)
Write-Host "Building for macOS ARM64..."
$env:GOOS = "darwin"; $env:GOARCH = "arm64"
go build -ldflags $LDFLAGS -o dist/ospy-darwin-arm64 ./cmd/ospy

# Windows AMD64
Write-Host "Building for Windows AMD64..."
$env:GOOS = "windows"; $env:GOARCH = "amd64"
go build -ldflags $LDFLAGS -o dist/ospy-windows-amd64.exe ./cmd/ospy

# Reset environment variables
Remove-Item env:GOOS -ErrorAction SilentlyContinue
Remove-Item env:GOARCH -ErrorAction SilentlyContinue

Write-Host "Build complete! Binaries are in dist/" -ForegroundColor Green
Get-ChildItem dist/ | Format-Table Name, Length, LastWriteTime
