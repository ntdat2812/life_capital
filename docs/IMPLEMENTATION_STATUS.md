# Implementation Status

Tài liệu này theo dõi tiến độ thực tế (Source Code) so với đặc tả tính năng trong `docs/features.md`. 
*Mục đích: Giúp các Agent AI đọc và biết chính xác cái gì đã làm xong, cái gì còn thiếu.*

## Status Categories
- 🟢 **Hoàn thành (Done)**: Đã code xong frontend, backend, test chạy ổn định.
- 🟡 **Đang phát triển (WIP / Partial)**: Mới làm được một phần, hoặc giao diện cứng chưa nối API.
- 🔴 **Chưa làm (Pending)**: Hoàn toàn chưa có code.

---

## Mức độ hoàn thiện theo Module

### Module 0 — Investor Profile & Life Timeline
- 🟡 **AI Phỏng vấn Onboarding**: Có giao diện chat (`OnboardingView.vue`), nhưng chưa linh hoạt (hardcode 8 bước).
- 🟢 **Profile Của Tôi**: Giao diện hiển thị thông tin đa Tab (Overview, Thu nhập, Phụ thuộc). Đã có chức năng Edit thủ công (Sửa Tuổi, FI Target...).
- 🟢 **Quản lý Dòng tiền (Income/Expense/Dependents)**: Đã có giao diện CRUD cho Income/Dependents. Đã cấu trúc Cashflow theo hướng vĩ mô (Lump Sum Budgeting cho Chi phí). Tự động tính tỷ lệ Thu nhập Thụ động.
- 🟢 **Life Timeline**: Đã hoàn thành giao diện hiển thị sự kiện (`TimelineView.vue`) và log sự kiện mới có tích hợp AI phân tích tác động tài chính (`LogEventView.vue`). Đã kết nối với Database lưu bảng `life_events`.
- 🟢 **Cascade Engine**: Đã hoàn thành luồng tạo version mới cho Profile. Đã cấu hình Transaction Manager xuyên suốt các Repositories để đảm bảo toàn vẹn dữ liệu khi cascade sự kiện thay đổi vào Profile, Incomes và Dependents.

### Module 1 — Dashboard
- 🟢 **Tổng quan tài sản (Net Worth Card)**: Đã hiển thị số thực tế từ DB.
- 🟢 **Biểu đồ phân bổ tài sản (Pie Chart)**: Đã hoàn thành (Sử dụng Chart.js).
- 🟢 **Tiến trình FI**: Đã hiển thị thanh tiến trình dựa trên mục tiêu hiện tại.
- 🔴 **Biểu đồ tăng trưởng (Line Chart)**: Đang chờ tính năng "Chốt sổ tháng" (Monthly Review) để có dữ liệu lịch sử vẽ biểu đồ đường.
- 🟢 **Cảnh báo chiến lược**: Đã có logic check cảnh báo (Lệch phân bổ tỷ lệ Nợ, tỷ trọng Tiền mặt).

### Module 2 — Assets & Liabilities Management
- 🟢 **Quản lý Tài sản (Assets)**: Đã hoàn thành CRUD, hỗ trợ nhiều phân loại (Vàng, Cổ phiếu, Tiền mặt...). Có tính năng nhập nhanh trọng lượng Vàng, cuộn vô hạn, bộ lọc, sắp xếp.
- 🟢 **Quản lý Khoản nợ (Liabilities)**: Đã hoàn thành CRUD, hỗ trợ nhiều phân loại. Tính năng nhập lãi suất, trả góp tháng.
- 🟢 **Tự động tính toán ròng**: Backend/Frontend tự động đồng bộ.

### Module 3 — Portfolio
- 🟢 **Danh mục đầu tư (Investable Assets)**: Đã hoàn thành API Backend và giao diện Frontend (Hiển thị tỷ trọng thực tế vs mục tiêu IPS).
- 🟢 **Watchlist**: Đã hoàn thành API Backend và giao diện Frontend (Cảnh báo vùng giá mua, xếp hạng độ tự tin).

### Module 4 — Goal-Aware Unified IPS
- 🟢 **Khởi tạo & Phân tích Chiến lược (AI)**: Đã hoàn thành màn hình IPS (`IPSView.vue`). AI có khả năng phân tích dòng tiền, người phụ thuộc và tài sản hiện có để đề xuất tỷ trọng phân bổ (`target_allocation`) và bài văn chiến lược (`detailed_strategy`).
- 🟢 **Chỉnh sửa linh hoạt (Manual Override)**: Cho phép người dùng chỉnh sửa bằng tay Tỷ trọng và Bài văn chiến lược. Có luồng cảnh báo rủi ro khi thay đổi.
- 🟡 **Cảnh báo lệch tỷ trọng (Rebalancing Alerts)**: Đã có DB Notification, nhưng chưa có luồng tự động quét và sinh ra Notification khi tài sản thị trường làm lệch tỷ trọng so với IPS. (Đang chờ)

### Module 5 -> 8 (Thesis, Decision Journal, Monthly Review)
- 🟢 **Investment Thesis**: Đã hoàn thành API lưu trữ, giao diện dạng lưới Grid và tính năng nhờ AI sinh Thesis tự động.
- 🟢 **Monthly Review**: Đã hoàn thành tính năng nhập vốn rải ngân, AI tổng hợp toàn bộ danh mục, IPS, nợ và gợi ý hành động.
- 🔴 **Decision Journal**: Chưa code.
### Module 9 — Authentication & Authorization
- 🟢 **Đăng ký (Signup)**: Đã xong, hỗ trợ multi-user.
- 🟢 **Đăng nhập (Login)**: Đã xong, kết nối DB với JWT Token.
- 🟢 **Bảo mật (JWT Flow)**: API được bảo vệ qua middleware.

---
## Backlog / Future Enhancements (Ghi chú tính năng tương lai)
- 📝 **Real-time Notifications**: Hiện tại Frontend đang dùng cơ chế Polling hoặc Fetch on Action để lấy thông báo. Trong tương lai, cân nhắc nâng cấp lên **Server-Sent Events (SSE)** hoặc **WebSockets** để có cảnh báo thời gian thực ngay khi Backend phân tích xong.

