# Life Capital — Frontend Layout & UI/UX

> Hướng dẫn thiết kế giao diện Vue 3, cấu trúc Layout, CSS Glassmorphic và Wireframes mẫu.

---

## 1. Giao Diện Bản Thiết Kế (UI/UX System)

Ứng dụng hướng tới phong cách tối giản, sang trọng (premium dark theme), tối ưu hóa trải nghiệm đọc dữ liệu tài chính dài hạn:
- **Theme chủ đạo**: Dark Mode làm mặc định.
- **Bảng màu (Hex codes)**:
  - Nền ứng dụng (Background): Slate 900 (`#0F172A`)
  - Bề mặt thẻ (Surface): Slate 800 với opacity (`rgba(30, 41, 59, 0.7)`)
  - Viền thẻ (Border): Translucent slate border (`rgba(255, 255, 255, 0.05)`)
  - Màu nhấn chính (Accent): Indigo 500 (`#6366F1`) & Violet 500 (`#8B5CF6`)
  - Trạng thái tích cực/Tăng trưởng (Growth/Success): Emerald 500 (`#10B981`)
  - Cảnh báo/Lệch tỷ trọng (Warning): Amber 500 (`#F59E0B`)
  - Nguy hiểm/Thesis bị vi phạm (Danger/Alert): Red 500 (`#EF4444`)

---

## 2. CSS Hiệu Ứng Glassmorphic

Tất cả các thẻ hiển thị (cards), sidebar và modal đều kế thừa class sau để tạo chiều sâu trực quan:

```css
.premium-card {
  background: rgba(30, 41, 59, 0.7);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  box-shadow: 0 8px 32px 0 rgba(0, 0, 0, 0.3);
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.premium-card:hover {
  transform: translateY(-2px);
  border-color: rgba(99, 102, 241, 0.2);
  box-shadow: 0 12px 40px 0 rgba(99, 102, 241, 0.05);
}
```

---

## 3. Cấu Trúc Bố Cục (Layout Hierarchy)

```
┌────────────────────────────────────────────────────────┐
│  Sidebar (Nav)  │  Header (Breadcrumb + Profile Alert) │
│  - Dashboard    ├──────────────────────────────────────┤
│  - Timeline     │                                      │
│  - Assets & Debts│  Main Content Area                  │
│  - Goals        │  (RouterView)                        │
│  - Portfolio    │                                      │
│  - Thesis       │                                      │
│  - Decisions    │                                      │
│  - Reviews      │                                      │
└─────────────────┴──────────────────────────────────────┘
```

- **Sidebar Layout**: Cố định bên trái đối với màn hình Desktop (tự động thu nhỏ thành dạng icon trên Tablet, chuyển xuống thanh tab bar phía dưới đối với Mobile).
- **Header**: Hiển thị trạng thái Onboarding hoặc cảnh báo "Life Event" chưa xử lý kèm theo lối tắt tới công cụ hỏi đáp AI Advisor.

---

## 4. UI Wireframes Đặc Trưng

### 4.1 Onboarding AI Interview (Trải nghiệm dạng hội thoại chat)
- Giao diện chia làm 2 phần:
  - **Trái (Chat Panel)**: Trình bày bong bóng chat từ AI Advisor và khung nhập câu trả lời của user. Phía dưới có thanh tiến trình 8 bước kèm các gợi ý nhanh (quick-select chips) như: `[Đã kết hôn]`, `[Chưa có con]`.
  - **Phải (Profile Preview)**: Hiển thị các thông tin đã thu thập dưới dạng realtime cards (như Net Worth dự kiến, danh sách dependents).

### 4.2 Dashboard chính
- **Cột Trái**: Bảng cân đối Net Worth rút gọn (Assets | Liabilities | Net Worth) và danh sách thẻ tiến trình các mục tiêu tài chính (Emergency Fund, VinFast Auto, Retire Early...) được rót vốn tự động theo cơ chế Waterfall.
- **Cột Phải**: Biểu đồ phân bổ tỷ trọng tài sản thực tế (Allocation Bar) và danh sách cảnh báo lệch mục tiêu IPS hoặc cảnh báo đòn bẩy nợ cao.
- **Dưới**: Nút hành động nổi bật "Start Monthly Review" màu Violet.

### 4.3 Giao diện Auth (Đăng nhập / Đăng ký)
- Sử dụng cấu trúc Card đặt ở giữa màn hình (centered layout).
- Form bao gồm các trường thông tin tiêu chuẩn (Email, Mật khẩu, Tên).
- Cung cấp tính năng "Đăng nhập bằng Google" nổi bật (Sử dụng `vue3-google-login`).
- Hiển thị cảnh báo lỗi bằng chữ đỏ nhạt viền đỏ khi nhập sai hoặc trùng lặp tài khoản.

### 4.4 Giao diện Investor Profile (Hồ sơ đầu tư)
- **Header**: Avatar, tên, tổng thu nhập/chi phí và điểm FI (Tự do tài chính).
- **Trái (Risk & Status)**: Risk Gauge (thước đo rủi ro từ 1-100 với các màu thay đổi từ xanh tới đỏ), Tình trạng hiện tại (Trạng thái hôn nhân, Số người phụ thuộc).
- **Phải (AI Insights & Constraints)**: Hiển thị các phân tích từ AI (Life Constraints) dưới dạng thẻ highlight và đề xuất phân bổ vốn.

### 4.5 Giao diện Assets & Liabilities (Quản lý Tài sản & Nợ)
- **Khu vực Top (Summary Cards)**:
  - 3 thẻ Glassmorphism lớn nằm ngang với viền, bóng đổ và biểu tượng chìm (watermark) cực kỳ nổi bật:
    - `Tổng Tài Sản`: Gradient Emerald (Xanh ngọc).
    - `Tổng Nợ`: Gradient Amber (Vàng hổ phách).
    - `Tài Sản Ròng (Net Worth)`: Card nổi bật nhất, bọc viền gradient Indigo/Violet, chữ siêu lớn.
- **Khu vực Toolbar (Bộ lọc - Filters)**:
  - Thanh công cụ nằm dưới thẻ Summary, chứa các Dropdown cho phép lọc riêng biệt theo danh mục (Category) cho cả Tài sản và Nợ.
- **Khu vực Bottom (2 Columns Grid)**:
  - **Cột Trái (Tài Sản)**: 
    - Tiêu đề "Tài sản của bạn" + Nút `[+ Thêm Tài Sản]`.
    - Danh sách tài sản được bọc trong một container cố định chiều cao (max-height), hỗ trợ cuộn vô hạn (infinite scrolling) thông qua Virtual Scroll.
    - Dữ liệu được gom nhóm theo Category. Tiêu đề nhóm có thể click để mở Popup chi tiết.
  - **Cột Phải (Khoản Nợ)**: 
    - Tiêu đề "Các khoản nợ" + Nút `[+ Thêm Khoản Nợ]`.
    - Danh sách khoản nợ tương tự cột trái (Gom nhóm, có thể click xem chi tiết).
  - Từng Item trong danh sách có nút `Sửa` và `Xóa` xuất hiện khi hover.
- **Modal Thêm Mới / Cập Nhật (Edit)**: 
  - Form nổi (Overlay) mờ ảo dùng chung cho cả tính năng Thêm và Sửa.
  - Dropdown chọn Category (Khóa không cho sửa khi đang ở chế độ Edit).
  - Tự động thay đổi các trường dữ liệu tùy theo Loại tài sản (Ví dụ: Vàng sẽ hiện form chọn thương hiệu/tuổi vàng thay vì nhập tay).
- **Modal Chi Tiết Nhóm (Group Detail Popup)**:
  - Hiện ra giữa màn hình khi click vào tiêu đề Nhóm (Ví dụ: Cổ phiếu).
  - Hiển thị danh sách các tài sản/nợ con thuộc nhóm đó.
  - Mỗi item đi kèm một thanh tiến trình (Progress Bar) trực quan để biểu thị tỷ trọng (%) của item đó so với tổng giá trị nội bộ nhóm.
