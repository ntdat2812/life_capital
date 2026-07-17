---
name: Backend Developer Skill
description: Go development, Echo framework, SQL database handling, context lifecycle, and error wrapping guidelines.
---

# Backend Developer Skill Guidelines

## 1. Go Idiomatic Development

- **No Magic Frameworks**: Do not hide Go's type-safety. Avoid global variables for states. Pass database pools and clients via dependency injection.
- **Factory Functions**: Implement constructors for all structs:
  ```go
  func NewHandler(svc Service) *Handler {
      return &Handler{service: svc}
  }
  ```

---

## 2. Echo Framework Standards

- **Route Registration**: All routes must be registered via a dedicated `registerRoutes()` function in `main.go`. Handlers are initialized in `main()` and passed into this function:
  ```go
  func registerRoutes(e *echo.Echo, healthH *handler.HealthHandler, profileH *handler.ProfileHandler) {
      api := e.Group("/api/v1")
      api.GET("/health", healthH.HealthCheck)
      api.GET("/profile", profileH.GetProfile)
  }
  ```
- **Handler Files**: Each domain has its own handler file under `internal/handler/` using the naming convention `{domain}_handler.go` (e.g., `health_handler.go`, `profile_handler.go`, `asset_handler.go`). Each handler struct is created with a factory function `NewXxxHandler(svc XxxService)`.
- **Route Grouping**: Define groups logically:
  ```go
  r := e.Group("/api/v1")
  r.Use(middleware.JWTWithConfig(...))
  r.GET("/profile", h.GetProfile)
  r.POST("/life-events", h.CreateLifeEvent)
  ```
- **Context Handling**: Always pass client request context to services and database requests:
  ```go
  ctx := c.Request().Context()
  profile, err := h.service.GetProfile(ctx, userID)
  ```
- **Error Binding**: Bind requests using `c.Bind(payload)`. Do not ignore binding errors. Return `HTTP 400` with clean JSON.


---

## 3. Database Queries using `pgx`

- **Prepared Queries**: Write pure SQL. Avoid raw string interpolation. Use parameter placeholders `$1, $2, ...` to prevent SQL injection.
- **Transactions**: For compound writes (like logging a life event and updating profile versions), run queries in a single PostgreSQL transaction:
  ```go
  tx, err := db.BeginTx(ctx, pgx.TxOptions{})
  defer tx.Rollback(ctx)
  // execute inserts/updates...
  tx.Commit(ctx)
  ```

---

## 4. Error Mapping

Map internal error variables to Echo HTTP errors:

```go
var (
    ErrNotFound     = errors.New("resource not found")
    ErrUnauthorized = errors.New("unauthorized access")
    ErrConflict     = errors.New("data state conflict")
)

func MapHTTPError(err error) error {
    if errors.Is(err, ErrNotFound) {
        return echo.NewHTTPError(http.StatusNotFound, err.Error())
    }
    if errors.Is(err, ErrConflict) {
        return echo.NewHTTPError(http.StatusConflict, err.Error())
    }
    return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
}
```

---

## 5. Swagger API Documentation Standards

Mọi Handler endpoint trong Go Echo phải sử dụng cú pháp chú thích (annotations) chuẩn của **swag** để tự động tạo Swagger docs:

```go
// GetProfile handles retrieving the active investor profile.
// @Summary      Lấy thông tin hồ sơ nhà đầu tư
// @Description  Trả về hồ sơ nhà đầu tư đang active cùng các thông số tài chính cơ bản.
// @Tags         profile
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.InvestorProfile
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      404  {object}  map[string]string "Not Found"
// @Router       /api/v1/profile [get]
func (h *Handler) GetProfile(c echo.Context) error { ... }
```

### Các bước biên dịch tài liệu Swagger:
1. Mỗi khi thay đổi hoặc thêm endpoint mới, chạy lệnh sau ở thư mục `backend/` để sinh lại folder `docs/`:
   ```bash
   swag init -g cmd/server/main.go
   ```
2. Thư mục `docs` được tự động sinh ra và phải được import dạng blank import (`_ "github.com/datnguyen/life_capital/backend/docs"`) trong `main.go` để đăng ký tài liệu.

---

## 6. Build & Run

- Main entrypoint is `backend/cmd/server/main.go`.
- Configuration values should be loaded from environment variables or a `.env` file using a library like `github.com/joho/godotenv`.


---

## 6. References & Design Documents
When writing backend code or setting up database tables, you MUST read:
- **Low-Level Logic Models**: [docs/lld.md](file:///Users/datnguyen/Dev/Projects/Go/life_capital/docs/lld.md)
- **Database Schema Specs**: [docs/erd.md](file:///Users/datnguyen/Dev/Projects/Go/life_capital/docs/erd.md)
- **Features & API Endpoints**: [docs/features.md](file:///Users/datnguyen/Dev/Projects/Go/life_capital/docs/features.md)

