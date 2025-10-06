# Tasks: Style Books Table to Occupy One-Fourth Screen

**Feature**: 002-style-the-table  
**Input**: Design documents from `/specs/002-style-the-table/`  
**Prerequisites**: plan.md, research.md, data-model.md, contracts/, quickstart.md

## Execution Flow Summary
1. Loaded plan.md: CSS Grid layout for 25%/75% split (desktop), stacked (mobile)
2. Tech Stack: JavaScript ES6+ (vanilla JS), CSS3, Chart.js (existing), Vitest
3. Structure: Monorepo `frontend/` workspace (HTML, CSS, JS changes only)
4. Contracts: Reuses existing `/api/stats` endpoint (no backend changes)
5. Data Model: ReadingStatistics entity (existing), SummaryCard + DashboardLayout components
6. Task Categories: Tests first (TDD), HTML structure, CSS layout, JS logic, validation
7. Dependencies: Tests → HTML → CSS → JS → Validation

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- All paths relative to repository root (`/Users/kristenwomack/code/reading-app`)

## Phase 3.1: Setup
- [x] **T001** Verify Vitest is installed and configured in `frontend/package.json` (no changes expected)

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

- [x] **T002 [P]** Create desktop layout test in `frontend/tests/ui.test.js`:
  - Test name: "dashboard layout shows 25%/75% split on desktop (>768px)"
  - Mock window.matchMedia to simulate 1440px viewport
  - Assert `.dashboard-layout` has `grid-template-columns: 1fr 3fr`
  - Assert `.summary-section` width is approximately 25% of container
  - Assert `.charts-section` width is approximately 75% of container

- [x] **T003 [P]** Create mobile layout test in `frontend/tests/ui.test.js`:
  - Test name: "dashboard layout stacks vertically on mobile (<768px)"
  - Mock window.matchMedia to simulate 375px viewport
  - Assert `.dashboard-layout` has `grid-template-columns: 1fr` OR `flex-direction: column`
  - Assert `.summary-section` width is 100% of container
  - Assert no horizontal scrolling

- [x] **T004 [P]** Create summary card content test in `frontend/tests/ui.test.js`:
  - Test name: "summary card displays correct statistics"
  - Mock API response: `{ year: 2025, total_books: 42, total_pages: 12450, avg_pages_per_book: 296 }`
  - Assert card title shows "2025 Summary"
  - Assert "Total Books" row shows "42"
  - Assert "Total Pages" row shows "12,450" (with comma separator)
  - Assert "Avg Pages/Book" row shows "296"

- [x] **T005 [P]** Create empty state test in `frontend/tests/ui.test.js`:
  - Test name: "summary card shows empty state when no books"
  - Mock API response: `{ year: 2024, total_books: 0, total_pages: 0, avg_pages_per_book: 0 }`
  - Assert card shows "No books tracked for this year" OR similar message
  - Assert no error rendering

## Phase 3.3: Core Implementation (ONLY after tests are failing)

### HTML Structure
- [x] **T006** Update `frontend/index.html`:
  - Wrap existing content in `<div class="dashboard-layout">`
  - Add `<aside class="summary-section">` with `<div class="summary-card">` inside
  - Move existing charts into `<main class="charts-section">`
  - Preserve all existing chart elements (canvas elements, etc.)

### CSS Layout
- [x] **T007** Create desktop grid layout in `frontend/styles/main.css`:
  - Add `.dashboard-layout` styles with `display: grid` and `grid-template-columns: 1fr 3fr`
  - Add `gap: 1.5rem` for spacing between sections
  - Add `align-items: start` to align tops of sections
  - Wrap in `@media (min-width: 769px)` media query

- [x] **T008** Create mobile stacked layout in `frontend/styles/main.css`:
  - Add `.dashboard-layout` styles with `display: flex` and `flex-direction: column`
  - Add `gap: 1rem` for vertical spacing
  - Set `.summary-section` width to 100%
  - Wrap in `@media (max-width: 768px)` media query

- [x] **T009** Style summary card component in `frontend/styles/main.css`:
  - Add `.summary-card` styles (border, padding, background color)
  - Add `.summary-title` styles for heading (e.g., "2025 Summary")
  - Add `.stats-list` styles for vertical stack container
  - Add `.stat-row` styles with flexbox for label/value alignment
  - Add `.stat-label` styles (left-aligned text)
  - Add `.stat-value` styles (right-aligned, bold/emphasized)
  - Add `.empty-state` and `.empty-message` styles for 0 books case

### JavaScript Logic
- [x] **T010** Create summary card rendering function in `frontend/src/ui.js`:
  - Add `renderSummaryCard(stats)` function
  - Takes ReadingStatistics object as parameter
  - Selects `.summary-card` element from DOM
  - Updates `.summary-title` with `${stats.year} Summary`
  - Renders stat rows for total_books, total_pages, avg_pages_per_book
  - Handles empty state when total_books === 0

- [x] **T011** Add number formatting function in `frontend/src/ui.js`:
  - Create `formatNumber(num)` function
  - Use `num.toLocaleString('en-US')` for comma separators
  - Apply to total_pages display value

- [x] **T012** Integrate summary card with year selection in `frontend/src/main.js`:
  - Call `renderSummaryCard(stats)` when stats API response received
  - Ensure summary updates when year selector changes
  - Coordinate with existing chart rendering (no breaking changes)

### Chart Visualization Update
- [x] **T013** Update chart from bar to line chart in `frontend/src/chart.js`:
  - Change `type: 'bar'` to `type: 'line'` in Chart.js config
  - Update dataset configuration:
    - Set `borderColor` to 'rgb(75, 192, 192)' (teal line color)
    - Set `backgroundColor` to 'rgba(75, 192, 192, 0.2)' (light fill under line)
    - Add `tension: 0.1` for slight curve in line (smooth appearance)
    - Add `fill: true` to show area under line
  - Keep existing responsive options and Y-axis configuration
  - Preserve chart destroy/recreate logic for year switching

## Phase 3.4: Polish & Validation
- [ ] **T014** Run all Vitest tests: `cd frontend && npm test`
  - Verify T002-T005 tests now PASS
  - Fix any failing tests before proceeding

- [ ] **T015** Execute manual validation from `specs/002-style-the-table/quickstart.md`:
  - Start backend: `cd backend && go run main.go`
  - Open `http://localhost:3000` in browser
  - Complete all validation steps (desktop layout, mobile layout, content accuracy, responsive transition)
  - Test on Chrome, Firefox, Safari, Edge
  - Verify no horizontal scrolling at any width
  - Verify line chart renders correctly in 75% section with connected data points

- [ ] **T016 [P]** Update `README.md` at repository root:
  - Add screenshot or description of new dashboard layout (optional)
  - Update "Features" section to mention summary card (if applicable)

- [ ] **T017** Final cleanup:
  - Remove any console.log statements added during development
  - Verify no unused CSS classes
  - Run linter if configured: `cd frontend && npm run lint` (if exists)

## Dependencies
```
T001 (setup verification)
  ↓
T002, T003, T004, T005 (tests - parallel)
  ↓
T006 (HTML structure)
  ↓
T007, T008, T009 (CSS layout - can overlap)
  ↓
T010, T011 (JS rendering logic)
  ↓
T012 (JS integration)
  ↓
T013 (chart update to line)
  ↓
T014 (test validation)
  ↓
T015 (manual validation)
  ↓
T016, T017 (polish - parallel)
```

## Parallel Execution Examples
```bash
# After T001, run tests in parallel:
# T002, T003, T004, T005 can all be written simultaneously

# After T006, run CSS tasks in parallel:
# T007 (desktop grid) and T008 (mobile stack) are independent
# T009 (summary card styles) can start after T007 OR T008 complete

# After T015, run polish tasks in parallel:
# T016 (README update) and T017 (cleanup) are independent
```

## Task Validation Checklist
- [x] All contracts accounted for (existing `/api/stats` reused, no new tests needed)
- [x] All UI components have tests (SummaryCard, DashboardLayout)
- [x] Tests come before implementation (T002-T005 before T006-T012)
- [x] Parallel tasks are truly independent (different files)
- [x] Each task specifies exact file path
- [x] No [P] tasks modify same file

## Notes
- **Backend**: No changes required - existing `/api/stats` endpoint provides all necessary data
- **Constitution compliance**: Zero new dependencies, TDD enforced, monorepo structure preserved
- **Performance goal**: Layout render < 16ms (validate in T014 with DevTools Performance tab)
- **Browser support**: CSS Grid is supported in all modern browsers (Chrome, Firefox, Safari, Edge)
- **Breakpoint**: Single breakpoint at 768px (mobile/tablet boundary)
