<!--
Sync Impact Report
==================
Version change: Initial → 1.0.0
Principles established:
  - I. Monorepo Structure
  - II. Minimal Dependencies
  - III. Test-First Development (TDD)
  - IV. API Contract Testing
  - V. Simplicity & YAGNI
Sections added:
  - Core Principles (5 principles)
  - Technology Stack
  - Development Standards
  - Governance
Templates status:
  ✅ spec-template.md: Aligned with principles
  ✅ plan-template.md: Constitution check section applies all principles
  ✅ tasks-template.md: Task phases support TDD and monorepo structure
Follow-up TODOs: None
-->

# Reading App Constitution

## Core Principles

### I. Monorepo Structure

All application components MUST reside in a single repository with clear separation of concerns:
- **Backend** (`/backend`): Go services, APIs, and business logic
- **Frontend** (`/frontend`): JavaScript/TypeScript UI components
- **Shared** (`/shared`): Common contracts, types, and documentation
- **Tests** (`/tests`): Integration and end-to-end test suites

Each workspace MUST:
- Be independently buildable and testable
- Have its own dependency management (go.mod, package.json)
- Expose clear interfaces to other workspaces
- Include README documenting its purpose and commands

**Rationale**: Monorepo structure ensures atomic changes across frontend and backend, simplifies contract synchronization, and provides single source of truth for the entire application while maintaining clear boundaries.

### II. Minimal Dependencies

Prefer standard library and well-established, stable dependencies over feature-rich frameworks:

**Backend (Go)**:
- Use Go standard library wherever possible (net/http, encoding/json, database/sql)
- External dependencies MUST be justified by significant value vs. maintenance cost
- No frameworks that dictate application structure (prefer libraries over frameworks)
- Database access: Use standard database/sql with minimal driver-only dependencies

**Frontend (JavaScript)**:
- Prefer vanilla JS or minimal frameworks (e.g., Alpine.js, Lit) over heavy frameworks
- Build tools: Keep minimal (esbuild preferred for speed)
- No utility libraries for functionality easily implemented (e.g., no lodash for simple array operations)

**Dependency Approval Criteria**:
- Actively maintained (commit within last 6 months)
- Clear, focused purpose (does one thing well)
- Minimal transitive dependencies
- Strong security track record

**Rationale**: Minimal dependencies reduce security surface area, decrease build times, simplify debugging, and prevent version conflicts. The application remains maintainable long-term without framework churn.

### III. Test-First Development (NON-NEGOTIABLE)

TDD is mandatory for all feature development:

1. **Write tests first**: Before any implementation code
2. **User approval**: Tests reviewed and approved as spec validation
3. **Red**: Run tests → confirm they fail
4. **Green**: Implement minimal code to pass tests
5. **Refactor**: Improve code while keeping tests green

**Test Requirements**:
- Unit tests for business logic and data transformations
- Contract tests for API endpoints (request/response validation)
- Integration tests for critical user flows
- All tests MUST be runnable in CI/CD

**No Exceptions**: Features without tests cannot be merged.

**Rationale**: TDD ensures specifications are testable, catches regressions early, serves as living documentation, and produces more maintainable code through better design.

### IV. API Contract Testing

All APIs MUST have explicit contracts validated by automated tests:

- **Contract Definition**: OpenAPI/JSON Schema or code-based contracts in `/shared/contracts`
- **Backend Validation**: Server implementation validates against contract
- **Frontend Validation**: Client consumes typed contracts (generated or checked)
- **Test Coverage**: Every endpoint MUST have contract test verifying request/response structure
- **Breaking Changes**: Require explicit version bump and migration documentation

**Contract Test Requirements**:
- Verify HTTP status codes
- Validate request/response body schemas
- Check required vs. optional fields
- Test error response formats
- Document example payloads

**Rationale**: Explicit contracts prevent frontend/backend drift, enable parallel development, catch integration issues early, and provide clear API documentation.

### V. Simplicity & YAGNI

Start with the simplest solution that works; avoid premature optimization and speculative features:

**MUST**:
- Implement only specified requirements (no "nice to have" additions)
- Choose straightforward data structures and algorithms first
- Write clear, readable code over clever code
- Document complex decisions with rationale

**MUST NOT**:
- Add abstraction layers without proven need (wait for third use case)
- Implement features "for future flexibility"
- Optimize without profiling data showing bottleneck
- Use design patterns without clear benefit

**When Complexity is Justified**:
- Performance data shows measurable user impact
- Security requirements demand additional safeguards
- External system constraints force architectural decisions

Document justification in `research.md` for each feature.

**Rationale**: Simple solutions are easier to understand, debug, and modify. Complexity should be earned through real requirements, not anticipated ones.

## Technology Stack

### Backend
- **Language**: Go (latest stable version)
- **HTTP**: Standard library `net/http` or minimal router (e.g., chi, gorilla/mux)
- **Database**: SQLite (development) / PostgreSQL (production) via `database/sql`
- **Testing**: Standard `testing` package + table-driven tests
- **Build**: Standard `go build` with Makefile for common tasks

### Frontend
- **Language**: JavaScript (ES6+) or TypeScript for type safety
- **Framework**: Minimal (vanilla JS, Alpine.js, or Lit web components)
- **Build**: esbuild or Vite (fast, minimal config)
- **Testing**: Vitest or standard browser test runner
- **Styling**: CSS with minimal preprocessing or Tailwind if design system needed

### Shared
- **Contracts**: JSON Schema or shared TypeScript definitions
- **Documentation**: Markdown in repository
- **Scripts**: Shell scripts or Makefiles for automation

## Development Standards

### Code Organization
- **Backend**: Standard Go project layout (`cmd/`, `internal/`, `pkg/`)
- **Frontend**: Component-based structure (`/components`, `/pages`, `/lib`)
- **Shared**: Versioned contracts with changelog

### Quality Gates
All changes MUST pass:
- Tests (100% of new code covered)
- Linting (golangci-lint for Go, eslint for JS)
- Contract validation (no breaking changes without version bump)
- Build verification (all workspaces build successfully)

### Documentation Requirements
- README per workspace with setup and testing instructions
- API documentation via OpenAPI spec or inline comments
- Architecture decisions recorded in `/docs/decisions/` (ADRs)
- Code comments for non-obvious logic (not for obvious code)

### Performance Standards
- Backend API response time: <100ms p95 for standard operations
- Frontend initial load: <2s on 3G connection
- Database queries: <50ms for common queries
- Build time: <30s for full monorepo build

## Governance

### Constitutional Authority
This constitution supersedes all other development practices, guidelines, and preferences. When in doubt, constitution principles apply.

### Amendment Process
1. Propose amendment with clear rationale
2. Document impact on existing features and templates
3. Update affected templates and documentation
4. Increment version using semantic versioning:
   - **MAJOR**: Breaking governance changes or principle removals
   - **MINOR**: New principles or significant expansions
   - **PATCH**: Clarifications, typos, non-semantic refinements
5. Commit with message format: `docs: amend constitution to v{VERSION} ({summary})`

### Compliance Verification
- All feature specifications MUST align with principles
- Implementation plans MUST include "Constitution Check" section
- Code reviews MUST verify principle adherence
- CI/CD pipeline MUST enforce quality gates

### Complexity Exceptions
If a principle violation is necessary:
1. Document justification in feature `research.md`
2. Identify specific principle conflict
3. Explain why standard approach is insufficient
4. Propose mitigation strategy
5. Obtain explicit approval before implementation

Complexity without justification will be rejected.

### Development Guidance
Runtime development guidance and agent-specific instructions are maintained separately:
- `.github/copilot-instructions.md` for GitHub Copilot
- `.specify/templates/agent-file-template.md` template for other agents

**Version**: 1.0.0 | **Ratified**: 2025-01-16 | **Last Amended**: 2025-01-16