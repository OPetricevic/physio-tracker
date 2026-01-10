; Inno Setup script for Physio Tracker (Windows x64)
; This script lives in: scripts\win\physio-x64.iss
; When copied into the bundle, it should be located at: release\physio-bundle\scripts\win\
; Bundle root is two levels up from this file.

[Setup]
AppName=Physio Tracker
AppVersion=1.0.0
DefaultDirName={pf}\PhysioTracker
DisableDirPage=yes
DisableProgramGroupPage=yes
OutputDir=.
OutputBaseFilename=physio-tracker-setup-x64
Compression=lzma
SolidCompression=yes
WizardStyle=modern
PrivilegesRequired=admin

[Files]
; App binary (bundle root)
Source: "..\..\server.exe"; DestDir: "{app}"; Flags: ignoreversion

; Frontend + assets (bundle root)
Source: "..\..\frontend\dist\*"; DestDir: "{app}\frontend\dist"; Flags: recursesubdirs ignoreversion
Source: "..\..\assets\fonts\*"; DestDir: "{app}\assets\fonts"; Flags: recursesubdirs ignoreversion

; uploads may be empty, so don't fail build
Source: "..\..\uploads\*"; DestDir: "{app}\uploads"; Flags: recursesubdirs ignoreversion skipifsourcedoesntexist

; Root script (bundle root)
Source: "..\..\start_windows.ps1"; DestDir: "{app}"; Flags: ignoreversion

; Windows scripts (these are in the SAME folder as this .iss: scripts\win\)
Source: "install_postgres.ps1"; DestDir: "{app}\scripts"; Flags: ignoreversion
Source: "service_install.ps1"; DestDir: "{app}\scripts"; Flags: ignoreversion
Source: "backup.ps1"; DestDir: "{app}\scripts"; Flags: ignoreversion
Source: "restore.ps1"; DestDir: "{app}\scripts"; Flags: ignoreversion

[Run]
; Install PostgreSQL if missing, then install service (uses bundled server.exe)
Filename: "powershell.exe"; Parameters: "-ExecutionPolicy Bypass -File ""{app}\scripts\install_postgres.ps1"""; StatusMsg: "Setting up PostgreSQL..."
Filename: "powershell.exe"; Parameters: "-ExecutionPolicy Bypass -File ""{app}\scripts\service_install.ps1"""; StatusMsg: "Installing Windows service..."

[Icons]
Name: "{commonprograms}\Physio Tracker"; Filename: "http://localhost:3600"
Name: "{commondesktop}\Physio Tracker"; Filename: "http://localhost:3600"

[UninstallRun]
; Remove service on uninstall (escaped braces)
Filename: "powershell.exe"; Parameters: "-ExecutionPolicy Bypass -Command ""if (Get-Service -Name PhysioTracker -ErrorAction SilentlyContinue) {{ Stop-Service PhysioTracker -ErrorAction SilentlyContinue; sc.exe delete PhysioTracker }}"""
