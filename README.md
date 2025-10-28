# Cộng đồng Học DevOps

Website được xây dựng bằng [Fiber](https://github.com/gofiber/fiber) giúp cộng đồng chia sẻ và đóng góp bài học DevOps. Ứng dụng sử dụng GORM với SQLite để lưu trữ tài khoản demo.

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

Ứng dụng sử dụng MySQL thông qua GORM. Cấu hình chuỗi kết nối qua biến môi trường `DATABASE_DSN`, ví dụ:

```
DATABASE_DSN="user:password@tcp(127.0.0.1:3306)/fiber_learning?charset=utf8mb4&parseTime=True&loc=Local"
```

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
