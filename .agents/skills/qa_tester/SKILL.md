---
name: QA Tester Skill
description: Testing standards, Go backend unit and integration test design, cascade engine validation, and automated checks.
---

# QA Tester Skill Guidelines

## 1. Go Unit Testing Standards

- Use the standard `testing` package with assertions from `github.com/stretchr/testify/assert`.
- Group test suites by context:
  ```go
  func TestCalculateFITimeLine(t *testing.T) {
      profile := Model.InvestorProfile{
          PrimaryIncome: 45000000.0,
          TotalMonthlyExpense: 20000000.0,
          FITargetAmount: 8000000000.0,
      }
      years := CalculateYearsToFI(profile)
      assert.Equal(t, 26.6, years) // simple static check example
  }
  ```

---

## 2. Cascade Engine Validation Flow

When testing `Life Event Cascades`, assert:
1. The cascade changes the `investor_profile` version.
2. Savings rate calculations and target dates adjust appropriately.
3. If risk appetite changes, the system creates a new draft IPS version.
4. Old versions of both profiles and IPS are preserved with `status = superseded`.

---

## 3. API Test Requirements

For all endpoints, integration tests should verify:
- **Success Case**: Valid payload returns expected JSON with `200 OK` or `201 Created`.
- **Bad Request**: Missing/invalid fields return `400 Bad Request`.
- **Authorization**: Accessing endpoints without a JWT Bearer token returns `401 Unauthorized`.
- **Resource Constraints**: Trying to double-log the same monthly review returns `409 Conflict`.

---

## 4. Manual Verification Scripts

Create QA scripts for testing key features manually:
- **Onboarding QA Checklist**: Validate that progressing through the 8 onboarding steps correctly maps variables to the resulting profile summary card.
- **Life Event QA Checklist**: Log a "Con đầu lòng" event, check that the dashboard displays the "Life Event Alert", check the "Life Timeline" tab lists the event, and check that the IPS tab has a draft version pending user review.

---

## 5. UI/UX Verification Guidelines

Verify the visual layout and user interactions to preserve premium design standards:
- **Glassmorphism Integrity Check**: Ensure card containers have `backdrop-filter: blur(12px)` and matching translucent borders (`rgba(255, 255, 255, 0.05)`). Card content must remain highly readable against background gradients.
- **Dark Mode Contrast**: Verify that text satisfies minimum contrast ratios in dark mode (Slate 50 on Slate 900 background). High-priority alerts must stand out using harmonized warning colors (Amber/Emerald/Red).
- **Responsive Layout Check**: View pages on mobile viewports (stacked cards layout) and tablet viewports (collapsed sidebar layout) to ensure no content overflow or text clipping.
- **Micro-animations & Interactive States**: Ensure all hoverable buttons, cards, and input fields trigger smooth hover transitions (200-300ms duration, ease-out behavior).
- **State Feedback**: Verify loading states (skeletons), empty states, and action results (error/success toast notifications) render correctly without visual jumpiness.

---

## 6. References & Design Documents
When validating functions, UI/UX elements, or writing tests, you MUST refer to:
- **Features List & Routing Grid**: [docs/features.md](file:///Users/datnguyen/Dev/Projects/Go/life_capital/docs/features.md)
- **Logic Specs & Auth Controls**: [docs/lld.md](file:///Users/datnguyen/Dev/Projects/Go/life_capital/docs/lld.md)
- **UI Mockups & CSS Standards**: [docs/frontend_layout.md](file:///Users/datnguyen/Dev/Projects/Go/life_capital/docs/frontend_layout.md)
