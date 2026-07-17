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
