# PowerShell backup script for Windows
# Creates a timestamped pg_dump (-Fc) and optionally prunes old backups.
# Env vars:
#   DATABASE_URL   - required, connection string to the target DB.
#   BACKUP_DIR     - optional, defaults to ../backups relative to this script.
#   RETAIN_DAYS    - optional, prune backups older than this (default: 14 days). Set to 0 to disable pruning.

if (-not $env:DATABASE_URL) {
    Write-Error "DATABASE_URL is required"
    exit 1
}

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$DefaultBackupDir = Join-Path (Split-Path -Parent $ScriptDir) "backups"
$BackupDir = if ($env:BACKUP_DIR) { $env:BACKUP_DIR } else { $DefaultBackupDir }

if (-not (Test-Path $BackupDir)) {
    New-Item -ItemType Directory -Force -Path $BackupDir | Out-Null
}

$timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
$backupFile = Join-Path $BackupDir ("physio_" + $timestamp + ".dump")

Write-Host "Creating backup: $backupFile"
pg_dump $env:DATABASE_URL -Fc -f $backupFile
Write-Host "Backup complete."

$retainDays = 14
if ($env:RETAIN_DAYS) {
    if ([int]::TryParse($env:RETAIN_DAYS, [ref]$null)) {
        $retainDays = [int]$env:RETAIN_DAYS
    }
}

if ($retainDays -gt 0) {
    $cutoff = (Get-Date).AddDays(-$retainDays)
    $oldBackups = Get-ChildItem -Path $BackupDir -Filter "physio_*.dump" | Where-Object { $_.LastWriteTime -lt $cutoff }
    foreach ($file in $oldBackups) {
        Write-Host "Pruning old backup: $($file.FullName)"
        Remove-Item -Force $file.FullName
    }
}

Write-Host "Done. Latest backup: $backupFile"
