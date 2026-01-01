$ErrorActionPreference = "Stop"

param(
  [string]$PgVersion = "15",
  [string]$ServiceName = "physio-postgres",
  [string]$SuperPassword = "physio",
  [string]$DbUser = "physio",
  [string]$DbPassword = "physio",
  [int]$Port = 5432
)

# Installs PostgreSQL silently if it is not already present, then ensures the
# physio database/user exist. Designed for the Physio Tracker installer.

function Test-ServiceExists {
  param([string]$Name)
  return [bool](Get-Service -Name $Name -ErrorAction SilentlyContinue)
}

$installDir = Join-Path $env:ProgramFiles "PhysioPostgres"
$dataDir = Join-Path ${env:ProgramData} "PhysioTracker\postgres"
$binDir = Join-Path $installDir "bin"
$psql = Join-Path $binDir "psql.exe"

if (Test-ServiceExists -Name $ServiceName) {
  Write-Host "PostgreSQL service '$ServiceName' already present. Skipping install."
} else {
  $installerUrl = "https://get.enterprisedb.com/postgresql/postgresql-$PgVersion.8-1-windows-x64.exe"
  $tmpDir = New-Item -ItemType Directory -Path ([System.IO.Path]::Combine([System.IO.Path]::GetTempPath(), "physio-pg")) -Force
  $installerPath = Join-Path $tmpDir "postgres_installer.exe"

  Write-Host "Downloading PostgreSQL $PgVersion from $installerUrl ..."
  Invoke-WebRequest -Uri $installerUrl -OutFile $installerPath

  if (-not (Test-Path $installerPath)) {
    throw "Failed to download PostgreSQL installer."
  }

  if (-not (Test-Path $dataDir)) {
    New-Item -ItemType Directory -Path $dataDir | Out-Null
  }

  $args = @(
    "--mode", "unattended",
    "--unattendedmodeui", "minimal",
    "--prefix", "`"$installDir`"",
    "--datadir", "`"$dataDir`"",
    "--superpassword", "`"$SuperPassword`"",
    "--servicename", "`"$ServiceName`"",
    "--serverport", "$Port"
  )

  Write-Host "Installing PostgreSQL to $installDir (service: $ServiceName, port: $Port)..."
  $process = Start-Process -FilePath $installerPath -ArgumentList $args -Wait -PassThru
  if ($process.ExitCode -ne 0) {
    throw "PostgreSQL installer failed with exit code $($process.ExitCode)"
  }
}

if (-not (Test-Path $psql)) {
  throw "psql not found at $psql. Installation may have failed."
}

$env:PATH = "$binDir;$env:PATH"

Write-Host "Ensuring database user '$DbUser' exists..."
$userExists = & $psql -U postgres -h localhost -p $Port -tc "SELECT 1 FROM pg_roles WHERE rolname='$DbUser';" | Select-String "1" -Quiet
if (-not $userExists) {
  & $psql -U postgres -h localhost -p $Port -c "CREATE USER $DbUser WITH PASSWORD '$DbPassword';"
}

Write-Host "Ensuring database 'physio' exists..."
$dbExists = & $psql -U postgres -h localhost -p $Port -tc "SELECT 1 FROM pg_database WHERE datname='physio';" | Select-String "1" -Quiet
if (-not $dbExists) {
  & $psql -U postgres -h localhost -p $Port -c "CREATE DATABASE physio OWNER $DbUser;"
}

Write-Host "PostgreSQL ready (service '$ServiceName', port $Port)."
