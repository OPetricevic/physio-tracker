$ErrorActionPreference = "Stop"

# Install/update Windows service to run physio server
$root = Split-Path -Parent $MyInvocation.MyCommand.Definition
$exe = Join-Path $root "..\server.exe"
if (-not (Test-Path $exe)) {
  $exe = Join-Path $root "..\server"
}
$svcName = "PhysioTracker"
$displayName = "Physio Tracker"
$port = $env:PORT
if (-not $port) { $port = "3600" }
$db = $env:DATABASE_URL
if (-not $db) { $db = "postgres://physio:physio@localhost:5433/physio?sslmode=disable" }

$envVars = "PORT=$port", "DATABASE_URL=$db", "FRONTEND_DIR=$($root)\..\frontend\dist"

if (Get-Service -Name $svcName -ErrorAction SilentlyContinue) {
  Stop-Service $svcName -ErrorAction SilentlyContinue
  sc.exe delete $svcName | Out-Null
}

sc.exe create $svcName binPath= "`"$exe`"" start= auto DisplayName= "$displayName" | Out-Null
sc.exe description $svcName "Physio Tracker local server" | Out-Null

# Set environment variables for the service (registry)
$regPath = "HKLM:\SYSTEM\CurrentControlSet\Services\$svcName"
Set-ItemProperty -Path $regPath -Name "Environment" -Value ($envVars -join [char]0) -Type MultiString

Start-Service $svcName
Write-Host "Service $svcName installed and started on port $port"
