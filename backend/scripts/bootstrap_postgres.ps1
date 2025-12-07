# PowerShell bootstrap for Windows
# Creates an app role and database if missing.
# Defaults can be overridden via environment variables:
#   DB_BOOTSTRAP_URL  - superuser connection string (default: postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable)
#   DB_NAME           - database name (default: physio)
#   DB_APP_USER       - app role/user (default: physio_app)
#   DB_APP_PASS       - app user password (default: physio_app_pass)

$bootstrapUrl = $env:DB_BOOTSTRAP_URL
if ([string]::IsNullOrWhiteSpace($bootstrapUrl)) { $bootstrapUrl = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" }

$dbName = $env:DB_NAME
if ([string]::IsNullOrWhiteSpace($dbName)) { $dbName = "physio" }

$appUser = $env:DB_APP_USER
if ([string]::IsNullOrWhiteSpace($appUser)) { $appUser = "physio_app" }

$appPass = $env:DB_APP_PASS
if ([string]::IsNullOrWhiteSpace($appPass)) { $appPass = "physio_app_pass" }

$sqlTemplate = @'
DO $$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = '%APPUSER%') THEN
      CREATE ROLE %APPUSER% LOGIN PASSWORD '%APPPASS%';
   END IF;
END
$$;

DO $$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_database WHERE datname = '%DBNAME%') THEN
      CREATE DATABASE %DBNAME% OWNER %APPUSER%;
   END IF;
END
$$;
'@

$sql = $sqlTemplate.
    Replace("%APPUSER%", $appUser).
    Replace("%APPPASS%", $appPass).
    Replace("%DBNAME%", $dbName)

Write-Host "Using bootstrap URL: $bootstrapUrl"
Write-Host "Ensuring role '$appUser' and database '$dbName' exist..."

psql $bootstrapUrl -v ON_ERROR_STOP=1 -c $sql

Write-Host "Bootstrap complete."
Write-Host "Set DATABASE_URL=postgres://$appUser:$appPass@localhost:5432/$dbName?sslmode=disable"
