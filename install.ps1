# idops installer for Windows
# Usage: irm https://raw.githubusercontent.com/nhh0718/idops/main/install.ps1 | iex

$ErrorActionPreference = "Stop"
$repo = "nhh0718/idops"
$binary = "idops"

# Detect architecture
$arch = if ([Environment]::Is64BitOperatingSystem) {
    if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }
} else { "amd64" }

# Get latest release
Write-Host "Fetching latest release..." -ForegroundColor Cyan
$release = Invoke-RestMethod "https://api.github.com/repos/$repo/releases/latest"
$version = $release.tag_name
Write-Host "  Version: $version" -ForegroundColor Green

# Find matching asset
$assetName = "${binary}_$($version.TrimStart('v'))_windows_${arch}.zip"
$asset = $release.assets | Where-Object { $_.name -eq $assetName }
if (-not $asset) {
    Write-Host "Error: No binary found for windows/$arch" -ForegroundColor Red
    exit 1
}

# Download
$tmpDir = Join-Path $env:TEMP "idops-install"
$zipPath = Join-Path $tmpDir $assetName
New-Item -ItemType Directory -Force -Path $tmpDir | Out-Null

Write-Host "Downloading $assetName..." -ForegroundColor Cyan
Invoke-WebRequest -Uri $asset.browser_download_url -OutFile $zipPath

# Extract
Write-Host "Extracting..." -ForegroundColor Cyan
Expand-Archive -Path $zipPath -DestinationPath $tmpDir -Force

# Install to user's local bin
$installDir = Join-Path $env:LOCALAPPDATA "idops"
New-Item -ItemType Directory -Force -Path $installDir | Out-Null

# Find binary (may be at root or inside a subdirectory)
$binaryFile = Get-ChildItem -Path $tmpDir -Recurse -Filter "$binary.exe" | Select-Object -First 1
if (-not $binaryFile) {
    Write-Host "Error: $binary.exe not found in archive" -ForegroundColor Red
    exit 1
}
Copy-Item $binaryFile.FullName (Join-Path $installDir "$binary.exe") -Force

# Find and install dashboard (dashboard-dist from release, or dashboard from source)
$dashboardSrc = Get-ChildItem -Path $tmpDir -Recurse -Directory -Filter "dashboard-dist" | Select-Object -First 1
if (-not $dashboardSrc) {
    $dashboardSrc = Get-ChildItem -Path $tmpDir -Recurse -Directory -Filter "dashboard" | Select-Object -First 1
}
if ($dashboardSrc) {
    $dashboardDest = Join-Path $installDir "dashboard-dist"
    if (Test-Path $dashboardDest) { Remove-Item -Recurse -Force $dashboardDest }
    Copy-Item $dashboardSrc.FullName $dashboardDest -Recurse -Force
    Write-Host "  Dashboard installed" -ForegroundColor Green
} else {
    Write-Host "  Warning: Dashboard not found in archive" -ForegroundColor Yellow
}

# Add to PATH if not already there
$userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($userPath -notlike "*$installDir*") {
    [Environment]::SetEnvironmentVariable("PATH", "$userPath;$installDir", "User")
    Write-Host "  Added $installDir to PATH" -ForegroundColor Yellow
    Write-Host "  Restart your terminal for PATH changes to take effect" -ForegroundColor Yellow
}

# Cleanup
Remove-Item -Recurse -Force $tmpDir

Write-Host ""
Write-Host "idops $version installed to $installDir\$binary.exe" -ForegroundColor Green
Write-Host "Run: idops --help" -ForegroundColor Cyan
