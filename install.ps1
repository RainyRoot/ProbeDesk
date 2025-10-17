# Requires PowerShell 5+
# Self-elevate if not running as Administrator
if (-not ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
    Write-Output "Restarting as Administrator..."
    $psi = New-Object System.Diagnostics.ProcessStartInfo
    $psi.FileName = "powershell.exe"
    $psi.Arguments = "-ExecutionPolicy Bypass -File `"$($MyInvocation.MyCommand.Definition)`""
    $psi.Verb = "runas"
    try {
        [System.Diagnostics.Process]::Start($psi) | Out-Null
    } catch {
        Write-Output "Installation aborted by user."
    }
    exit
}

Write-Output "Running with Administrator privileges."

# --- Installation variables ---
$exe = ".\probedesk.exe"
$targetDir = "C:\Tools"
$target = Join-Path $targetDir "probedesk.exe"

# --- Create target directory ---
if (!(Test-Path $targetDir)) {
    New-Item -ItemType Directory -Path $targetDir | Out-Null
}

# --- Copy binary ---
Copy-Item $exe $target -Force
Write-Output "ProbeDesk installed in $targetDir"

# --- PATH update (User scope) ---
$oldPath = [Environment]::GetEnvironmentVariable("PATH", [EnvironmentVariableTarget]::User)
if ($oldPath -notlike "*$targetDir*") {
    [Environment]::SetEnvironmentVariable("PATH", $oldPath + ";$targetDir", [EnvironmentVariableTarget]::User)
    Write-Output "$targetDir added to PATH. Restart your terminal to use 'probedesk' globally."
}

# --- Autocomplete installation ---
$profilePath = Join-Path $env:USERPROFILE "Documents\PowerShell\Microsoft.PowerShell_profile.ps1"

if (!(Test-Path $profilePath)) {
    New-Item -ItemType File -Path $profilePath -Force | Out-Null
}

$autocomplete = @"
\$flags = @("system","ipconfig","netuse","products","vpn","services","users","usb","trace","remote","report","flushdns","winget-update","scan-health","check-health","restore-health","autocomplete-install")

Register-ArgumentCompleter -CommandName "probedesk" -ScriptBlock {
    param(\$commandName, \$parameterName, \$wordToComplete, \$commandAst, \$fakeBoundParameter)
    if (\$commandAst.CommandElements.Count -gt 1 -and \$commandAst.CommandElements[1].Value -eq "win") {
        \$flags | Where-Object { \$_ -like "\$wordToComplete*" } |
            ForEach-Object { 
                [System.Management.Automation.CompletionResult]::new(\$_, \$_, 'ParameterValue', \$_) 
            }
    }
}
"@

if (-not (Get-Content $profilePath | Select-String "Register-ArgumentCompleter -CommandName `"probedesk`"")) {
    Add-Content -Path $profilePath -Value $autocomplete
    Write-Output "Autocomplete for ProbeDesk installed. Restart PowerShell to activate."
} else {
    Write-Output "Autocomplete already installed."
}

Write-Output "✅ Installation completed successfully."
Pause
