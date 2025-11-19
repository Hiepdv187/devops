# Cộng đồng Học DevOps

Website được xây dựng bằng [Fiber](https://github.com/gofiber/fiber) giúp cộng đồng chia sẻ và đóng góp bài học DevOps. Ứng dụng sử dụng GORM với PostgreSQL (Supabase) để lưu trữ dữ liệu.

## Yêu cầu

- Go 1.21+

## Cài đặt và chạy

```bash
# Cài dependencies
go mod tidy

# Chạy server
go run .
```

Server mặc định chạy tại `http://localhost:3000`.

## Cấu hình Database (Supabase)

### Bước 1: Tạo Supabase Project

1. Truy cập [https://supabase.com](https://supabase.com) và tạo tài khoản
2. Tạo project mới và lưu lại **Database Password**
3. Vào **Settings > Database** để lấy thông tin kết nối

### Bước 2: Cấu hình .env

Tạo file `.env` từ `.env.example` và điền thông tin Supabase:

```env
DB_HOST=db.xxxxxxxxxxxxx.supabase.co
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_supabase_password
DB_NAME=postgres
DB_SSLMODE=require
```

Hoặc sử dụng connection string trực tiếp:

```env
DATABASE_DSN=postgresql://postgres:your_password@db.xxxxx.supabase.co:5432/postgres
```

📖 **Chi tiết hướng dẫn setup**: Xem file [SUPABASE_SETUP.md](./SUPABASE_SETUP.md)

Sau khi kết nối thành công, hệ thống sẽ tự động migrate schema và thêm tài khoản mẫu:

- Email: `admin@hocdevops.community`
- Mật khẩu: `devops123`

Các tính năng chính:

- Giao diện web với trang chủ, danh sách bài viết, chi tiết bài viết, đăng ký/đăng nhập và biểu mẫu đóng góp nội dung.
- Hệ thống session lưu đăng nhập, hỗ trợ đăng xuất và flash message thông báo.
- API REST cho thao tác đăng ký, đăng nhập, tạo bài viết, bình luận.

Các endpoint quan trọng:

- `POST /auth/register`: tạo tài khoản mới cho cộng đồng.
- `POST /auth/login`: đăng nhập, trả về thông tin người dùng.
- `POST /auth/logout`: đăng xuất.
- `POST /posts`: tạo bài viết mới sau khi có `author_id` hợp lệ hoặc người dùng đã đăng nhập.
- `POST /posts/:id/comments`: thêm bình luận cho bài viết.

Endpoint đăng nhập nhận payload dạng JSON:

```json
{
  "email": "admin@hocdevops.community",
  "password": "devops123"
}
```

Trả về thông tin người dùng sau khi xác thực thành công. Với giao diện web, bạn có thể vào `/auth/register` và `/auth/login` để thao tác bằng form.

## Cấu trúc thư mục

```
.
├── internal
│   └── handlers       # Logic xử lý request và dữ liệu demo
├── public             # Static assets (CSS, hình ảnh)
├── views
│   ├── layouts        # Template layout chính
│   └── pages          # Trang con
├── go.mod / go.sum
└── main.go
```

## Đóng góp

1. Fork project và tạo branch mới.
2. Thêm hoặc chỉnh sửa nội dung bài học trong `internal/handlers` và template trong `views`.
3. Mở pull request mô tả rõ thay đổi.

Liên hệ quản trị viên: `hello@hocdevops.community`.

## Backup & Migration

Supabase tự động backup database hàng ngày. Để migration hoặc backup thủ công, xem hướng dẫn chi tiết tại [SUPABASE_SETUP.md](./SUPABASE_SETUP.md).