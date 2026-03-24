---
name: rpiv.validate
description: Validate implementation against plan, verify success criteria, identify issues
disable-model-invocation: true
---

# Validate Plan

Verify that an implementation plan was correctly executed. Check all success criteria and identify
deviations.

**Tracker**: At the start of this workflow, update `docs/.tracker.json` — set the task status to
`validating`.

## Process

### Step 1: Gather Context

1. If in an existing conversation, review what was implemented. If starting fresh, discover through
   git and codebase analysis.
2. Locate the plan (use provided path, search recent commits, or ask user).
3. Read the plan FULLY (no limit/offset).
4. Run `make check test` or equivalent comprehensive checks.
5. Check recent commits: `git log --oneline -n 20` and `git diff` covering implementation commits.

### Step 2: Systematic Validation

For each phase in the plan:

1. **Verify completion** — check that plan checkmarks (- [x]) reflect actual code state.
2. **Run automated verification** — execute each command from "Automated Verification." Document
   pass/fail.
3. **Assess manual criteria** — list what needs human testing with clear steps.
4. **Think about edge cases** — error handling, missing validations, potential regressions.

Spawn parallel research tasks if starting fresh and need to discover what was implemented.

### Step 3: Generate Report

```markdown
## Validation Report: [Plan Name]

### Implementation Status
- Phase N: [Name] - [Fully/Partially implemented]

### Automated Verification Results
- [Command]: [pass/fail]

### Code Review Findings
- Matches plan: [what aligns]
- Deviations: [what differs and why]
- Potential issues: [concerns]

### Manual Testing Required
- [ ] [Specific test with steps]

### Recommendations
- [Actionable items]
```

## Important

- Run all automated checks — don't skip verification commands.
- Be honest about incomplete items or shortcuts taken.
- This works best after commits are made, as it can analyze git history.
