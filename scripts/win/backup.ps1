$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $MyInvocation.MyCommand.Definition
$backupDir = Join-Path $root "..\backups"
if (-not (Test-Path $backupDir)) { New-Item -ItemType Directory -Path $backupDir | Out-Null }

$db = $env:DATABASE_URL
if (-not $db) { $db = "postgres://physio:physio@localhost:5433/physio?sslmode=disable" }

$timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
$file = Join-Path $backupDir "backup_$timestamp.sql"

Write-Host "Running pg_dump to $file"
& pg_dump $db > $file
Write-Host "Backup created at $file"
