
# Implementation Plan: Style Books Table to Occupy One-Fourth Screen

**Branch**: `002-style-the-table` | **Date**: 2025-01-16 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-style-the-table/spec.md`

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path
   → If not found: ERROR "No feature spec at {path}"
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → Detect Project Type from file system structure or context (web=frontend+backend, mobile=app+api)
   → Set Structure Decision based on project type
3. Fill the Constitution Check section based on the content of the constitution document.
4. Evaluate Constitution Check section below
   → If violations exist: Document in Complexity Tracking
   → If no justification possible: ERROR "Simplify approach first"
   → Update Progress Tracking: Initial Constitution Check
5. Execute Phase 0 → research.md
   → If NEEDS CLARIFICATION remain: ERROR "Resolve unknowns"
6. Execute Phase 1 → contracts, data-model.md, quickstart.md, agent-specific template file (e.g., `CLAUDE.md` for Claude Code, `.github/copilot-instructions.md` for GitHub Copilot, `GEMINI.md` for Gemini CLI, `QWEN.md` for Qwen Code, or `AGENTS.md` for all other agents).
7. Re-evaluate Constitution Check section
   → If new violations: Refactor design, return to Phase 1
   → Update Progress Tracking: Post-Design Constitution Check
8. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
9. STOP - Ready for /tasks command
```

**IMPORTANT**: The /plan command STOPS at step 7. Phases 2-4 are executed by other commands:
- Phase 2: /tasks command creates tasks.md
- Phase 3-4: Implementation execution (manual or via tools)

## Summary
This feature updates the existing reading tracker webpage to display a condensed summary statistics table that occupies 25% of screen width on desktop (alongside existing charts at 75% width) and expands to full width on mobile devices. The table will show only top-level aggregated statistics (total books, pages, averages) without individual book rows, maintaining a single dashboard view for all reading data. This is a CSS-focused change to restructure the existing UI layout.

## Technical Context
**Language/Version**: JavaScript ES6+ (vanilla JS), CSS3  
**Primary Dependencies**: Chart.js (existing), no new dependencies  
**Storage**: N/A (existing books.json data source, no changes)  
**Testing**: Vitest (existing frontend test framework)  
**Target Platform**: Modern browsers (Chrome, Firefox, Safari, Edge) with responsive design
**Project Type**: web (frontend + backend monorepo)  
**Performance Goals**: Instant layout render (<16ms), smooth responsive transitions  
**Constraints**: Must maintain existing functionality, no breaking changes to API or data flow  
**Scale/Scope**: Single dashboard page with 2 sections (25% table + 75% charts on desktop)

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Initial Check (Before Phase 0) ✅

### I. Monorepo Structure ✅
- **Status**: PASS
- **Rationale**: Changes confined to existing `frontend/` workspace (CSS in `styles/main.css`, minor HTML/JS in `src/` if needed). No new workspaces or cross-workspace dependencies introduced.

### II. Minimal Dependencies ✅
- **Status**: PASS
- **Rationale**: Zero new dependencies. Uses existing CSS3 capabilities (Flexbox/Grid for layout, media queries for responsive design). No CSS frameworks or preprocessors required.

### III. Test-First Development (TDD) ✅
- **Status**: PASS
- **Rationale**: Will write visual regression tests or layout tests first using Vitest to verify 25% width on desktop, full width on mobile, and proper statistics display before implementing CSS changes.

### IV. API Contract Testing ✅
- **Status**: PASS (N/A)
- **Rationale**: No API changes. Existing `/api/stats?year=YYYY` endpoint provides necessary aggregated data. Frontend consumes existing contract.

### V. Simplicity & YAGNI ✅
- **Status**: PASS
- **Rationale**: Simplest solution uses CSS Grid or Flexbox for 2-column layout (25%/75% split) with media query for mobile breakpoint. No premature abstractions or frameworks.

---

### Post-Design Check (After Phase 1) ✅

**Re-evaluation Results**: All principles remain satisfied after design completion.

### I. Monorepo Structure ✅
- **Design Impact**: Confirmed - only `frontend/styles/main.css`, `frontend/index.html`, and `frontend/src/ui.js` modified. Backend and shared contracts untouched.

### II. Minimal Dependencies ✅
- **Design Impact**: Confirmed - zero new dependencies in design. CSS Grid (native), no utility libraries, no new npm packages.

### III. Test-First Development (TDD) ✅
- **Design Impact**: Test strategy defined (see quickstart.md and Phase 2 task ordering). Tests will be written before implementation (T001-T004 before T005-T011).

### IV. API Contract Testing ✅
- **Design Impact**: Confirmed - reuses existing `/api/stats` contract (documented in contracts/api-stats.md). No contract changes = no contract test updates needed.

### V. Simplicity & YAGNI ✅
- **Design Impact**: Confirmed - design uses CSS Grid with single breakpoint (768px). No complexity added. Research.md documents simplest viable approach.

**Constitution Compliance**: All principles satisfied at both checkpoints. No violations to document.

## Project Structure

### Documentation (this feature)
```
specs/[###-feature]/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
├── contracts/           # Phase 1 output (/plan command)
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)
```
backend/
├── internal/
│   ├── books/           # No changes
│   ├── handlers/        # No changes
│   └── stats/           # No changes (provides aggregated stats)
└── main.go              # No changes

frontend/
├── src/
│   ├── ui.js            # Modify: Update DOM structure for summary table
│   ├── api-client.js    # No changes (uses existing /api/stats)
│   ├── chart.js         # No changes
│   └── main.js          # Minor: Coordinate layout initialization
├── styles/
│   └── main.css         # Modify: Add layout styles (25%/75% grid, responsive)
├── index.html           # Modify: Update markup for summary table container
└── tests/
    └── ui.test.js       # Add: Layout and responsive behavior tests

shared/
└── contracts/           # No changes (existing API contracts sufficient)
```

**Structure Decision**: Web application (Option 2). This feature modifies only the frontend workspace within the existing monorepo. Changes are isolated to CSS (layout), HTML (markup structure), and UI rendering logic (JavaScript). Backend and contracts remain unchanged as the existing `/api/stats` endpoint already provides the necessary aggregated statistics data.

## Phase 0: Outline & Research
1. **Extract unknowns from Technical Context** above:
   - For each NEEDS CLARIFICATION → research task
   - For each dependency → best practices task
   - For each integration → patterns task

2. **Generate and dispatch research agents**:
   ```
   For each unknown in Technical Context:
     Task: "Research {unknown} for {feature context}"
   For each technology choice:
     Task: "Find best practices for {tech} in {domain}"
   ```

3. **Consolidate findings** in `research.md` using format:
   - Decision: [what was chosen]
   - Rationale: [why chosen]
   - Alternatives considered: [what else evaluated]

**Output**: research.md with all NEEDS CLARIFICATION resolved

## Phase 1: Design & Contracts
*Prerequisites: research.md complete*

1. **Extract entities from feature spec** → `data-model.md`:
   - Entity name, fields, relationships
   - Validation rules from requirements
   - State transitions if applicable

2. **Generate API contracts** from functional requirements:
   - For each user action → endpoint
   - Use standard REST/GraphQL patterns
   - Output OpenAPI/GraphQL schema to `/contracts/`

3. **Generate contract tests** from contracts:
   - One test file per endpoint
   - Assert request/response schemas
   - Tests must fail (no implementation yet)

4. **Extract test scenarios** from user stories:
   - Each story → integration test scenario
   - Quickstart test = story validation steps

5. **Update agent file incrementally** (O(1) operation):
   - Run `.specify/scripts/bash/update-agent-context.sh copilot`
     **IMPORTANT**: Execute it exactly as specified above. Do not add or remove any arguments.
   - If exists: Add only NEW tech from current plan
   - Preserve manual additions between markers
   - Update recent changes (keep last 3)
   - Keep under 150 lines for token efficiency
   - Output to repository root

**Output**: data-model.md, /contracts/*, failing tests, quickstart.md, agent-specific file

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
1. **CSS Layout Tasks**:
   - Create CSS Grid layout styles (`.dashboard-layout` with 25%/75% split)
   - Add responsive media query for mobile breakpoint (<768px)
   - Style summary card component (`.summary-card`, stat rows)

2. **HTML Structure Tasks**:
   - Update `index.html` to wrap existing content in grid layout
   - Add summary card container with semantic HTML
   - Preserve existing chart containers in charts section

3. **JavaScript UI Tasks**:
   - Update `ui.js` to populate summary card with stats data
   - Add number formatting function (comma separators)
   - Handle empty state rendering (0 books)

4. **Testing Tasks** (TDD - tests before implementation):
   - Write layout tests (verify 25%/75% split on desktop)
   - Write responsive tests (verify stacking on mobile)
   - Write content tests (verify stat values match API)
   - Write empty state tests

**Ordering Strategy**:
1. **Tests first** (TDD principle):
   - T001-T004: Write all UI tests (failing initially)
2. **HTML structure** (foundation):
   - T005: Update index.html markup
3. **CSS layout** (visual structure):
   - T006: Add grid layout styles
   - T007: Add responsive styles
   - T008: Style summary card component
4. **JavaScript logic** (data binding):
   - T009: Update ui.js to populate summary card
   - T010: Add number formatting
   - T011: Handle empty state
5. **Validation**:
   - T012: Run all tests (should pass)
   - T013: Execute quickstart.md manual validation

**Task Dependencies**:
- Tests (T001-T004) are independent [P]
- HTML changes (T005) must precede CSS styling
- CSS tasks (T006-T008) can be partially parallel [P]
- JS tasks (T009-T011) depend on HTML structure
- Final validation (T012-T013) depends on all implementation tasks

**Estimated Output**: 13-15 numbered, ordered tasks in tasks.md

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)  
**Phase 4**: Implementation (execute tasks.md following constitutional principles)  
**Phase 5**: Validation (run tests, execute quickstart.md, performance validation)

## Complexity Tracking
*Fill ONLY if Constitution Check has violations that must be justified*

**No violations detected.** All constitutional principles are satisfied by the design.

This feature adheres to:
- Monorepo structure (frontend-only changes)
- Minimal dependencies (zero new dependencies)
- TDD (tests before implementation)
- Contract stability (reuses existing API)
- Simplicity (CSS Grid with single breakpoint)


## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command)
- [x] Phase 1: Design complete (/plan command)
- [x] Phase 2: Task planning complete (/plan command - describe approach only)
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS
- [x] All NEEDS CLARIFICATION resolved
- [x] Complexity deviations documented (N/A - no deviations)

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*
