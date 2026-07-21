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

- **Sidebar Layout**: Cố định bên trái đối với màn hình Desktop (tự động thu nhỏ thành dạng icon trên Tablet, chuyển xuống thanh tab bar phía dưới đối với Mobile). Tab Notifications độc lập đã được gỡ bỏ.
- **Header**: Nằm cố định phía trên (Sticky Header). Chứa biểu tượng Cái chuông (Bell Icon) cho Notifications (Popover/Drawer) và hình đại diện User. Hiển thị trạng thái Onboarding hoặc cảnh báo "Life Event" chưa xử lý.

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

### 4.4 Giao diện Investor Profile (Hồ sơ đầu tư) & Quản lý Dòng tiền
- **Bố cục Tab (Tabs Layout)**: Chuyển đổi trang Profile thành cấu trúc đa Tab:
  - **Tab Tổng Quan (Overview)**:
    - **Header**: Avatar, tên, điểm FI (Tự do tài chính) + Nút `[Chỉnh sửa Hồ Sơ]` mở Modal cho phép sửa tay (Tuổi, Mục tiêu FI, Tình trạng hôn nhân...).
    - **Trái (Risk & Status)**: Risk Gauge (thước đo rủi ro từ 1-100 với các màu thay đổi từ xanh tới đỏ), Tình trạng hiện tại.
    - **Phải (AI Insights & Constraints)**: Hiển thị các phân tích từ AI (Life Constraints) dưới dạng thẻ highlight và Tổng quan Thu/Chi.
  - **Tab Dòng Tiền (Income Streams)**:
    - Danh sách thẻ (Cards) hiển thị các nguồn thu nhập (Lương, Cho thuê nhà...). Có nhãn phân biệt Chủ động / Thụ động.
    - Nút `[+ Thêm Nguồn Thu Nhập]`. Từng thẻ có nút Sửa/Xóa.
  - **Tab Người Phụ Thuộc (Dependents)**:
    - Danh sách thẻ hiển thị thông tin người phụ thuộc kèm chi phí nuôi dưỡng hàng tháng.
    - Nút `[+ Thêm Người Phụ Thuộc]`. Từng thẻ có nút Sửa/Xóa.

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

### 4.6 Giao diện Portfolio & Watchlist (Danh mục đầu tư & Theo dõi)
- **Top Bar (Controls)**:
  - Nút Switch chuyển đổi giữa 2 chế độ: **[Portfolio]** và **[Watchlist]**.
- **Chế độ Portfolio (Lấy dữ liệu từ bảng Assets)**:
  - **Summary Panel**: 
    - Tổng giá trị đầu tư, Lãi/Lỗ (P/L) tổng hợp, và một biểu đồ Doughnut nhỏ hiển thị Tỷ trọng danh mục.
  - **Bảng Holding (Data Table)**:
    - Bảng hiển thị: Tên tài sản (Mã), Giá vốn (Avg Price), Giá hiện tại, Khối lượng, Tổng giá trị, Tỷ trọng (%), Lãi/Lỗ.
    - Cột Hành động: `[Xem Thesis]`, `[Chỉnh sửa]`.
- **Chế độ Watchlist**:
  - **Bảng Theo dõi**:
    - Hiển thị các tài sản đang đưa vào tầm ngắm nhưng chưa mua.
    - Bảng gồm: Mã, Tên, Giá hiện tại, Giá mục tiêu mua (Target Price), Điểm chất lượng (Quality Score), Độ ưu tiên.
    - Phía bên phải mỗi dòng có nút `[Viết Thesis]`.

### 4.7 Giao diện Investment Thesis (Luận điểm Đầu tư)
- **Màn hình lưới (Grid View)**:
  - Hiển thị danh sách các thẻ (Cards) Thesis. Mỗi thẻ đại diện cho 1 tài sản (VD: HPG, BTC).
  - Thẻ hiển thị: Tên mã, Điểm tin cậy (Conviction Score 1-10), và một câu Tóm tắt luận điểm (Thesis Summary).
  - Nút Nổi to ở góc dưới: `[+ AI Phân Tích Mã Mới]`.
- **Modal Sinh Thesis bằng AI**:
  - Nhập mã tài sản (Ví dụ: FPT).
  - Nút bấm `[🪄 AI Tạo Luận Điểm]`. Hệ thống sẽ hiện Skeleton Loading Glassmorphic.
  - Kết quả trả ra ngay trong Modal để User duyệt trước khi lưu.
- **Trang Chi Tiết (Detail View)**:
  - Chia làm 2 cột:
    - **Cột Trái (Nội dung chính)**: "Tại sao tôi sở hữu?" (Văn bản dài), Lợi thế cạnh tranh (Moats), Động lực (Catalysts), Rủi ro (Risks). Mỗi mục là một thẻ Box có màu viền tương ứng (Xanh dương cho Moat, Xanh lá cho Catalyst, Đỏ cho Risk).
    - **Cột Phải (Metrics & Rules)**: Điểm Conviction (Dạng đồng hồ đo Gauge), Điều kiện bắt buộc bán (Sell Rules). Cột này nằm cố định (Sticky) khi cuộn.
