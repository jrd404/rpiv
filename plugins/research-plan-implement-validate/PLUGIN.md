# Research → Plan → Implement → Validate

A Claude Code plugin implementing the "frequent intentional compaction" workflow for solving
hard problems in complex codebases.

## Overview

This plugin provides slash commands and specialized sub-agents that structure AI-assisted
development around context management:

1. **Research** (`/research`) — Understand the codebase with parallel sub-agents
2. **Plan** (`/plan`) — Create detailed implementation plans interactively
3. **Implement** (`/implement`) — Execute plans phase-by-phase with verification
4. **Validate** (`/validate`) — Verify implementation matches plan intent

## Commands

| Command | Purpose |
|---|---|
| `/research` | Codebase research with parallel sub-agents |
| `/plan` | Interactive implementation planning |
| `/implement` | Execute approved plans phase-by-phase |
| `/validate` | Verify implementation correctness |
| `/iterate` | Update existing plans with feedback |
| `/commit` | Atomic git commits with user approval |
| `/describe-pr` | Generate PR descriptions from diffs |
| `/handoff` | Create session handoff documents |
| `/resume` | Resume work from handoff |
| `/debug` | Debug issues (logs, git, state) |
| `/oneshot` | Quick research+plan without formal artifacts |

## Agents

| Agent | Model | Purpose |
|---|---|---|
| codebase-locator | Sonnet | Find WHERE code lives |
| codebase-analyzer | Sonnet | Understand HOW code works |
| codebase-pattern-finder | Sonnet | Find similar patterns as templates |
| web-researcher | Sonnet | Strategic web search and synthesis |

## Artifact Output

Commands write research documents, plans, and handoffs to project-local directories:

```
docs/
├── research/    # From /research
├── plans/       # From /plan
└── handoffs/    # From /handoff
```

## Key Principles

- **Documentarians, not critics** — Agents document what EXISTS, never suggest improvements
- **Parallel sub-agent spawning** — Research spawns locator + analyzer + pattern-finder concurrently
- **Human review at highest leverage** — Review research and plans, not just code
- **No open questions in final deliverables** — Plans are complete and actionable
- **Phase-by-phase implementation** — Implement pauses after each phase for verification
