param (
    [string]$File = "backup.sql",
    [string]$Container = "fiber-learning-db",
    [string]$User = "root",
    [string]$Password = "rootpass",
    [string]$Database = "fiber_learning"
)

Write-Host ""
Write-Host "=== Bat dau khoi phuc CSDL MySQL ===" -ForegroundColor Cyan

if (-not (Test-Path $File)) {
    Write-Host "File '$File' khong ton tai!" -ForegroundColor Red
    exit 1
}


Write-Host "Dang khoi phuc du lieu vao '$Database'..." -ForegroundColor Green

try {
    $sqlContent = Get-Content -Path $File -Raw -Encoding UTF8
    $bytes = [System.Text.Encoding]::UTF8.GetBytes($sqlContent)
    $stream = New-Object System.IO.MemoryStream
    $stream.Write($bytes, 0, $bytes.Length)
    $stream.Seek(0, 'Begin') | Out-Null

    $processInfo = New-Object System.Diagnostics.ProcessStartInfo
    $processInfo.FileName = "docker"
    $processInfo.Arguments = "exec -i $Container mysql --default-character-set=utf8mb4 -u $User -p$Password $Database"
    $processInfo.RedirectStandardInput = $true
    $processInfo.UseShellExecute = $false
    $processInfo.CreateNoWindow = $true

    $process = New-Object System.Diagnostics.Process
    $process.StartInfo = $processInfo
    $process.Start() | Out-Null

    $stream.CopyTo($process.StandardInput.BaseStream)
    $process.StandardInput.Close()
    $process.WaitForExit()

    Write-Host "`nDone!" -ForegroundColor Green
}
catch {
    Write-Host "`nCo loi xay ra khi khoi phuc:" -ForegroundColor Red
    Write-Host $_.Exception.Message
}

Write-Host ""
Write-Host "=== Hoan tat ===" -ForegroundColor Cyan
