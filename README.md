# Life Capital (WealthOS) 🚀

Life Capital (WealthOS) là một **Hệ điều hành Quản lý Tài sản Cá nhân (Personal Wealth Operating System)**. Dự án được thiết kế để giúp người dùng theo dõi và hoạch định bức tranh tài chính cá nhân một cách toàn diện, thông minh, với sự hỗ trợ đắc lực từ AI (Gemini / Llama 3).

## 🌐 Live Deployments

Dự án hiện đang được deploy và hoạt động thực tế trên cloud:

* **🖥️ Frontend (Web App):** [https://life-capital-nine.vercel.app/](https://life-capital-nine.vercel.app/)
* **⚙️ Backend API (Swagger Docs):** [https://life-capital.onrender.com/swagger/index.html](https://life-capital.onrender.com/swagger/index.html)

*(Lưu ý: Do backend sử dụng gói miễn phí của Render, nếu bạn truy cập lần đầu tiên sau một khoảng thời gian dài không sử dụng, server có thể sẽ mất khoảng 30-50 giây để khởi động lại. Hãy kiên nhẫn tải lại trang nhé!)*

## 🛠 Tech Stack

* **Frontend:** Vue 3 (Composition API), Vite, Pinia, Tailwind CSS.
* **Backend:** Go (Golang 1.22+), Echo Framework, pgx.
* **Database:** PostgreSQL.
* **AI Integration:** Google Gemini & Groq (Llama 3).

## 📖 Tính năng chính (Key Features)

- **Hồ sơ đầu tư:** Xác định khẩu vị rủi ro và mục tiêu tài chính cá nhân.
- **Theo dõi Dòng tiền & Tài sản:** Quản lý Thu nhập, Người phụ thuộc, Tài sản (Assets) và Tiêu sản/Nợ (Liabilities).
- **Phân tích Sự kiện cuộc sống (Life Events):** AI tự động phân tích và tính toán lại kế hoạch tài chính mỗi khi bạn gặp một biến cố hay sự kiện mới trong đời.
- **Chính sách Đầu tư (IPS) & Review hàng tháng:** Đưa ra các khuyến nghị, quy tắc đầu tư và tạo báo cáo cập nhật hàng tháng.
