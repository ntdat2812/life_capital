# Life Capital — Low-Level Design (LLD)

> Tài liệu đặc tả luồng xử lý chi tiết (Low-Level Design) cho các tính năng trọng tâm của hệ điều hành.

---

## 1. Cơ Chế Life Event Cascade Engine

Cascade Engine chịu trách nhiệm đồng bộ trạng thái tài chính khi đời sống của nhà đầu tư có biến động.

```
[Life Event Ghi Nhận] 
         │
         ▼
[AI Phân Tích Tác Động] (Tính toán Delta thu nhập, chi phí, rủi ro)
         │
         ▼
[Tạo Bản Ghi investor_profiles v(N+1)]
         │
         ▼
[Kiểm tra điều kiện IPS] 
         ├── Không có biến động lớn ──> [Kết thúc cascade]
         └── Có biến động đáng kể ───> [Tạo bản nháp investment_policies v(M+1)]
                                                   │
                                                   ▼
                                      [Gửi thông báo User phê duyệt]
                                                   │
                                                   ▼
                                      [User Approve] ──> [Active IPS mới]
                                                               │
                                                               ▼
                                                  [Update target allocations]
                                                               │
                                                               ▼
                                                  [Trigger rebalance portfolio]
```

### Thuật toán xử lý (Go pseudo-code)
```go
func CascadeLifeEvent(ctx context.Context, event model.LifeEvent) error {
    tx, err := db.Begin(ctx)
    defer tx.Rollback()

    // 1. Lấy profile active hiện tại
    currentProfile, err := profileRepo.GetActive(ctx, event.UserID)

    // 2. Gọi AI phân tích Delta tác động
    aiResult, err := aiGateway.AnalyzeLifeEvent(ctx, currentProfile, event)

    // 3. Clone profile và tạo version mới với thay đổi
    newProfile := cloneAndApply(currentProfile, aiResult)
    newProfile.Version = currentProfile.Version + 1
    newProfile.Status = "active"
    newProfile.TriggerEventID = event.ID
    
    // Set old profile status to superseded
    err = profileRepo.UpdateStatus(ctx, currentProfile.ID, "superseded")
    err = profileRepo.Save(ctx, newProfile)

    // 4. Kiểm tra xem có cần cập nhật IPS không
    if aiResult.RequiresIPSUpdate {
        draftIPS := model.InvestmentPolicy{
            UserID: event.UserID,
            Version: currentIPSVersion + 1,
            Status: "draft",
            ProfileVersion: newProfile.Version,
            TriggerEventID: event.ID,
            TargetAllocation: aiResult.SuggestedAllocation,
            AIRationale: aiResult.Rationale,
        }
        err = ipsRepo.Save(ctx, draftIPS)
        // Kích hoạt thông báo đẩy lên UI
    }

    return tx.Commit()
}
```

---

## 2. Quản lý Tài sản & Nợ (Assets & Liabilities)

### 2.1. Thiết kế dữ liệu Tài sản biến động
Để hỗ trợ việc tích hợp lấy giá Real-time từ API bên thứ 3 trong tương lai (Phase 2):
- Các tài sản dạng biến động (Cổ phiếu, Quỹ, Vàng, Crypto) được cấu trúc để người dùng nhập `quantity` (Số lượng) và `ticker` (Mã tài sản).
- Hệ thống hỗ trợ nhập tay `current_price` (Giá thị trường hiện tại) ở Phase 1.
- `current_value` = `quantity` * `current_price` (Tính toán tự động tại giao diện hoặc bằng Job nền trong tương lai).

### 2.2. API & Logic Lấy Dữ Liệu (Pagination & Filtering)
- Backend cung cấp các API `GET /api/v1/wealth/assets` và `GET /api/v1/wealth/liabilities`.
- **Query Parameters**: Hỗ trợ lọc (filtering) theo `category` và phân trang với `limit`, `offset`.
- **Cơ chế Phân trang Frontend**: Sử dụng virtual scrolling (lắng nghe sự kiện `scroll` trên container `max-h`). Dữ liệu sẽ tự động được append (fetch thêm) khi người dùng cuộn xuống gần cuối danh sách mà không cần bấm chuyển trang truyền thống.
- **Tính năng Chỉnh sửa**: Hỗ trợ `PUT /api/v1/wealth/assets/:id` để tái cập nhật thông tin tài sản (Ví dụ: Chỉnh sửa giá mua, thay đổi số lượng, hoặc nhập tổng vốn đầu tư thay cho giá vốn trung bình).

### 2.3. Quản trị đa tiền tệ (Currency)
- Cột `base_currency` mặc định được lưu ở bảng `users` (`VARCHAR(3) DEFAULT 'VND'`). Mọi tính toán tổng tài sản ròng (Net Worth) trên Dashboard và bảng tổng hợp sẽ được quy đổi về `base_currency` này. Các API và giao diện đều sử dụng biến toàn cục này để format số tiền hiển thị.

---

## 3. Chu kỳ Monthly Review Loop

Mỗi tháng một lần, hệ thống thực hiện phân tích tổng thể:

```
[Bắt đầu chu kỳ] ──> [Life Event Check] (User log/xác nhận các biến cố)
                            │
                            ▼
                     [Nhập thông tin dòng tiền] (Capital mới, thu nhập/chi phí)
                            │
                            ▼
                     [Import/Cập nhật số dư tài sản] (Assets & Holdings)
                            │
                            ▼
                     [AI Synthesis Processing]
                            │
                            ▼
                     [Trả kết quả Báo Cáo AI] (Phân tích IPS, thesis, valuation)
                            │
                            ▼
                     [Đề Xuất Phân Bổ Vốn] (Ví dụ: FPT 12M, MBB 3M)
                            │
                            ▼
                     [User chấp nhận đề xuất] ──> [Tự động log vào Decision Journal]
```

### Quy trình tổng hợp Prompt cho AI Monthly Review:
1. **Thu thập dữ liệu**: Backend query toàn bộ thông tin:
   - Profile hiện tại (Income, Expense, Risk Score).
   - Danh mục hiện tại (Holdings, giá vốn, tỷ trọng thực tế).
   - Mục tiêu IPS (Target Allocation, quy tắc).
   - Luận điểm đầu tư đang active của các mã đang hold.
   - Kết quả kinh doanh kỳ mới nhất (Earnings Review).
   - Watchlist (Mức định giá, target price).
2. **Xây dựng Context**: Kết hợp các JSON thành một payload có cấu trúc.
3. **Gọi LLM (Thông qua AI Provider Pattern hỗ trợ nhiều models khác nhau)**: Yêu cầu đưa ra quyết định giải ngân tối ưu cho `new_investment_amount` dựa trên các quy tắc cứng của IPS và logic mềm của Thesis.

---

## 3. Quy trình Đăng Nhập & Phân Quyền (Authentication)

Hệ thống hỗ trợ nhiều người dùng (multi-user) và được bảo mật bằng cơ chế chuẩn:
- **Đăng ký/Đăng nhập**: Sử dụng mật khẩu băm thông qua thuật toán `bcrypt` với cost là `12` lưu tại bảng `users`.
- **Session**: Cấp phát mã JWT (JSON Web Token) chứa `user_id` và thời gian hết hạn (ExpiresAt: 72 giờ).
- **Middleware**: Echo router sử dụng `middleware.JWTWithConfig` để chặn tất cả các request ngoại trừ `/api/v1/auth/login` và `/api/v1/onboarding/start`.
- **Cơ sở dữ liệu**: Tất cả các truy vấn đều filter theo `user_id` được trích xuất từ JWT context để bảo vệ dữ liệu.

---

## 4. Thiết Kế Lưu Trữ Suggestion & Luồng Load Context của AI

Để đảm bảo hiệu năng và tính nhất quán của dữ liệu, quá trình load context và lưu các khuyến nghị của AI được thiết kế theo các pha trạng thái cụ thể:

### 4.1 Cơ chế Load Context lên AI (Context Aggregation)
Trước khi gửi Prompt tới LLM (sử dụng AI Provider Pattern hỗ trợ linh hoạt các LLM Models), backend thực hiện gom cụm dữ liệu đồng thời bằng `goroutine` để tối ưu hóa thời gian phản hồi:

```
                  ┌───> GetActiveProfile() ──────┐
                  ├───> GetCurrentHoldings() ────┤
[Request Trigger] ├───> GetActiveIPS() ──────────┼───> [Aggregate into Go Struct] ───> [Serialize JSON Context] ───> [Send LLM]
                  ├───> GetActiveTheses() ───────┤
                  ├───> GetNewEarningsReviews() ─┤
                  └───> GetWatchlistItems() ─────┘
```

1. **Context Struct (Go)**:
   ```go
   type AIReviewContext struct {
       Profile           model.InvestorProfile   `json:"profile"`
       Holdings          []model.Holding         `json:"holdings"`
       Policy            model.InvestmentPolicy  `json:"policy"`
       Theses            []model.Thesis          `json:"theses"`
       RecentEarnings    []model.EarningsReview  `json:"recentEarnings"`
       Watchlist         []model.WatchlistItem   `json:"watchlist"`
       Liabilities       []model.Liability       `json:"liabilities"`
       Goals             []model.Goal            `json:"goals"`
       MonthlyInput      model.MonthlyInput      `json:"monthlyInput"`
   }
   ```
2. **Serializing**: Struct này được serialize sang định dạng JSON sạch và tiêm trực tiếp vào system instructions làm dữ liệu nền tảng.

### 4.2 Nơi lưu trữ đề xuất của AI (Lifecycle & Persistence)

Dữ liệu do AI sinh ra (khuyến nghị mua/bán, phân tích tác động rủi ro) được quản lý qua 3 trạng thái:

```
[1. Draft State]             [2. User Decision State]          [3. Execution State]
AI sinh khuyến nghị           User duyệt qua UI                 Hệ thống cập nhật danh mục
Lưu JSONB vào DB              Accepted / Rejected JSONB         Ghi Decision Journal + Asset
```

#### Trạng thái 1: Draft (Nháp)
- Khi AI vừa hoàn tất phân tích, dữ liệu thô dạng cấu trúc được lưu trực tiếp vào cột JSONB `ai_recommendations` trong bảng `monthly_reviews`.
- Các nội dung văn bản phân tích (ví dụ: `ai_net_worth_analysis`, `ai_portfolio_analysis`) được lưu vào các cột TEXT tương ứng trong bảng `monthly_reviews` để hiển thị báo cáo mà không cần gọi lại API AI.

#### Trạng thái 2: User Decision (Ra quyết định)
- Trên giao diện Vue, người dùng check `[Accept]` hoặc `[Reject]` cho từng đề xuất giải ngân (Ví dụ: Chấp nhận mua FPT 12M, Từ chối mua REE 3M).
- Lựa chọn này gửi về backend cập nhật vào 2 trường JSONB: `accepted_recommendations` và `rejected_recommendations` của bảng `monthly_reviews` để làm dữ liệu lưu trữ đối chiếu.

#### Trạng thái 3: Execution (Thực thi)
- Với các đề xuất được `Accept`, backend tự động thực hiện:
  1. Tạo hoặc cập nhật giá trị tài sản trong bảng `assets` và danh mục `portfolio_holdings`.
  2. Ghi một bản ghi vào bảng nhật ký quyết định `decision_journal` với `decision_type = 'buy'` hoặc `'sell'`, tự động liên kết khóa ngoại tới `monthly_review_id` để tiện cho việc đánh giá lại (Decision Review) sau này.

#### Trạng thái 4: Raw Logs (Nhật ký AI)
- Toàn bộ Prompt gửi đi và Response thô trả về kèm theo thông tin model name, số token sử dụng được lưu vào bảng `ai_conversations` với `context_type = 'monthly_review'` nhằm mục đích debug và tối ưu prompts.

---

## 5. Thiết Kế Mục Tiêu Tài Chính (Goal-Based Investing Waterfall)

Hệ thống không quản lý mục tiêu tĩnh mà mô hình hóa chúng thành cơ chế **Goal Waterfall** động. Toàn bộ dòng tài sản ròng (Net Worth) được rót chảy lần lượt qua các mục tiêu tài chính dựa trên mức độ ưu tiên:

### 5.1 Quy trình phân phối Net Worth theo thứ tự ưu tiên (Waterfall Algorithm)
1. **Sắp xếp mục tiêu**: Backend/Frontend sắp xếp danh sách `financial_goals` theo mức độ ưu tiên tăng dần (`priority` từ 1 đến 5).
2. **Rót vốn lũy kế**:
   - Khởi tạo `RemainingNetWorth = TotalAssets - TotalLiabilities`.
   - Duyệt qua từng mục tiêu:
     - Số tiền phân bổ cho mục tiêu: `Allocated = Min(RemainingNetWorth, Goal.TargetAmount)`.
     - Tiến độ mục tiêu (%): `Progress = (Allocated / Goal.TargetAmount) * 100`.
     - Trừ số tiền đã phân bổ: `RemainingNetWorth = RemainingNetWorth - Allocated`.
3. **Phản ứng Realtime**: Khi người dùng CRUD tài sản hoặc khoản nợ, Net Worth thay đổi sẽ lập tức trigger chạy lại thuật toán Waterfall này, cập nhật trạng thái tiến độ của toàn bộ mục tiêu mà không cần sửa thủ công.

---

## 6. Mô Hình Hóa Khoản Nợ & Tác Động Dòng Tiền (Debt Cashflow & Liabilities)

Sự hiện diện của các khoản nợ trong bảng cân đối kế toán ảnh hưởng sâu sắc đến dòng tiền và tốc độ tích sản:

### 6.1 Nghĩa vụ trả nợ định kỳ (Debt Service Impact)
- **Công thức tính dòng tích lũy ròng định kỳ**:
  ```
  AdjustedMonthlySavings = (TotalIncome * BaseSavingsRate) - Sum(Liabilities.MonthlyPayment)
  ```
- **Tác động Timeline**: Nếu nghĩa vụ trả nợ định kỳ chiếm dụng phần lớn dòng tiền tiết kiệm, `AdjustedMonthlySavings` giảm mạnh ➡️ Số năm ước tính để Net Worth đạt mức mục tiêu tự do tài chính (FI Target) tự động kéo dài ra.
- **AI Warning Triggers**:
  - `Debt-to-Asset Ratio = TotalLiabilities / TotalAssets`.
  - Nếu tỷ số đòn bẩy vượt quá ngưỡng trần rủi ro định sẵn (mặc định **35%**):
    1. AI Committee tự động kích hoạt trạng thái Cảnh báo rủi ro đòn bẩy cao.
    2. Điều chỉnh thuật toán gợi ý rebalance: Đề xuất giải ngân mới của tháng sẽ dồn **80-100% dòng tiền dư** để ưu tiên trả nợ gốc trước khi giải ngân vào các lớp tài sản tăng trưởng như Cổ phiếu.

---

## 7. Luồng Xác Thực (Authentication & JWT Flow)

### 7.1 Đăng ký (Signup)
```
POST /api/v1/auth/signup
Body: { "email": "...", "name": "...", "password": "..." }

1. Kiểm tra xem email đã được sử dụng hay chưa.
   - Nếu email đã tồn tại → trả về 409 Conflict.
2. Hash mật khẩu bằng bcrypt (cost = 12).
3. INSERT vào bảng `users`.
4. Trả về 201 Created (không trả token — user phải login riêng).
```

### 7.2 Đăng nhập (Login)
```
POST /api/v1/auth/login
Body: { "email": "...", "password": "..." }

1. Tìm user theo email.
   - Không tìm thấy → 401 Unauthorized.
2. So sánh mật khẩu bằng bcrypt.CompareHashAndPassword().
   - Không khớp → 401 Unauthorized.
3. Tạo JWT token chứa claims: { "user_id": uuid, "exp": now + 24h }.
4. Ký token bằng secret key (biến môi trường JWT_SECRET).
5. Trả về 200 OK với body: { "token": "eyJhbG..." }.
```

### 7.3 Middleware Bảo Vệ Route
```go
// Áp dụng cho tất cả route group /api/v1/* (trừ /auth/*)
func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    1. Đọc header Authorization: Bearer <token>.
    2. Parse và validate JWT token bằng JWT_SECRET.
    3. Nếu hợp lệ → inject user_id vào echo.Context → gọi next(c).
    4. Nếu không hợp lệ hoặc hết hạn → trả về 401 Unauthorized.
}
```

### 7.4 Cấu hình Database (Docker Compose)
- **Database Engine**: PostgreSQL 16 (Alpine image).
- **Container Name**: `life_capital_db`.
- **Credentials mặc định (local dev)**:
  - `POSTGRES_DB`: `life_capital`
  - `POSTGRES_USER`: `lifecap`
  - `POSTGRES_PASSWORD`: `lifecap_secret`
- **Port**: `5432:5432`
- **Persistence**: Named volume `pgdata` đảm bảo dữ liệu không mất khi container bị xóa.
