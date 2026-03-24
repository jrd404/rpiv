---
name: rpiv.plan
description: Create detailed implementation plans with thorough research and iteration
model: opus
disable-model-invocation: true
---

# Implementation Plan

Create detailed implementation plans through an interactive, iterative process. Be skeptical,
thorough, and collaborative.

**Tracker**: At the start of this workflow, update `docs/.tracker.json` — set the task status to
`planning`.

## Initial Response

When invoked:

1. If a file path or description was provided, read those files FULLY and begin research.
2. If no parameters, ask the user for the task description, constraints, and links to related
   research. Then wait for input.

## Process Steps

### Step 1: Context Gathering & Initial Analysis

1. **Read all mentioned files FULLY** (no limit/offset). Do this yourself in the main context BEFORE
   spawning any sub-tasks.

2. **Spawn parallel research agents**:
   - **codebase-locator** — find all files related to the task
   - **codebase-analyzer** — understand current implementation
   - **codebase-pattern-finder** — find similar features to model after

3. **Read all files identified by research** FULLY into main context.

4. **Cross-reference** requirements with actual code. Note discrepancies and assumptions.

5. **Present understanding and focused questions** — show what you found with file:line references.
   Only ask questions you genuinely cannot answer through code investigation.

### Step 2: Research & Discovery

After initial clarifications:

1. If the user corrects a misunderstanding, spawn new research to verify — do not blindly accept.
2. Spawn parallel sub-tasks for deeper investigation. Be specific about directories in prompts.
3. Wait for ALL sub-tasks to complete.
4. Present findings and design options with pros/cons. Ask which approach to pursue.

### Step 3: Plan Structure Development

1. Propose the phase structure as an outline. Get feedback before writing details.

### Step 4: Detailed Plan Writing

Write the plan to `docs/plans/YYYY-MM-DD_HH-MM-SS_description.md` (kebab-case description).

**Template:**

````markdown
# [Feature/Task Name] Implementation Plan

## Overview

[What we're implementing and why]

## Current State Analysis

[What exists, what's missing, key constraints]

### Key Discoveries:

- [Finding with file:line reference]
- [Pattern to follow]

## What We're NOT Doing

[Out-of-scope items]

## Implementation Approach

[High-level strategy]

## Phase 1: [Descriptive Name]

### Overview

[What this phase accomplishes]

### Changes Required:

#### 1. [Component/File Group]

**File**: `path/to/file.ext` **Changes**: [Summary]

### Success Criteria:

#### Automated Verification:

- [ ] Tests pass: `make test`
- [ ] Type checking passes: `npm run typecheck`
- [ ] Linting passes: `make lint`

#### Manual Verification:

- [ ] Feature works as expected in UI
- [ ] Edge cases verified manually
- [ ] No regressions in related features

**Gate**: All automated AND manual criteria must pass before advancing to the next phase. Pause for
human confirmation after automated checks pass.

---

## Phase 2: [Descriptive Name]

[Same structure...]

---

## Testing Strategy

[Unit, integration, and manual testing approach]

## References

- Related research: `docs/research/[relevant].md`
- Similar implementation: `[file:line]`
````

### Step 5: Review

Present the draft location. Iterate based on feedback until the user is satisfied.

## Guidelines

1. **Be skeptical** — question vague requirements, verify with code, don't assume.
2. **Be interactive** — get buy-in at each step. Don't write the full plan in one shot.
3. **No open questions in the final plan** — if you hit an open question, STOP. Research or ask
   immediately. The plan must be complete and actionable before finalizing.
4. **Every phase needs both automated and manual success criteria.** The template above shows the
   required format. Automated = commands agents can run. Manual = requires human testing.
