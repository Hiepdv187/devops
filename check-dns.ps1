# Kiểm tra DNS resolution cho Supabase host
$host = "db.gtdxzzzibtyhnwhyfwuo.supabase.co"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "KIỂM TRA SUPABASE HOST" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Host: $host" -ForegroundColor Yellow
Write-Host ""

try {
    Write-Host "Đang kiểm tra DNS resolution..." -ForegroundColor Gray
    $result = Resolve-DnsName -Name $host -ErrorAction Stop
    Write-Host "✅ DNS OK! Host tồn tại" -ForegroundColor Green
    Write-Host ""
    $result | Format-Table -AutoSize
} catch {
    Write-Host "❌ DNS FAILED! Host không tồn tại" -ForegroundColor Red
    Write-Host ""
    Write-Host "Lỗi: $_" -ForegroundColor Red
    Write-Host ""
    Write-Host "🔍 Khắc phục:" -ForegroundColor Yellow
    Write-Host "1. Truy cập: https://supabase.com/dashboard" -ForegroundColor White
    Write-Host "2. Kiểm tra project có tồn tại không" -ForegroundColor White
    Write-Host "3. Vào Settings > Database" -ForegroundColor White
    Write-Host "4. Copy lại Connection string (URI) chính xác" -ForegroundColor White
    Write-Host "5. Host phải có dạng: db.xxxxxxxxxxxxx.supabase.co" -ForegroundColor White
    Write-Host ""
    Write-Host "⚠️  Lưu ý: Project Reference ID trong URL phải CHÍNH XÁC" -ForegroundColor Yellow
}
