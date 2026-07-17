---
name: UX Designer Skill
description: Core visual principles, glassmorphic themes, UI assets styling, component layout guidelines, responsive grids, and micro-animation specifications.
---

# UX/UI Designer Skill Guidelines

## 1. Visual & Aesthetic Principles
Life Capital requires a premium, state-of-the-art visual appearance that looks and feels like a modern SaaS platform (WealthOS). Do not design basic, generic, or simple templates.

---

## 2. Design Tokens & Styling Rules

### 2.1 Color Palette
- **Base Background**: Slate 900 (`#0F172A`)
- **Card/Surface Background**: Slate 800 (`#1E293B`) with partial opacity (`rgba(30, 41, 59, 0.7)`)
- **Accent Elements**: Indigo 500 (`#6366F1`) or Violet 500 (`#8B5CF6`) for actions, buttons, highlights.
- **Alert & Growth Indicators**:
  - Emerald 500 (`#10B981`) for positive returns, targets completed, or healthy states.
  - Amber 500 (`#F59E0B`) for warnings (e.g., target allocation drift).
  - Red 500 (`#EF4444`) for critical errors, broken theses, or risks.

### 2.2 Glassmorphism & Translucency
- Use translucent backgrounds and borders to create depth:
  ```css
  background: rgba(30, 41, 59, 0.7);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.05);
  box-shadow: 0 8px 32px 0 rgba(0, 0, 0, 0.37);
  ```

### 2.3 Typography Hierarchy
- **Header Fonts**: Outfit or Inter (bold, uppercase/letter-spaced tracking for secondary labels).
- **Body & Data**: Inter for text, JetBrains Mono (or a clean sans-serif mono font) for all currency values, allocations, and percentages to align values vertically.

---

## 3. Premium Layout & Component Rules

- **Information Densities**: Group information using cards. Avoid endless list pages.
- **Profile Timeline**: Display Life Timeline with custom chronological icons, connecting vertical paths, and clear hover cards detailing profile version drifts.
- **Wizard Interactivity**: The 8-step Onboarding Interview and 4-step Monthly Review must have a conversational message bubbles interface, including a smooth vertical typing indicator and animated transition items.
- **Active Alerts Panel**: Placed top-right or prominently on the dashboard. Uses soft gradient warning borders instead of harsh solid colors.

---

## 4. Micro-animations & Transitions

- Apply transitions to all states: `hover`, `active`, `focus`.
- **Standard Transition Timing**: `all 200ms cubic-bezier(0.4, 0, 0.2, 1)`.
- Use a slight scale effect (`scale-102`) or shadow lift for cards when hovered.

---

## 5. References & Design Documents
When creating mockups, templates, or styles, you MUST refer to:
- **UI Layout Specs**: [docs/frontend_layout.md](file:///Users/datnguyen/Dev/Projects/Go/life_capital/docs/frontend_layout.md)
- **Features & Screen Routes**: [docs/features.md](file:///Users/datnguyen/Dev/Projects/Go/life_capital/docs/features.md)

