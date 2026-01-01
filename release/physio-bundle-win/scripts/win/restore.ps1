$ErrorActionPreference = "Stop"

param(
  [string]$File
)

if (-not $File -or -not (Test-Path $File)) {
  Write-Error "Provide a valid backup file path: .\restore.ps1 -File path\to\backup.sql"
}

$db = $env:DATABASE_URL
if (-not $db) { $db = "postgres://physio:physio@localhost:5433/physio?sslmode=disable" }

Write-Host "Restoring from $File"
& psql $db < $File
Write-Host "Restore complete."
