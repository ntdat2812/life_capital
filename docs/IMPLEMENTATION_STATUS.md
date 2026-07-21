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
- 🔴 **Danh mục cổ phiếu**: Chưa code.
- 🔴 **Khuyến nghị & Conviction**: Chưa code.

### Module 4 -> 8 (IPS, Thesis, Decision Journal, Monthly Review)
- 🔴 **Tất cả**: Chưa code.

### Module 9 — Authentication & Authorization
- 🟢 **Đăng ký (Signup)**: Đã xong, hỗ trợ multi-user.
- 🟢 **Đăng nhập (Login)**: Đã xong, kết nối DB với JWT Token.
- 🟢 **Bảo mật (JWT Flow)**: API được bảo vệ qua middleware.

