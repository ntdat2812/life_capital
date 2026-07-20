---
name: Frontend Developer Skill
description: Vue 3 Composition API, Pinia state stores, glassmorphism, transitions, layout layouts, and API client setup.
---

# Frontend Developer (Vue) Skill Guidelines

## 1. Vue 3 Composition API & `<script setup>`

- Always write stores and components using the Vue 3 Composition API.
- **Language**: All user-facing text in the UI must be in Vietnamese. You can use English for internal code logic, variables, and components.
- Use explicit TypeScript annotations or reactive objects:
  ```html
  <script setup>
  import { ref, onMounted } from 'vue'
  import { useProfileStore } from '@/stores/profile'

  const profileStore = useProfileStore()
  const loading = ref(true)

  onMounted(async () => {
    await profileStore.fetchProfile()
    loading.value = false
  })
  </script>
  ```

---

## 2. Pinia State Design

- Keep API interactions inside stores (separation of concerns).
- Store states must represent clear business entities:
  - `authStore.js`: authentication token, user session.
  - `profileStore.js`: active investor profile, life timeline list.
  - `assetStore.js`: current assets categories, cash value tracking.

---

## 3. Premium Glassmorphic Styling Guidelines

Life Capital demands a modern, state-of-the-art visual appearance:
- **Colors**: Sleek dark theme using Slate (`#0F172A`), Indigo (`#6366F1`), and Emerald (`#10B981`) for positive growth signals.
- **Glassmorphism**: Use backdrop filters for surfaces:
  ```css
  .glass-card {
    background: rgba(30, 41, 59, 0.7);
    backdrop-filter: blur(12px);
    border: 1px solid rgba(255, 255, 255, 0.05);
    border-radius: 12px;
  }
  ```
- **Typography**: Inter for primary text, monospace or JetBrains Mono for financial figures.
- **Animations**: Transitions should be subtle and clean (e.g., `transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1)`).

---

## 4. API Client Integration

- Use Axios or native fetch wrapped in a client interceptor file (`src/lib/api.js`).
- Automatically inject the JWT token from `authStore` into the `Authorization` header.
- Cleanly intercept `401 Unauthorized` responses to redirect the user to `/login`.

---

## 5. References & Design Documents
When building views, templates, or styles, you MUST refer to the following documents. However, **do not let them strictly bind your thinking**. You are encouraged to be creative and propose innovative solutions that elevate the product's quality according to your specific frontend expertise:
- **Layout & CSS Tokens**: [docs/frontend_layout.md](file:///Users/datnguyen/Dev/Projects/Go/life_capital/docs/frontend_layout.md)
- **Screens & Routing Grid**: [docs/features.md](file:///Users/datnguyen/Dev/Projects/Go/life_capital/docs/features.md)
- **Onboarding/Cascade Logic**: [docs/lld.md](file:///Users/datnguyen/Dev/Projects/Go/life_capital/docs/lld.md)

