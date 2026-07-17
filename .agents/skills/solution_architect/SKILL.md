---
name: Solution Architect Skill
description: System architecture, database design patterns, API design paradigms, and AI prompt orchestration guidelines.
---

# Solution Architect (SA) Skill Guidelines

## 1. System Architecture Diagram

```
[ Vue 3 Client ] ---> [ Echo HTTP Router ] ---> [ Database Repository (pgx) ] ---> [ PostgreSQL 16 ]
                             |
                             v
               [ Life Event Cascade Engine ]
                             |
                             v
                 [ AI Gateway Service ] ---> [ Gemini 2.5 Flash / Ollama ]
```

---

## 2. Database Schema Guidelines

### 2.1 Schema Versioning
- Use standard SQL migrations inside `backend/migrations/`.
- All tables must explicitly map foreign key constraints and define cascades appropriately.
- PostgreSQL JSONB is used for structured metadata (e.g., `life_constraints` in `investor_profiles`, `ai_recommendations` in `monthly_reviews`). Do not use JSONB for primary business entities.

### 2.2 Temporal & Versioning Design
- **Version Columns**: `investor_profiles` and `investment_policies` tables use integer incrementing version columns (`version`).
- **Trigger event links**: `investor_profiles` must store `trigger_event_id` to link back to the `life_events` table record.

---

## 3. Life Event Cascade Engine Logic

When a life event is logged:
1. **Save Event**: Persist the `life_event` record.
2. **AI Impact Analysis**: Fire AI request mapping the event payload against the latest `investor_profile`.
3. **Trigger Profile Update**: Create a new version of `investor_profile` containing modified cash flow, updated risk scores, and new constraints.
4. **IPS Check**: If the profile changes risk tolerance or income significantly, flag `requires_ips_update = true` and create a draft of `investment_policy` for user review.

---

## 4. API Design Standards

- Route endpoints explicitly. Use Echo sub-routing for resource groups (e.g., `v1 := e.Group("/api/v1")`).
- Validate incoming JSON payloads using Go's `go-playground/validator` integrated into Echo context binders.
- Response payloads must be enveloped in consistent schemas.

---

## 5. AI Prompt Orchestration

- Prompts should use a multi-step design. Construct the prompt dynamically by combining the current profile context, existing holdings, and watchlist parameters.
- Provide JSON Schemas in system instructions to guarantee type-safe JSON returns from LLM models.
- Always implement a fallback model wrapper (e.g. Gemini 2.5 Flash as primary, Local Ollama or Claude as fallback).

---

## 6. Documentation & Architecture Sync

Every architectural or database schema change must be documented synchronously:
- **DB Schema/Index Changes**: Update [erd.md](file:///Users/datnguyen/Dev/Projects/Go/life_capital/docs/erd.md) (`docs/erd.md`) to reflect tables and fields before applying migrations.
- **Workflow/Cascade Logic Changes**: Update [lld.md](file:///Users/datnguyen/Dev/Projects/Go/life_capital/docs/lld.md) (`docs/lld.md`) to keep logic blocks and Go pseudocode models accurate.
- **Core Requirements & Calculations**: Update [business.md](file:///Users/datnguyen/Dev/Projects/Go/life_capital/docs/business.md) (`docs/business.md`) if financial rules or product scopes evolve.
- **API Spec Changes**: Update the API Design section in [features.md](file:///Users/datnguyen/Dev/Projects/Go/life_capital/docs/features.md) (`docs/features.md`) và đồng bộ trực tiếp chú thích Swagger trong mã nguồn Go.


