$exe = ".\probedesk.exe"
$target = "C:\Tools\probedesk.exe"

if (!(Test-Path "C:\Tools")) {
    New-Item -ItemType Directory -Path "C:\Tools"
}

Copy-Item $exe $target -Force
Write-Output "ProbeDesk installed in C:\Tools"

# PATH 
$oldPath = [Environment]::GetEnvironmentVariable("PATH", [EnvironmentVariableTarget]::User)
if ($oldPath -notlike "*C:\Tools*") {
    [Environment]::SetEnvironmentVariable("PATH", $oldPath + ";C:\Tools", [EnvironmentVariableTarget]::User)
    Write-Output "C:\Tools added to PATH. Restart your terminal to use 'probedesk' globally."
}

# Autocomplete installation
$profilePath = Join-Path $env:USERPROFILE "Documents\PowerShell\Microsoft.PowerShell_profile.ps1"

if (!(Test-Path $profilePath)) {
    New-Item -ItemType File -Path $profilePath -Force
}

$autocomplete = @"
$flags = @("system","ipconfig","netuse","products","vpn","services","users","usb","trace","remote","report","flushdns","winget-update","scan-health","check-health","restore-health","autocomplete-install")

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

# Only add if not already present
if (-not (Get-Content $profilePath | Select-String "Register-ArgumentCompleter -CommandName `"probedesk`"")) {
    Add-Content -Path $profilePath -Value $autocomplete
    Write-Output "Autocomplete for ProbeDesk installed. Restart PowerShell to activate."
} else {
    Write-Output "Autocomplete already installed."
}