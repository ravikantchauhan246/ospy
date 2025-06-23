# Release script for Ospy (PowerShell version)

param(
    [Parameter(Mandatory=$true)]
    [string]$Version,
    [switch]$SkipBuild,
    [switch]$Draft
)

$ErrorActionPreference = "Stop"

Write-Host "Creating release for Ospy $Version..." -ForegroundColor Green

# Validate version format (semantic versioning)
if ($Version -notmatch '^v?\d+\.\d+\.\d+(-[a-zA-Z0-9]+)?$') {
    Write-Error "Version must follow semantic versioning (e.g., v1.0.0 or 1.0.0)"
    exit 1
}

# Ensure version starts with 'v'
if (-not $Version.StartsWith('v')) {
    $Version = "v$Version"
}

# Create releases directory
$RELEASE_DIR = "releases"
if (!(Test-Path $RELEASE_DIR)) {
    New-Item -ItemType Directory -Name $RELEASE_DIR | Out-Null
}

# Build if not skipped
if (-not $SkipBuild) {
    Write-Host "Building binaries..." -ForegroundColor Cyan
    $env:VERSION = $Version
    & powershell -ExecutionPolicy Bypass -File scripts/build.ps1 -Version $Version
    Remove-Item env:VERSION -ErrorAction SilentlyContinue
}

# Create release directory for this version
$VERSION_DIR = "$RELEASE_DIR/$Version"
if (Test-Path $VERSION_DIR) {
    Write-Host "Removing existing release directory..." -ForegroundColor Yellow
    Remove-Item $VERSION_DIR -Recurse -Force
}
New-Item -ItemType Directory -Path $VERSION_DIR | Out-Null

# Package binaries
Write-Host "Packaging binaries..." -ForegroundColor Cyan

$binaries = @(
    @{Name="ospy-linux-amd64"; Archive="ospy-$Version-linux-amd64.tar.gz"}
    @{Name="ospy-linux-arm64"; Archive="ospy-$Version-linux-arm64.tar.gz"}
    @{Name="ospy-darwin-amd64"; Archive="ospy-$Version-darwin-amd64.tar.gz"}
    @{Name="ospy-darwin-arm64"; Archive="ospy-$Version-darwin-arm64.tar.gz"}
    @{Name="ospy-windows-amd64.exe"; Archive="ospy-$Version-windows-amd64.zip"}
)

foreach ($binary in $binaries) {
    $binaryPath = "dist/$($binary.Name)"
    $archivePath = "$VERSION_DIR/$($binary.Archive)"
    
    if (Test-Path $binaryPath) {
        Write-Host "Packaging $($binary.Name)..." -ForegroundColor Yellow
        
        if ($binary.Archive.EndsWith('.zip')) {
            # Windows - create ZIP
            Compress-Archive -Path $binaryPath -DestinationPath $archivePath -Force
        } else {
            # Unix systems - create tar.gz (requires tar on Windows)
            if (Get-Command tar -ErrorAction SilentlyContinue) {
                & tar -czf $archivePath -C dist $binary.Name
            } else {
                Write-Warning "tar not found, creating zip instead of tar.gz"
                $zipPath = $archivePath -replace '\.tar\.gz$', '.zip'
                Compress-Archive -Path $binaryPath -DestinationPath $zipPath -Force
            }
        }
    } else {
        Write-Warning "Binary not found: $binaryPath"
    }
}

# Copy additional files
Write-Host "Copying additional files..." -ForegroundColor Cyan

$additionalFiles = @("README.md", "LICENSE", "configs/config.yaml")
foreach ($file in $additionalFiles) {
    if (Test-Path $file) {
        Copy-Item $file $VERSION_DIR/
    }
}

# Create checksums
Write-Host "Generating checksums..." -ForegroundColor Cyan
$checksumFile = "$VERSION_DIR/checksums.txt"
Get-ChildItem "$VERSION_DIR/*.tar.gz", "$VERSION_DIR/*.zip" | ForEach-Object {
    $hash = (Get-FileHash $_.FullName -Algorithm SHA256).Hash.ToLower()
    "$hash  $($_.Name)" | Add-Content $checksumFile
}

Write-Host "Release $Version created successfully!" -ForegroundColor Green
Write-Host "Files are in: $VERSION_DIR" -ForegroundColor Yellow

# Show release summary
Write-Host "`nRelease Summary:" -ForegroundColor Cyan
Get-ChildItem $VERSION_DIR | Format-Table Name, Length, LastWriteTime

if ($Draft) {
    Write-Host "`nDraft release created. Review files before publishing." -ForegroundColor Yellow
} else {
    Write-Host "`nRelease is ready for distribution!" -ForegroundColor Green
}
