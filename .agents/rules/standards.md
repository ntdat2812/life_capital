---
name: Life Capital Engineering Standards
description: Canonical coding, naming, technology, API, database, and testing standards for the Life Capital project.
---

# Life Capital Engineering Standards

> Single source of truth for **how** we build Life Capital (WealthOS).

---

## 1. Technology Stack

| Layer | Choice | Notes |
|---|---|---|
| **Backend Language** | **Go 1.22+** | Clean concurrency, explicit typing, simple deployment |
| **Backend Framework** | **Echo (v4)** | Lightweight HTTP router with robust middleware |
| **Frontend Framework** | **Vue 3 (Vite)** | Composition API (`<script setup>`) with Pinia for state management |
| **Frontend Routing** | **Vue Router** | Standard router, history mode |
| **Styling** | **Tailwind CSS 4 + custom Vanilla CSS** | Clean glassmorphism & dark theme |
| **Database** | **PostgreSQL 16** | Strong relational integrity + JSONB support |
| **AI Integration** | **Gemini 2.5 Flash API** | Prioritize free tiers, low-cost API calls, or local models via Ollama. Allow key injection. |
| **Authentication** | **JWT (Header-based) + Bcrypt** | Simple login flow, single-user deployment |

---

## 2. Coding Standards

### Backend (Go)
- **Layered Architecture**: `handler -> service -> repository`.
- **Constructors**: Provide explicit factory functions (`NewService`, `NewRepository`) to inject dependencies.
- **Context Propagation**: Always pass `context.Context` to repositories and services for timeout and cancellation.
- **Error Handling**: Use explicit error returning. Map core business/domain errors (e.g. `ErrNotFound`, `ErrConflict`, `ErrInvalidInput`) to HTTP status codes in the handler or custom Echo error handler.
- **Panic Recovery**: Ensure the Echo recovery middleware is active. Never allow application panics in handlers.

### Frontend (Vue 3)
- **Composition API**: Leverage `<script setup>` style for consistency.
- **State Management**: Keep API calls inside Pinia actions. Components should only read state and call actions.
- **Component Design**: Modular, reusable components under `src/components/`, layout templates under `src/views/`.
- **Vanilla CSS / Custom Transitions**: Maximize CSS transitions for a premium, buttery-smooth experience (200-300ms transitions).

---

## 3. Naming Standards

### Go Naming
- **File Names**: snake_case (e.g., `user_handler.go`, `portfolio_repository.go`).
- **Structs/Interfaces**: PascalCase (e.g., `PortfolioService`, `AssetRepository`).
- **Functions/Methods**: PascalCase (e.g., `GetAssetByTicker`, `LogLifeEvent`).
- **Variables**: camelCase (e.g., `currentNetWorth`, `newInvestmentAmount`).

### Vue Naming
- **Components**: PascalCase (e.g., `MetricCard.vue`, `LifeTimeline.vue`).
- **Views**: PascalCase (e.g., `DashboardView.vue`, `PortfolioView.vue`).
- **Stores**: camelCase with suffix `Store` (e.g., `authStore.js`, `profileStore.js`).
- **Files/Assets**: kebab-case or camelCase as appropriate.

---

## 4. API Design Standards

- **RESTful endpoints** under `/api/v1/`.
- **JSON bodies**: CamelCase JSON fields to match Vue client preferences (use struct tags: `json:"netWorth"`).
- **HTTP Verbs**: 
  - `GET` for reads.
  - `POST` for resource creation (e.g., `/api/v1/life-events`).
  - `PUT`/`PATCH` for updates.
  - `DELETE` for removals/archivals.
- **HTTP Status Codes**:
  - `200 OK` / `201 Created`.
  - `400 Bad Request` for validation errors.
  - `401 Unauthorized` for expired/invalid tokens.
  - `403 Forbidden` for invalid access levels.
  - `404 Not Found` for resource missing.
  - `409 Conflict` for state violations (e.g., target allocations exceeding 100%).
  - `422 Unprocessable Entity` for business rules checks.
- **Swagger Documentation**: Mọi API endpoint mới được tạo phải có chú thích Swagger chi tiết (Summary, Description, Tags, Param, Success, Failure, Router). Tài liệu API phải được render tự động qua Swagger UI tại đường dẫn `/swagger/*`.

---

## 5. Testing & Definition of Done

A task is **DONE** only when:
1. **Architecture check**: Business logic is inside services, route bindings in handlers, query commands in repository.
2. **Schema and Cascade Rules validation**: New events trigger the correct updates via the `Life Event Cascade Engine`.
3. **Frontend compatibility**: Vue store matches API payload structure.
4. **Clean codebase**: Code compiles without errors, no unused variables or deprecated types, code is formatted with `gofmt` and CSS is clean.
5. **Single User Security check**: Endpoint requires a JWT token (except onboarding and login).
