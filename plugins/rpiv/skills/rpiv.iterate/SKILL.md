---
name: rpiv.iterate
description: Iterate on existing implementation plans with thorough research and updates
model: opus
disable-model-invocation: true
---

# Iterate Implementation Plan

Update existing implementation plans based on user feedback. Be skeptical — verify technical
feasibility with code research before making changes.

**Tracker**: At the start of this workflow, update `docs/.tracker.json` — set the task status to
`planning`.

## Initial Response

- If NO plan file provided: ask for the path.
- If plan file provided but NO feedback: read the plan FULLY, then ask what to change.
- If BOTH provided: proceed directly to Step 1.

## Process Steps

### Step 1: Read and Understand

1. Read the existing plan FULLY (no limit/offset).
2. Parse the requested changes. Determine if they require codebase research.

### Step 2: Research If Needed

Only spawn research if the changes require new technical understanding.

- Use **codebase-locator**, **codebase-analyzer**, **codebase-pattern-finder** as needed.
- Be specific about directories in prompts.
- Read identified files FULLY. Wait for all sub-tasks to complete.

### Step 3: Confirm Approach

Present your understanding of the requested changes and what you found. Get user confirmation before
editing.

### Step 4: Update the Plan

1. Use the Edit tool for surgical changes — not wholesale rewrites.
2. Maintain existing structure unless explicitly changing it.
3. Keep file:line references accurate.
4. If adding phases, follow the existing pattern (automated + manual success criteria per phase).
5. If modifying scope, update "What We're NOT Doing."

### Step 5: Review

Present what changed and ask if further adjustments are needed.

## Guidelines

1. **Be surgical** — precise edits, preserve good content, don't over-engineer.
2. **No open questions** — if a change raises questions, research or ask immediately. Do not update
   the plan with unresolved questions.
3. **Verify before accepting** — don't blindly accept change requests that seem problematic. Point
   out conflicts with existing phases.
