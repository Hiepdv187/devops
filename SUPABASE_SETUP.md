# Hướng dẫn kết nối Supabase Database

## 📋 Tổng quan

Project này đã được cấu hình sẵn để kết nối với **Supabase PostgreSQL Database**. Supabase là một nền tảng Backend-as-a-Service (BaaS) mã nguồn mở, cung cấp PostgreSQL database với nhiều tính năng mạnh mẽ.

## 🚀 Các bước setup

### Bước 1: Tạo project trên Supabase

1. Truy cập [https://supabase.com](https://supabase.com)
2. Đăng ký/Đăng nhập tài khoản (có thể dùng GitHub)
3. Click **"New Project"**
4. Điền thông tin:
   - **Name**: Tên project của bạn (ví dụ: `wedevops`)
   - **Database Password**: Tạo mật khẩu mạnh và **LƯU LẠI** (sẽ cần dùng sau)
   - **Region**: Chọn region gần nhất (ví dụ: `Southeast Asia (Singapore)`)
5. Click **"Create new project"** và đợi vài phút để Supabase khởi tạo

### Bước 2: Lấy thông tin kết nối Database

1. Vào project vừa tạo, click **Settings** (icon bánh răng ⚙️) ở sidebar bên trái
2. Click **Database** trong menu Settings
3. Scroll xuống phần **"Connection string"**
4. Chọn tab **"URI"** 
5. Copy connection string có dạng:
   ```
   postgresql://postgres:[YOUR-PASSWORD]@db.xxxxxxxxxxxxx.supabase.co:5432/postgres
   ```

**Lưu ý:** 
- Phần `xxxxxxxxxxxxx` là Project Reference ID của bạn
- `[YOUR-PASSWORD]` là password bạn đã tạo ở Bước 1
- Nếu quên password, có thể reset tại: **Settings > Database > Reset database password**

### Bước 3: Cấu hình file .env

1. Mở file `.env` trong project (đã được tạo sẵn)
2. **Option 1 - Điền từng biến riêng lẻ (KHUYẾN NGHỊ):**

```env
DB_HOST=db.xxxxxxxxxxxxx.supabase.co
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=postgres
DB_SSLMODE=require
```

3. **Option 2 - Dùng connection string trực tiếp:**

```env
DATABASE_DSN=postgresql://postgres:your_password@db.xxxxxxxxxxxxx.supabase.co:5432/postgres
```

**Lưu ý:**
- Thay `xxxxxxxxxxxxx` bằng Project Reference ID thực tế
- Thay `your_password` bằng password database thực tế
- **KHÔNG commit file `.env`** vào Git (đã có trong `.gitignore`)

### Bước 4: Chạy ứng dụng

```bash
# Cài đặt dependencies (nếu chưa cài)
go mod download

# Chạy ứng dụng
go run main.go
```

Ứng dụng sẽ tự động:
- ✅ Kết nối tới Supabase database
- ✅ Chạy migration (tạo các bảng cần thiết)
- ✅ Seed dữ liệu demo user

### Bước 5: Kiểm tra kết nối thành công

Nếu mọi thứ OK, bạn sẽ thấy log:
```
✓ Database connected successfully
✓ Migration completed
```

Kiểm tra database trên Supabase:
1. Vào **Table Editor** trong Supabase Dashboard
2. Bạn sẽ thấy các bảng đã được tạo:
   - `users`
   - `posts`
   - `comments`
   - `annotations`
   - `books`
   - `book_pages`
   - `highlights`

## 📊 Xem dữ liệu trong Supabase

Supabase cung cấp **Table Editor** rất trực quan:

1. Click **Table Editor** ở sidebar
2. Chọn bảng muốn xem (ví dụ: `users`)
3. Có thể:
   - Xem, thêm, sửa, xóa dữ liệu trực tiếp
   - Chạy SQL queries
   - Export data

## 🔐 Bảo mật

### Không commit .env file
File `.env` chứa thông tin nhạy cảm, **KHÔNG BAO GIỜ** commit vào Git:

```bash
# Kiểm tra .gitignore đã có dòng này
.env
```

### Sử dụng Row Level Security (RLS)

Supabase hỗ trợ RLS để bảo vệ dữ liệu. Tham khảo:
- [Supabase RLS Documentation](https://supabase.com/docs/guides/auth/row-level-security)

## 🛠️ Troubleshooting

### Lỗi "connection refused"
- ✅ Kiểm tra DB_HOST có đúng không (phải có dạng `db.xxxxx.supabase.co`)
- ✅ Kiểm tra internet connection
- ✅ Kiểm tra Supabase project có đang hoạt động không

### Lỗi "password authentication failed"
- ✅ Kiểm tra DB_PASSWORD có đúng không
- ✅ Reset password tại Settings > Database nếu cần

### Lỗi "SSL connection required"
- ✅ Đảm bảo `DB_SSLMODE=require` hoặc `sslmode=require` trong connection string

### Lỗi "database does not exist"
- ✅ DB_NAME phải là `postgres` (default database của Supabase)

## 📚 Tài liệu tham khảo

- [Supabase Documentation](https://supabase.com/docs)
- [Supabase Database](https://supabase.com/docs/guides/database)
- [GORM Documentation](https://gorm.io/docs/)

## 🔄 Migration từ Neon sang Supabase

Nếu bạn đang có dữ liệu trên Neon và muốn chuyển sang Supabase:

### Option 1: Dump & Restore (Khuyến nghị)

```bash
# 1. Dump data từ Neon
pg_dump "postgresql://neondb_owner:npg_ZCva2xmGOt3g@ep-odd-morning-a7ox4rz0-pooler.ap-southeast-2.aws.neon.tech/wedevops" > backup.sql

# 2. Restore vào Supabase
psql "postgresql://postgres:YOUR_PASSWORD@db.xxxxx.supabase.co:5432/postgres" < backup.sql
```

### Option 2: Tạo mới và để app tự migrate

Application sẽ tự động tạo schema và seed demo data khi chạy lần đầu.

## 💡 Tips

1. **Free tier của Supabase** bao gồm:
   - 500 MB database storage
   - Unlimited API requests
   - 50,000 monthly active users
   - Row Level Security
   - Auto-generated APIs

2. **Connection pooling**: Supabase tự động có connection pooling, giúp app xử lý nhiều requests hiệu quả hơn

3. **Backups**: Supabase tự động backup database hàng ngày (trên paid plans), hoặc bạn có thể manual backup qua dashboard

## 📞 Hỗ trợ

Nếu gặp vấn đề, có thể:
- Check logs của ứng dụng
- Xem Supabase logs tại **Logs** section trong dashboard
- Tham khảo [Supabase Community](https://github.com/supabase/supabase/discussions)

---

**Chúc bạn setup thành công! 🎉**
