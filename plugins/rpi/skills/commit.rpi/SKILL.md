---
name: commit.rpi
description: Create git commits with user approval and no Claude attribution
disable-model-invocation: true
---

# Commit Changes

You are tasked with creating git commits for the changes made during this session.

## Commit Message Format

Use **Conventional Commits** with imperative mood. Focus on _why_, not _what_.

```
<type>[optional scope]: <description>

[optional body]
```

Types: `feat`, `fix`, `refactor`, `docs`, `test`, `chore`, `ci`, `perf`, `build`, `style`

## Process

1. **Assess changes:**
   - Run `git status` and `git diff` to understand modifications
   - Decide whether changes should be one commit or multiple logical commits

2. **Draft commit(s):**
   - Group related files together
   - Write Conventional Commit messages
   - Keep commits focused and atomic

3. **Get user approval:**
   - List files and commit message(s)
   - **Wait for explicit user confirmation before committing**

4. **Execute upon confirmation:**
   - Use `git add` with specific files (never use `-A` or `.`)
   - Create commits with approved messages
   - Show the result with `git log --oneline -n [number]`

## Rules

- **NEVER add co-author lines or Claude attribution** — no `Co-Authored-By`, no "Generated with
  Claude" footers. Commits are authored solely by the user.
- **NEVER commit** `.claude/` directories or config files matching
  `/^(\.agents|CLAUDE|AGENTS(\.override)?|GEMINI|TEAM_GUIDE)\.md$/`. These are local config files
  and must not be tracked in version control.
- **Always get user approval** on commit messages before executing.
