; Inno Setup script for Physio Tracker (Windows)
; This assumes you already built the bundle via `make package` into release/physio-bundle
; Adjust paths as needed before building the installer.

[Setup]
AppName=Physio Tracker
AppVersion=1.0.0
DefaultDirName={pf}\PhysioTracker
DisableDirPage=yes
DisableProgramGroupPage=yes
OutputDir=.
OutputBaseFilename=physio-tracker-setup
Compression=lzma
SolidCompression=yes
WizardStyle=modern

[Files]
Source: "..\release\physio-bundle\server.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\release\physio-bundle\server"; DestDir: "{app}"; Flags: ignoreversion skipifsourcedoesntexist
Source: "..\release\physio-bundle\frontend\dist\*"; DestDir: "{app}\frontend\dist"; Flags: recursesubdirs ignoreversion
Source: "..\release\physio-bundle\assets\fonts\*"; DestDir: "{app}\assets\fonts"; Flags: recursesubdirs ignoreversion
Source: "..\release\physio-bundle\uploads\*"; DestDir: "{app}\uploads"; Flags: recursesubdirs ignoreversion
Source: "..\scripts\start_windows.ps1"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\scripts\win\service_install.ps1"; DestDir: "{app}\scripts"; Flags: ignoreversion
Source: "..\scripts\win\backup.ps1"; DestDir: "{app}\scripts"; Flags: ignoreversion
Source: "..\scripts\win\restore.ps1"; DestDir: "{app}\scripts"; Flags: ignoreversion

[Run]
; Install service (will use bundled server.exe)
Filename: "powershell.exe"; Parameters: "-ExecutionPolicy Bypass -File ""{app}\scripts\service_install.ps1"""; StatusMsg: "Installing Windows service..."

[Icons]
Name: "{commonprograms}\Physio Tracker"; Filename: "http://localhost:3600"
Name: "{commondesktop}\Physio Tracker"; Filename: "http://localhost:3600"

[UninstallRun]
; Remove service on uninstall
Filename: "powershell.exe"; Parameters: "-ExecutionPolicy Bypass -Command ""if (Get-Service -Name PhysioTracker -ErrorAction SilentlyContinue) { Stop-Service PhysioTracker -ErrorAction SilentlyContinue; sc.exe delete PhysioTracker }"""

