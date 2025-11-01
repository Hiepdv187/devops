param (
    [string]$File = "backup.sql",
    [string]$DbHost = "znovxl.h.filess.io",
    [int]$Port = 3306,
    [string]$User = "Wedevops_queenplant",
    [string]$Password = "856333f8b461857adb56cfaa544f250bfae28f9c",
    [string]$Database = "Wedevops_queenplant"
)

Write-Host ""
Write-Host "=== Bat dau khoi phuc CSDL MySQL Online ===" -ForegroundColor Cyan

# Kiem tra tham so bat buoc
if ([string]::IsNullOrWhiteSpace($DbHost)) {
    Write-Host "Loi: Can cung cap DbHost (dia chi server database)!" -ForegroundColor Red
    Write-Host "Vi du: .\mysql-restore-online.ps1 -DbHost 'db.example.com' -User 'username' -Password 'password'" -ForegroundColor Yellow
    exit 1
}

if ([string]::IsNullOrWhiteSpace($User)) {
    Write-Host "Loi: Can cung cap User (ten nguoi dung)!" -ForegroundColor Red
    exit 1
}

if ([string]::IsNullOrWhiteSpace($Password)) {
    Write-Host "Loi: Can cung cap Password (mat khau)!" -ForegroundColor Red
    exit 1
}

if (-not (Test-Path $File)) {
    Write-Host "Loi: File '$File' khong ton tai!" -ForegroundColor Red
    exit 1
}

Write-Host "Thong tin ket noi:" -ForegroundColor Yellow
Write-Host "  Host: $DbHost" -ForegroundColor White
Write-Host "  Port: $Port" -ForegroundColor White
Write-Host "  User: $User" -ForegroundColor White
Write-Host "  Database: $Database" -ForegroundColor White
Write-Host ""

# Kiem tra mysql command co san khong
$mysqlPath = Get-Command mysql -ErrorAction SilentlyContinue
if (-not $mysqlPath) {
    Write-Host "Loi: Khong tim thay lenh 'mysql'. Vui long cai dat MySQL Client!" -ForegroundColor Red
    Write-Host "Tai ve tai: https://dev.mysql.com/downloads/mysql/" -ForegroundColor Yellow
    exit 1
}

Write-Host "Dang kiem tra ket noi toi database..." -ForegroundColor Green

# Test ket noi
try {
    $testQuery = "SELECT 1"
    $testResult = & mysql -h $DbHost -P $Port -u $User -p"$Password" -e $testQuery 2>&1
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Loi: Khong the ket noi toi database!" -ForegroundColor Red
        Write-Host $testResult -ForegroundColor Red
        exit 1
    }
    
    Write-Host "Ket noi thanh cong!" -ForegroundColor Green
}
catch {
    Write-Host "Loi: Khong the ket noi toi database!" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "Dang xu ly file backup..." -ForegroundColor Green

# Doc va xu ly file backup - thay the collation khong tuong thich
$sqlContent = Get-Content -Path $File -Raw -Encoding UTF8

# Thay the cac collation MySQL 8.0 thanh phien ban tuong thich
Write-Host "Dang chuyen doi collation cho tuong thich voi MySQL 5.7..." -ForegroundColor Yellow
$sqlContent = $sqlContent -replace 'utf8mb4_0900_ai_ci', 'utf8mb4_unicode_ci'
$sqlContent = $sqlContent -replace 'utf8_0900_ai_ci', 'utf8_unicode_ci'

Write-Host "Dang khoi phuc du lieu vao database '$Database'..." -ForegroundColor Green
Write-Host "Vui long cho, qua trinh nay co the mat vai phut..." -ForegroundColor Yellow

try {
    # Pipe noi dung da xu ly vao mysql
    $sqlContent | & mysql -h $DbHost -P $Port -u $User -p"$Password" --default-character-set=utf8mb4 $Database 2>&1 | Tee-Object -Variable output
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "`nKhoi phuc thanh cong!" -ForegroundColor Green
        
        # Hien thi thong tin database sau khi restore
        Write-Host "`nThong tin database sau khi khoi phuc:" -ForegroundColor Cyan
        $tableCount = & mysql -h $DbHost -P $Port -u $User -p"$Password" -D $Database -e "SHOW TABLES;" -s 2>&1
        if ($tableCount) {
            $tables = ($tableCount | Measure-Object -Line).Lines
            Write-Host "  So luong bang: $tables" -ForegroundColor White
        }
    }
    else {
        Write-Host "`nCo loi xay ra khi khoi phuc!" -ForegroundColor Red
        if ($output) {
            Write-Host $output -ForegroundColor Red
        }
        exit 1
    }
}
catch {
    Write-Host "`nCo loi xay ra khi khoi phuc:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "=== Hoan tat ===" -ForegroundColor Cyan
Write-Host ""
