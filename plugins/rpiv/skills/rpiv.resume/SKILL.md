---
name: rpiv.resume
description: Resume work from handoff document with context analysis and validation
disable-model-invocation: true
---

# Resume Work from Handoff

Resume work from a handoff document. Validate that the handoff state matches reality before
proceeding.

**Tracker**: At the start of this workflow, update `docs/.tracker.json` — set the task status to
`implementing`.

## Initial Response

- **Path provided**: Read the handoff FULLY. Read any research/plan docs it links to under `docs/`
  (do NOT use a sub-agent for these critical files).
- **Ticket number provided**: Find the most recent handoff in `docs/handoffs/TICKET-NUM/` by
  filename timestamp. If not found, ask the user for the path.
- **No parameters**: List contents of `docs/handoffs/` and ask which to resume.

## Process Steps

### Step 1: Read and Analyze Handoff

1. Read the handoff FULLY (no limit/offset). Extract tasks, statuses, learnings, artifacts, and
   action items.
2. Spawn focused research tasks to verify current state against handoff claims.
3. Wait for ALL sub-tasks to complete.
4. Read critical files from "Learnings" and "Recent changes" sections.

### Step 2: Present Analysis

Present:
- Original tasks and their verified current status
- Key learnings validated (still valid / changed)
- Recent changes verified (present / missing / modified)
- Recommended next actions based on handoff + current state
- Any conflicts or regressions found

Get confirmation before proceeding.

### Step 3: Execute

1. Create a todo list from the action items.
2. Begin with the first approved task.
3. Apply learnings and patterns from the handoff throughout.
4. Never assume handoff state matches current state — verify file references and patterns before
   acting on them.
