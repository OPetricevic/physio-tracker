$ErrorActionPreference = "Stop"

# This script starts the local server and serves the built frontend.
# Defaults can be overridden via environment variables.

$root = Split-Path -Parent $MyInvocation.MyCommand.Definition

if (-not $env:FRONTEND_DIR) { $env:FRONTEND_DIR = Join-Path $root "frontend/dist" }
if (-not $env:PORT) { $env:PORT = "3600" }
if (-not $env:DATABASE_URL) { $env:DATABASE_URL = "postgres://physio:physio@localhost:5433/physio?sslmode=disable" }

$serverPath = Join-Path $root "server.exe"
if (-not (Test-Path $serverPath)) {
  $serverPath = Join-Path $root "server"
}

if (-not (Test-Path $serverPath)) {
  Write-Error "server binary not found in $root. Build/package first."
}

Write-Host "Starting server on :$($env:PORT) (DATABASE_URL=$($env:DATABASE_URL))"
& $serverPath
