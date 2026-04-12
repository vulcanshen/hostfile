# hostfile uninstaller for Windows
# Usage: irm https://raw.githubusercontent.com/vulcanshen/hostfile/main/uninstall.ps1 | iex

$ErrorActionPreference = "Stop"

$installDir = "$env:LOCALAPPDATA\hostfile"

if (-not (Test-Path "$installDir\hostfile.exe")) {
    Write-Host "hostfile not found in $installDir" -ForegroundColor Red
    exit 1
}

Remove-Item $installDir -Recurse -Force
Write-Host "removed $installDir" -ForegroundColor Green

# Remove from PATH
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -like "*$installDir*") {
    $newPath = ($userPath -split ";" | Where-Object { $_ -ne $installDir }) -join ";"
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    Write-Host "removed $installDir from PATH" -ForegroundColor Green
}

# Remove saved snapshots
$saveDir = Join-Path $env:USERPROFILE ".hostfile"
if (Test-Path $saveDir) {
    $answer = Read-Host "Remove saved snapshots in $saveDir? [y/N]"
    if ($answer -eq "y" -or $answer -eq "Y") {
        Remove-Item $saveDir -Recurse -Force
        Write-Host "removed $saveDir" -ForegroundColor Green
    } else {
        Write-Host "kept $saveDir" -ForegroundColor Yellow
    }
}

Write-Host ""
Write-Host "hostfile uninstalled. Restart your terminal for PATH changes to take effect." -ForegroundColor Green
