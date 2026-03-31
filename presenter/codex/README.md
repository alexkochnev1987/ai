# Codex Multi-Agent Development Workflow

This folder defines role separation for implementing the AI Web Presenter.

## Roles

1. **Orchestrator** (`orchestrator.md`)
   - Breaks work into tasks and phases.
   - Assigns implementation and validation to other agents.
   - Tracks dependencies and release readiness.

2. **Developer** (`developer.md`)
   - Implements product features.
   - Writes/updates code and migrations.
   - Documents technical decisions.

3. **Tester** (`tester.md`)
   - Validates core flows end-to-end.
   - Runs functional, regression, and error-path checks.
   - Reports defects with reproducible steps.

4. **Reviewer** (`reviewer.md`)
   - Performs code review and architecture checks.
   - Verifies style, security, and maintainability.
   - Produces merge recommendation report.

## Operating model

- Orchestrator is responsible for task routing and acceptance criteria.
- Developer cannot self-approve delivery without Tester and Reviewer sign-off.
- Tester focuses on behavior correctness; Reviewer focuses on code quality and risk.
- All reports should link back to implementation phase(s) from `../IMPLEMENTATION_PLAN.md`.
