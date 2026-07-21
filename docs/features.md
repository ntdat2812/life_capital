# Life Capital — Full Feature List & Screen Specifications

> Tài liệu mô tả chi tiết tất cả các tính năng cần phát triển và danh sách màn hình tương ứng.

---

## 1. Feature Specifications (Chi tiết tính năng)

### 1.1 Module 0 — Investor Profile & Life Timeline (Nền tảng)
- **AI Phỏng vấn Onboarding**: Hệ thống hội thoại 8 bước hướng dẫn người dùng cung cấp thông tin tài chính cá nhân. AI phân tích và lưu vào `investor_profiles` dạng cấu trúc.
- **Dòng thời gian Life Timeline**: Cho phép xem lịch sử sự kiện biến động cuộc sống (đổi việc, kết hôn, mua nhà) sắp xếp theo trình tự thời gian.
- **Động cơ Cascade**: Khi có sự kiện đời sống mới:
  1. AI tính toán sự thay đổi thu nhập/chi phí và phân loại khẩu vị rủi ro mới.
  2. Tạo bản ghi Profile version mới ở trạng thái `superseded` đối với bản cũ.
  3. Đề xuất điều chỉnh IPS tương ứng.

### 1.2 Module 1 — Dashboard
- **Tổng quan tài sản (Net Worth)**: Hiển thị tổng Net Worth = Assets - Liabilities. Biểu đồ đường thể hiện sự tăng trưởng Net Worth qua các tháng.
- **Tiến trình FI**: Thể hiện tỷ lệ % tiến gần tới số tự do tài chính. Hiển thị số năm dự kiến để đạt mục tiêu dựa trên Savings Rate thực tế.
- **Cảnh báo chiến lược**: Hiển thị các cảnh báo quan trọng như: Lệch phân bổ IPS, cổ phiếu chạm giá mua watchlist, kết quả kinh doanh mới yếu đi.

### 1.3 Module 2 — Assets & Liabilities Management
- **Quản lý đa dạng lớp tài sản (Assets)**: CRUD các lớp tài sản gồm Tiền mặt, Tiền gửi tiết kiệm, Vàng, Cổ phiếu, Bất động sản.
- **Quản lý các khoản nợ (Liabilities)**: CRUD các khoản nợ vay mua nhà, mua xe, nợ thẻ tín dụng kèm theo dư nợ còn lại và số tiền trả góp hàng tháng.
- **Tự động tính toán ròng**: Tính toán tổng Assets, tổng Liabilities, tỷ lệ đòn bẩy (Debt-to-Asset) và Net Worth thực tế.

### 1.4 Module 3 — Portfolio & Market Data Sync
- **Danh mục tài sản**: Bảng thống kê chi tiết các mã đang hold, số lượng, giá vốn, giá hiện tại, tỷ trọng thực tế và tỷ trọng mục tiêu (Target Allocation).
- **Tự động cập nhật giá**: Sử dụng Background Worker (Cronjob) và tích hợp Public APIs (TCBS, SSI, Yahoo Finance) để cào giá trị tài sản hàng ngày (cuối ngày).
- **Khuyến nghị & Conviction**: Điểm tin cậy (Conviction Score 1-10) của nhà đầu tư cho cổ phiếu và khuyến nghị hành động tương ứng (BUY/HOLD/SELL).
- **Hệ thống Thông báo (Notification System)**: Bắn cảnh báo tự động khi giá chạm ngưỡng mua (Watchlist) hoặc tỷ trọng tài sản vượt quá giới hạn rủi ro.

### 1.5 Module 4 — Investment Policy Statement (IPS)
- **Quy tắc đầu tư (IPS)**: Thiết lập mục tiêu chiến lược dài hạn, tỷ trọng tối đa/tối thiểu cho từng lớp tài sản (Target Allocation), quy tắc mua/bán và giới hạn rủi ro. (Đóng vai trò là "Luật chơi Vĩ mô" cho toàn bộ tài sản).
- **AI Tự động Đề xuất (AI Generator)**: Kể cả khi chưa có tài sản, AI sẽ tự động sinh ra một bản Nháp IPS (Draft) dựa trên Investor Profile và Life Event để gợi ý tỷ trọng nên phân bổ. User có toàn quyền tinh chỉnh (Manual Override) và lưu thành bản chính thức.
- **Phiên bản (Versioning)**: IPS lưu lịch sử dưới dạng phiên bản (v1, v2, v3) để đối chiếu sự thay đổi tư duy đầu tư qua thời gian.

### 1.6 Module 5 — Investment Thesis (Luận điểm đầu tư)
- **Tại sao tôi sở hữu? (Why I Own)**: Ghi lại lý do nắm giữ cốt lõi của từng doanh nghiệp/mã cụ thể. (Đóng vai trò là "Lý do Vi mô" cho 1 tài sản).
- **AI Viết hộ Thesis**: Người dùng không chuyên chỉ cần nhập mã cổ phiếu, AI sẽ tự động phân tích doanh nghiệp và sinh ra một bản Thesis tiêu chuẩn (bao gồm Lợi thế cạnh tranh - Moats, Động lực - Catalysts, và Rủi ro) để user tham khảo.
- **Sell Rules**: Điều kiện bắt buộc phải bán (Thesis bị phá vỡ, quản trị doanh nghiệp xấu đi).

### 1.7 Module 6 — Earnings Review (Đánh giá BCTC)
- **AI Đọc BCTC**: Sau mỗi quý, người dùng nhập link hoặc paste nội dung Báo cáo tài chính/KQKD của doanh nghiệp.
- **Tự động đối chiếu Thesis**: AI tự động đọc BCTC mới nhất và so sánh thẳng với các Luận điểm (Thesis) đã viết ở Module 5. Nếu phát hiện Moat bị phá vỡ hoặc tăng trưởng giảm, AI sẽ bắn Notification cảnh báo yêu cầu user xem xét bán.
- **Chấm điểm**: AI chấm điểm sức khỏe quý (1-10).

### 1.8 Module 7 — Decision Journal (Nhật ký quyết định)
- **Ghi nhận hành động**: Lưu lại nhật ký mua/bán, số tiền, giá, lý do ra quyết định và trạng thái tâm lý (FOMO, bình tĩnh, lo sợ).
- **Đánh giá lại**: Sau 3/6/12 tháng, hệ thống nhắc nhở người dùng review lại quyết định đó là đúng hay sai để rút ra bài học.

### 1.9 Module 8 — Monthly Review
- **Nhập liệu tháng**: Nhập số vốn mới mang đi đầu tư, thu nhập/chi phí thực tế trong tháng.
- **AI Report**: AI quét toàn bộ danh mục, IPS, Thesis và KQKD mới để đưa ra báo cáo tổng quan tháng cùng lệnh phân bổ tiền mới cụ thể.

### 1.10 Module 9 — Authentication & Authorization
- **Đăng ký (Signup)**: Multi-user system. Nếu email đăng ký đã tồn tại, API sẽ trả về lỗi `409 Conflict`.
- **Đăng nhập (Login)**: Nhập email + mật khẩu. Backend kiểm tra mật khẩu bằng bcrypt, nếu khớp sẽ trả về JWT access token.
- **JWT Token Flow**: Mọi API endpoint (trừ `/api/v1/auth/signup`, `/api/v1/auth/login`) yêu cầu header `Authorization: Bearer <token>`. Token hết hạn sau 24 giờ.
- **Logout**: Phía client xóa token khỏi localStorage. Không cần API phía server (stateless JWT).

---

## 2. Screen Specifications (Danh sách màn hình)

| ID | Tên màn hình | Đường dẫn Route (Vue Router) | Thành phần chính |
|---|---|---|---|
| 1 | Dashboard | `/` | Net Worth Card, Allocation Chart, FI Progress, Alerts |
| 2 | Investor Profile | `/profile` | Profile Summary, Key Metrics, Risk Gauge, AI Summary |
| 3 | Life Timeline | `/profile/timeline` | Timeline View, Event Cards, Impact Indicators |
| 4 | Log Life Event | `/profile/timeline/new` | Event Form, AI Impact Preview, Cascade |
| 5 | Dependents | `/profile/dependents` | Dependent Cards, Monthly Cost Summary |
| 6 | Income Streams | `/profile/income` | Income Table, Active/Passive Split, Trend |
| 7 | Assets & Liabilities | `/assets` | Assets & Liabilities tables, Add forms, Allocation Bar |
| 8 | Portfolio | `/portfolio` | Holdings Table, Allocation Bar, Import button |
| 9 | IPS Current | `/ips` | Policy Document View, Target Allocations |
| 10 | Thesis Grid | `/thesis` | Company Cards Grid, Filters |
| 11 | Thesis Detail | `/thesis/:ticker` | Full Thesis View, Conviction Gauge |
| 12 | Decision Journal | `/decisions` | Timeline View, Filters |
| 13 | Monthly Review | `/review` | Status Overview, Start Review CTA |
| 14 | Active Review | `/review/:month` | Step Wizard (Life Check -> Input -> AI -> Recs) |
| 15 | Onboarding Interview | `/onboarding/interview` | Conversational AI Chat, Progress Bar |
| 16 | Login | `/login` | Form đăng nhập credentials |
| 17 | Financial Goals | `/goals` | Goal Waterfall Cards, Add Goal Form |
| 18 | Signup | `/signup` | Form đăng ký tài khoản (multi-user) |
