# rpiv

A CLI for deploying Claude Code plugins across machines and projects.

`rpiv` packages the **rpiv** (Research-Plan-Implement-Validate) workflow as a self-contained Go binary.
Run `go install` on any machine and deploy the full framework to any project's `.claude/` directory
— no repo clone needed.

## Background

This project extracts and generalizes the "frequent intentional compaction" workflow originally
developed at [HumanLayer](https://github.com/humanlayer/humanlayer) for managing AI-assisted
development in complex codebases. The core idea: structure work into research, planning,
implementation, and validation phases, using parallel sub-agents to manage context efficiently and
placing human review at the highest-leverage points (research findings and plans, not just code).

The CLI follows the Cobra + hidden `__completer` pattern for dynamic bash tab completion.

## What's included

The embedded plugin provides 13 skills and 4 specialized sub-agents, using the proper Claude Code
plugin format with `.claude-plugin/plugin.json`.

**Skills:**

- `/rpiv.init`
- `/rpiv.research`
- `/rpiv.plan`
- `/rpiv.implement`
- `/rpiv.validate`
- `/rpiv.iterate`
- `/rpiv.commit`
- `/rpiv.pr`
- `/rpiv.handoff`
- `/rpiv.resume`
- `/rpiv.debug`
- `/rpiv.oneshot`
- `/rpiv.status`

**Agents:**

- `codebase-locator`
- `codebase-analyzer`
- `codebase-pattern-finder`
- `web-researcher`

**Workflow tracking:** A `PostToolUse` hook auto-updates `docs/.tracker.json` when research docs,
plans, or handoffs are written. Use `/rpiv.status` to see all active work at a glance.

## Install

```bash
go install github.com/jarrodchung/rpiv/cmd/rpiv@latest
```

Or build from source:

```bash
make build
# binary is at bin/rpiv
```

## Usage

### Deploy the plugin to a project

```bash
# Install all skills and agents to the current project
rpiv install all

# Install only skills (no agents)
rpiv install skills

# Install to your user-level .claude/ directory
rpiv install all --scope user

# Preview what would be installed
rpiv install all --dry-run
```

### Manage deployments

```bash
# Check status of installed files
rpiv status

# Update to latest embedded versions (skips locally modified files)
rpiv update

# Force-update everything, overwriting local modifications
rpiv update --force

# Remove all installed files
rpiv uninstall
```

### Explore what's available

```bash
# List available skills
rpiv list skills

# List available agents
rpiv list agents

# List available plugins
rpiv list plugins
```

### Tab completion

```bash
# Bash — add to .bashrc
complete -C 'rpiv __completer' rpiv

# Zsh — generate and install
rpiv completion zsh > ~/.zsh/completions/_rpiv
```

## How it works

Plugin assets (skills, agents, hooks, and scripts) are embedded into the binary at build time via
`go:embed`. When you run `rpiv install`, it copies these files into the target `.claude/` directory
and writes a `.rpiv-manifest.json` tracking checksums and versions. Subsequent `rpiv update` and
`rpiv status` commands use the manifest to detect drift, skip locally modified files, and apply
updates cleanly.

### Plugin format

The rpiv plugin uses the Claude Code plugin format:

```
plugins/rpiv/
├── .claude-plugin/plugin.json    # Plugin manifest (name: "rpiv")
├── skills/*/SKILL.md             # 13 workflow skills
├── agents/*.md                   # 4 specialized sub-agents
├── hooks/hooks.json              # PostToolUse hook for tracker
└── scripts/update-tracker.sh     # Tracker auto-update script
```

When installed, skills are invocable as `/rpiv.research`, `/rpiv.plan`, etc.
