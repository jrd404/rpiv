---
name: dashboard.rpi
description: Display workflow state from docs/.tracker.json for all active tasks
disable-model-invocation: true
---

# Workflow Dashboard

Display the current state of all tracked RPI workflow tasks.

## Process

1. **Read tracker file**:
   - Read `docs/.tracker.json`
   - If the file does not exist, scan `docs/research/`, `docs/plans/`, and `docs/handoffs/` to
     bootstrap an initial tracker from discovered artifacts

2. **For each task, display**:
   - Title
   - Status (with visual indicator)
   - Current phase (if set)
   - Branch
   - Time since last update
   - Latest artifact path
   - Notes (if any)

3. **Sort and group**:
   - Sort by `updated` descending
   - Group: **Active** first (researching, planning, implementing, validating, review), then
     **Paused**, then **Done**

4. **Flag stale tasks**:
   - Any task not updated in 48+ hours gets a stale warning

5. **Suggest next action per task**:
   - `research-complete` -> "Run `/rpi:plan` to create an implementation plan"
   - `plan-complete` -> "Run `/rpi:implement <plan-path>` to begin implementation"
   - `paused` -> "Run `/rpi:resume <latest-handoff-path>` to continue"
   - `implementing` -> "Continue implementation or run `/rpi:handoff` to pause"
   - `validating` -> "Complete validation or run `/rpi:commit` when ready"
   - `review` -> "PR is open — merge when approved"

## Output Format

```
## Workflow Dashboard

### Active
| Task | Status | Phase | Branch | Last Updated | Next Action |
|------|--------|-------|--------|-------------|-------------|
| Auth middleware refactor | implementing | Phase 2 of 3 | feature/auth-refactor | 2h ago | Continue implementation |

### Paused
| Task | Status | Latest Handoff | Last Updated | Next Action |
|------|--------|---------------|-------------|-------------|
| API migration | paused | docs/handoffs/2026-03-14_16-30-00_api-migration.md | 3d ago (stale) | `/rpi:resume docs/handoffs/...` |

### Done
| Task | Completed | Branch |
|------|-----------|--------|
| Config cleanup | 2026-03-10 | feature/config-cleanup |
```

## Bootstrapping from Artifacts

If `docs/.tracker.json` is missing, scan for artifacts to build initial state:

1. List files in `docs/research/`, `docs/plans/`, `docs/handoffs/`
2. Extract description slugs from filenames
3. Match artifacts to tasks by slug similarity
4. Infer status from latest artifact type:
   - Only research doc -> `research-complete`
   - Has plan doc -> `plan-complete`
   - Has handoff doc -> `paused`
5. Create `docs/.tracker.json` with inferred state
6. Inform user: "Bootstrapped tracker from existing artifacts. Review and adjust if needed."

## Status Reference

| Status              | Meaning                    |
| ------------------- | -------------------------- |
| `researching`       | Research in progress       |
| `research-complete` | Research doc written       |
| `planning`          | Plan in progress           |
| `plan-complete`     | Plan doc written           |
| `implementing`      | Implementation in progress |
| `validating`        | Validation in progress     |
| `review`            | PR open                    |
| `done`              | Completed                  |
| `paused`            | Handoff created            |
