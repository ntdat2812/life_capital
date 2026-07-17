---
name: Business Analyst Skill
description: Core product requirements, user workflows, business logic, and specification standards for Life Capital.
---

# Business Analyst (BA) Skill Guidelines

## 1. Role Objective
Ensure that every feature, flow, and UI requirement matches the core vision: aiding a single user in making optimal personal wealth decisions on a monthly cadence based on their life situation.

---

## 2. Core Domain Concepts

### 2.1 Investor Profile & Life Timeline
The profile is dynamic. It represents the investor's current state (Dependents, Income Streams, Financial Goals, Risk Tolerance, Philosophy).
- **Life Events** are logged when a major event occurs.
- **Rule**: Every logged Life Event must trigger an AI analysis assessing the impact on the portfolio strategy and re-versions the Investor Profile.

### 2.2 Investment Policy Statement (IPS)
- The IPS defines rules (Target Allocations, Buy/Sell Rules, Risk Limits).
- **Rule**: Every IPS version is linked to a specific version of the Investor Profile. Changing the profile triggers a suggestion to update the IPS.

### 2.3 The Monthly Review Loop
Happens once a month:
1. **Life Event Verification**: User checks if any life events occurred.
2. **New Capital Input**: User inputs new investment amount, changes in monthly income/expenses.
3. **Portfolio Import/Update**: Update current holdings.
4. **AI Generation**: AI synthesizes and produces concrete recommendations (Buy ticker A with amount X, Buy ticker B with amount Y).

---

## 3. Specifications for Onboarding Flow
The onboarding must be a conversational chat interview consisting of 8 specific steps:
1. **Personal Information**: Age, career, status.
2. **Income & Cash Flow**: Primary salary + secondary side hustles.
3. **Current Assets**: Cash, gold, deposits, stocks.
4. **Dependents**: Children, parents.
5. **Financial Goals**: Retiring target, target net worth.
6. **Risk Appetite**: Behavior during crashes.
7. **Investment Style**: Value, Growth, dividend focus.
8. **Life Constraints**: Debt, mortgages, health costs.

*Result*: Generates a structured JSON profile and triggers the first IPS draft.

---

## 4. Key Stories to Implement & Verify
- **Log Life Event**: As a user, I want to add a life event (e.g. "Tăng lương 30%") so that my overall FI target date is recalculated and my target allocations are adjusted.
- **Update Thesis**: As a user, I want to keep an investment thesis updated per holding so that my monthly AI recommendations are aligned with my qualitative convictions.
- **Review Outcomes**: As a user, I want to check my logged decision journal entries after 6 months to score my past decision quality.

---

## 5. References & Design Documents
When working as a Business Analyst, you MUST check and update the following documents:
- **Core Vision & Workflows**: [docs/business.md](file:///Users/datnguyen/Dev/Projects/Go/life_capital/docs/business.md)
- **Features Spec & Screens**: [docs/features.md](file:///Users/datnguyen/Dev/Projects/Go/life_capital/docs/features.md)

