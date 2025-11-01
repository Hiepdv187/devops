# Cấu hình Database Online

## Thông tin đã cập nhật

### 1. Go Application (`database.go`)
Ứng dụng Go giờ đây hỗ trợ kết nối database thông qua biến môi trường:

```bash
DB_HOST=znovxl.h.filess.io
DB_PORT=3306
DB_USER=Wedevops_queenplant
DB_PASSWORD=856333f8b461857adb56cfaa544f250bfae28f9c
DB_NAME=Wedevops_queenplant
```

### 2. PowerShell Restore Script (`mysql-restore-online.ps1`)
Script đã được cấu hình sẵn với thông tin database online.

## Cách sử dụng

### Chạy ứng dụng Go với database online

**Windows PowerShell:**
```powershell
$env:DB_HOST="znovxl.h.filess.io"
$env:DB_PORT="3306"
$env:DB_USER="Wedevops_queenplant"
$env:DB_PASSWORD="856333f8b461857adb56cfaa544f250bfae28f9c"
$env:DB_NAME="Wedevops_queenplant"

go run main.go
```

**Hoặc tạo file `.env`:**
```bash
# .env
DB_HOST=znovxl.h.filess.io
DB_PORT=3306
DB_USER=Wedevops_queenplant
DB_PASSWORD=856333f8b461857adb56cfaa544f250bfae28f9c
DB_NAME=Wedevops_queenplant
```

Sau đó load file `.env` trước khi chạy app (cần cài thư viện `godotenv`):
```go
import "github.com/joho/godotenv"

func init() {
    godotenv.Load()
}
```

### Khôi phục database

Script đã có sẵn thông tin kết nối, chỉ cần chạy:

```powershell
.\mysql-restore-online.ps1
```

Hoặc với file backup khác:
```powershell
.\mysql-restore-online.ps1 -File "my-backup.sql"
```

Hoặc override thông tin kết nối:
```powershell
.\mysql-restore-online.ps1 -DbHost "other-host.com" -User "otheruser" -Password "otherpass"
```

## Lưu ý bảo mật

⚠️ **QUAN TRỌNG**: 
- File `mysql-restore-online.ps1` hiện chứa thông tin nhạy cảm (password)
- **KHÔNG** commit file này lên Git
- Nên thêm vào `.gitignore`:
  ```
  mysql-restore-online.ps1
  .env
  ```

## Kiểm tra kết nối

Test kết nối MySQL từ command line:
```powershell
mysql -h znovxl.h.filess.io -P 3306 -u Wedevops_queenplant -p
# Nhập password: 856333f8b461857adb56cfaa544f250bfae28f9c
```

Sau khi kết nối thành công:
```sql
SHOW DATABASES;
USE Wedevops_queenplant;
SHOW TABLES;
```

## Cấu trúc DSN

Format DSN cho MySQL:
```
user:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True&loc=Local
```

Ví dụ với thông tin hiện tại:
```
Wedevops_queenplant:856333f8b461857adb56cfaa544f250bfae28f9c@tcp(znovxl.h.filess.io:3306)/Wedevops_queenplant?charset=utf8mb4&parseTime=True&loc=Local
```
