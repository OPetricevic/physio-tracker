# PowerShell migration runner for Windows
# Requires: DATABASE_URL environment variable set to the target DB connection string.

if (-not $env:DATABASE_URL) {
    Write-Error "DATABASE_URL is required"
    exit 1
}

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$RootDir = Split-Path -Parent $ScriptDir

psql $env:DATABASE_URL -v ON_ERROR_STOP=1 -f (Join-Path $RootDir "migrations/0001_init.sql")
psql $env:DATABASE_URL -v ON_ERROR_STOP=1 -f (Join-Path $RootDir "migrations/0002_seed_admin.sql")

Write-Host "Migrations applied."
