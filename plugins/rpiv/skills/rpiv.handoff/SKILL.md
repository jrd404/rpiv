---
name: rpiv.handoff
description: Create handoff document for transferring work to another session
disable-model-invocation: true
---

# Create Handoff

Write a handoff document to transfer your work context to another agent session. Be thorough but
concise — compact and summarize without losing key details.

**Tracker**: At the start of this workflow, update `docs/.tracker.json` — set the task status to
`paused`.

## Process

### 1. Filepath & Metadata

- **Path**: `docs/handoffs/YYYY-MM-DD_HH-MM-SS_description.md` (kebab-case description, 24h time)
- **With ticket**: `docs/handoffs/TICKET-NUM/YYYY-MM-DD_HH-MM-SS_TICKET-NUM_description.md`
- **Gather**: git branch, commit hash, repo name
- Example: `docs/handoffs/PROJ-2166/2025-01-08_13-55-22_PROJ-2166_create-context-compaction.md`

### 2. Write the handoff

Use this template:

    ---
    date: [ISO 8601 with timezone]
    git_commit: [hash]
    branch: [branch]
    repository: [repo]
    topic: "[Feature/Task Name]"
    tags: [handoff, <relevant-component-names>]
    status: complete
    last_updated: [YYYY-MM-DD]
    type: handoff
    ---

    # Handoff: {concise description}

    ## Task(s)
    {Tasks and their statuses (completed/in-progress/planned). Reference plan docs
    and current phase if applicable.}

    ## Critical References
    {2-3 most important spec/design docs. Leave blank if none.}

    ## Recent Changes
    {Changes you made, in file:line syntax.}

    ## Learnings
    {Patterns, root causes, gotchas the next agent should know. Include file paths.}

    ## Artifacts
    {Exhaustive list of produced/updated artifacts as file paths.}

    ## Action Items & Next Steps
    {What the next agent should do, based on task statuses.}

    ## Other Notes
    {Anything that doesn't fit above but is useful context.}

### 3. Save

Write to the filepath from step 1. Respond with:

    Handoff created! Resume with: /rpiv.resume path/to/handoff.md

## Important

- **More information, not less** — the template is a minimum.
- **Prefer file:line references over code snippets** — avoid large blocks or diffs unless debugging
  an error.
