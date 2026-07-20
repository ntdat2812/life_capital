# Welcome to Life Capital (WealthOS)

> This document is the primary entry point for all AI agents working on this project. It outlines the codebase layout, architectural decisions, coding style, and workflow standards.

---

## 1. Project Overview

**Life Capital** (also known as **WealthOS**) is a Personal Wealth Operating System.
- **Core Goal**: Help the user answer: *"With all my current assets, what should I do this month to get closer to financial independence?"*
- **Scope**: Multi-user, self-hosted, focused on monthly analysis rhythm rather than daily trading.
- **Key Paradigm**: Module 0 - Investor Profile & Life Timeline is the foundational layer. Every recommendation, risk limit, or target allocation must start from the current profile status and cascade when Life Events are logged.

---

## 2. Tech Stack

- **Backend**: Go (Golang) 1.22+ using **Echo Framework**.
- **Frontend**: Vue 3 (Vite, Pinia, Vue Router) using **Tailwind CSS 4** and Vanilla CSS components.
- **Database**: PostgreSQL 16 (relational data, JSONB for flexible attributes).
- **AI Integration**: Prioritizes low-cost/free tiers (e.g., Gemini 2.5 Flash, local models via Ollama) and supports custom user API keys.
- **Authentication**: JWT-based authentication with bcrypt hashed password. Multi-user signup/login dashboard.

---

## 3. Project Structure

```
life_capital/
├── frontend/                          # Vue 3 App (Vite)
│   ├── src/
│   │   ├── assets/                    # Styling, fonts, icons
│   │   ├── components/                # Reusable UI components
│   │   ├── router/                    # Vue Router configuration
│   │   ├── stores/                    # Pinia stores (auth, profile, assets, portfolio)
│   │   ├── views/                     # Screen views (Dashboard, Profile, Portfolio, etc.)
│   │   ├── App.vue
│   │   └── main.js
│   ├── package.json
│   └── vite.config.js
│
├── backend/                           # Go Backend (Echo)
│   ├── cmd/
│   │   └── server/
│   │       └── main.go                # Application Entry Point & Route Registration
│   ├── internal/
│   │   ├── handler/                   # Echo HTTP Controllers (one file per domain: {domain}_handler.go)
│   │   ├── service/                   # Business Logic & Cascade Engine
│   │   ├── repository/                # Database Queries (SQL / pgx)
│   │   ├── model/                     # Domain & DB Entities
│   │   ├── ai/                        # Prompts, AI Client Wrapper
│   │   ├── config/                    # Config structures & env loading
│   │   └── middleware/                # Auth, CORS, Logger
│   ├── docs/                          # Auto-generated Swagger API docs (swag init)
│   ├── migrations/                    # SQL migration scripts
│   ├── go.mod
│   └── go.sum
│
├── .agents/                           # AI Agent specific rules and skills
│   ├── rules/
│   │   └── standards.md               # Technology, coding, and naming conventions
│   └── skills/
│       ├── business_analyst/          # Product design and requirements guidelines
│       ├── solution_architect/        # Architecture, DB schema, and API standards
│       ├── ux_designer/               # Styling tokens, responsive grid, visual design rules
│       ├── backend_developer/         # Go, SQL, Echo server build patterns
│       ├── frontend_developer/        # Vue 3, Pinia, CSS styles guidelines
│       └── qa_tester/                 # Automated, manual, and visual interface checks
│
├── docs/                              # Project Design Documentation (Modular)
│   ├── business.md                    # Business vision and flows
│   ├── erd.md                         # SQL Schema & DB architecture
│   ├── lld.md                         # Low level design of key features
│   ├── frontend_layout.md             # UI layout details & templates
│   └── features.md                    # Breakdown of system features
│
├── docker-compose.yml
└── README.md
```

---

## 4. Coding Conventions Quick Summary

1. **Pure Separation of Concerns**: Controllers only handle requests/responses. Services manage all business logic, particularly the Life Event Cascade. Repositories manage database operations.
2. **Explicit Layering**: Go backend uses strict dependency injection (`Handler -> Service -> Repository`).
3. **Database Integrity**: PostgreSQL handles relationship constraints. JSONB is used only for unstructured metadata (e.g., raw AI conversation logs, custom constraints).
4. **Vue 3 Best Practices**: Use Composition API (`<script setup>`), Pinia for global state, and keep components modular.
5. **No fully qualified imports in Go**: Always use standard imports.
6. **Robust Error Handling**: Go returns clean, typed domain errors mapped to HTTP status codes inside middleware handlers.
7. **Swagger API Documentation**: Every API endpoint MUST have Swagger annotations (`@Summary`, `@Description`, `@Tags`, `@Router`, etc.). Run `swag init -g cmd/server/main.go` in `backend/` after any API changes. Swagger UI is served at `/swagger/*`.

---

## 5. Active Agent Guidelines

- **Mandatory Read Step**: BEFORE making any edits or starting a task, you MUST check the `docs/` folder. Read the corresponding design document based on your active role:
  - **Business Analyst**: Read [docs/business.md](file:///Users/datnguyen/Dev/Projects/Go/life_capital/docs/business.md) & [docs/features.md](file:///Users/datnguyen/Dev/Projects/Go/life_capital/docs/features.md).
  - **Solution Architect**: Read [docs/erd.md](file:///Users/datnguyen/Dev/Projects/Go/life_capital/docs/erd.md) & [docs/lld.md](file:///Users/datnguyen/Dev/Projects/Go/life_capital/docs/lld.md).
  - **Developer (Backend/Frontend)**: Read [docs/lld.md](file:///Users/datnguyen/Dev/Projects/Go/life_capital/docs/lld.md) & [docs/frontend_layout.md](file:///Users/datnguyen/Dev/Projects/Go/life_capital/docs/frontend_layout.md).
  - **QA Tester**: Read [docs/features.md](file:///Users/datnguyen/Dev/Projects/Go/life_capital/docs/features.md) & [docs/lld.md](file:///Users/datnguyen/Dev/Projects/Go/life_capital/docs/lld.md).
- Refer to `.agents/rules/standards.md` before making any coding modifications. However, **do not let the documents strictly bind your thinking**; you are encouraged to be creative and propose innovative solutions that align with the product's vision.
- Consult the specific roles under `.agents/skills/` depending on the current task (e.g. database change -> Solution Architect; Go code change -> Backend Developer).
- **Project Documentation Updates**: Every design or database changes must be reflected in the relevant files inside the `docs/` directory. Solution Architect updates architecture/DB schemas, QA updates tests, etc.
- **Swagger Sync**: After adding or modifying any backend API endpoint, regenerate Swagger docs by running `swag init -g cmd/server/main.go` inside `backend/` and commit the updated `backend/docs/` folder.
- Always verify changes with unit tests and document updates in the `/walkthrough.md` file.
