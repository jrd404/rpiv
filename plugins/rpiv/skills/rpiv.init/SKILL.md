---
name: rpiv.init
description: Initialize or repair the rpiv framework structure in the current repository — creates tracker, plan/handoff directories, migrates existing plans, and establishes workflow rules.
disable-model-invocation: true
---

# Initialize RPIV Framework

Set up or repair the rpiv workflow structure in the current repository.

## Process

### 1. Discover current state

Confirm this is a git repo. Check whether `docs/.tracker.json`, `docs/plans/`, and `docs/handoffs/`
exist.

### 2. Create missing structure

Create any missing directories: `docs/plans/`, `docs/handoffs/`.

### 3. Initialize or repair tracker

If `docs/.tracker.json` does not exist, create it: `{"version": 1, "tasks": {}}`.

If it exists, validate it has `version` and `tasks` keys. Fix if malformed.

**Task entry schema**: `title` (string), `status` (see below), `phase` (string|null), `branch`
(string|null), `created` (ISO 8601), `updated` (ISO 8601), `artifacts` (`{research, plan, handoffs[]}`), `notes` (string|null).

**Valid statuses**: `researching`, `research-complete`, `planning`, `plan-complete`, `implementing`,
`validating`, `review`, `done`, `paused`.

### 4. Migrate existing plans

Check if there is in-progress work with a plan that only lives in conversation context. If so:

1. Ask the user to confirm the current plan.
2. Write it to `docs/plans/YYYY-MM-DD_HH-MM-SS_description.md` using the standard plan template
   (phased structure with automated + manual success criteria per phase).
3. Add a task entry to the tracker.

Skip if nothing to migrate.

### 5. Present workflow rules

After setup:

```
RPIV framework initialized.

Workflow rules:
  1. Research before planning — plans reference file:line
     discoveries, not assumptions.
  2. Plan before implementing — follow an approved plan from
     docs/plans/.
  3. Phases are sequential gates — ALL success criteria
     (automated + manual) must pass and user must confirm
     before advancing.
  4. Validate before merging — run /rpiv.validate against
     the plan before opening a PR.
  5. Handoff before abandoning context — create a handoff
     so the next session can resume cleanly.

These are hard gates enforced by the rpiv skills.
```

### 6. Repair mode

If already initialized and valid, audit:

- Tasks with missing artifact files (tracker references a nonexistent doc)
- Artifact files not tracked (docs exist without tracker entries)
- Malformed tracker entries (missing required fields)

Fix what can be fixed automatically. Report what needs manual attention.

## Important

- Does NOT create `docs/research/` — that is created on-demand by `/rpiv.research`.
- If the repo already has a fully valid rpiv structure, confirm it and show the current status.
- The workflow rules are not suggestions. Skills must refuse to advance past unverified gates.
