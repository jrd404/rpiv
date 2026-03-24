---
name: rpiv.research
description: Research codebase comprehensively using parallel sub-agents
model: opus
disable-model-invocation: true
---

# Research Codebase

Conduct comprehensive codebase research by spawning parallel sub-agents and synthesizing findings.

**Tracker**: At the start of this workflow, update `docs/.tracker.json` — set the task status to
`researching`.

## Initial Setup

When invoked, ask for the research question. Then wait for the user's query.

## Steps

1. **Read mentioned files first**: If the user references specific files, read them FULLY (no
   limit/offset) yourself in the main context BEFORE spawning any sub-tasks.

2. **Decompose the question**: Break it into composable research areas. Create a todo list to track
   subtasks.

3. **Spawn parallel research agents**:
   - **codebase-locator** — find what exists and where
   - **codebase-analyzer** — understand how implementations work
   - **codebase-pattern-finder** — discover similar patterns and conventions
   - Each agent knows its job — tell it what to find, not how to search

4. **Wait for ALL agents to complete**, then synthesize:
   - Connect findings across components
   - Include specific file:line references
   - Highlight patterns, connections, and architectural decisions

5. **Gather metadata**: date/time, git commit, branch, repo name.

6. **Write research document** to `docs/research/YYYY-MM-DD_HH-MM-SS_description.md`:

   ```markdown
   ---
   date: [ISO 8601 with timezone]
   researcher: [name]
   git_commit: [hash]
   branch: [branch]
   repository: [repo]
   topic: "[Question/Topic]"
   tags: [research, <relevant-component-names>]
   status: complete
   last_updated: [YYYY-MM-DD]
   last_updated_by: [name]
   ---

   # Research: [Topic]

   ## Research Question
   [Original query]

   ## Summary
   [High-level findings]

   ## Detailed Findings
   ### [Component/Area]
   - Finding with reference ([file.ext:line](link))

   ## Code References
   - `path/to/file:123` - Description

   ## Open Questions
   [Areas needing further investigation]
   ```

7. **GitHub permalinks**: If on main branch or commit is pushed, replace local file references with
   `https://github.com/{owner}/{repo}/blob/{commit}/{file}#L{line}` permalinks.

8. **Present summary** to the user. Ask if they have follow-up questions.

9. **Follow-ups**: Append to the same document under `## Follow-up Research [timestamp]`. Update
   `last_updated` and `last_updated_by` in frontmatter.

## Important

- Always run fresh codebase research — never rely solely on existing documents.
- Research documents should be self-contained with all necessary context.
- Keep the main agent focused on synthesis; delegate deep file reading to sub-agents.
- Never write the document with placeholder values — gather all metadata first.
