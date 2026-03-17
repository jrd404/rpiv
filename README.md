# cc

A CLI for deploying Claude Code plugins across machines and projects.

`cc` packages the **research-plan-implement-validate** workflow as a self-contained Go binary. Run
`go install` on any machine and deploy the full framework to any project's `.claude/` directory — no
repo clone needed.

## Background

This project extracts and generalizes the "frequent intentional compaction" workflow originally
developed at [HumanLayer](https://github.com/humanlayer/humanlayer) for managing AI-assisted
development in complex codebases. The core idea: structure work into research, planning,
implementation, and validation phases, using parallel sub-agents to manage context efficiently and
placing human review at the highest-leverage points (research findings and plans, not just code).

The CLI follows the Cobra + hidden `__completer` pattern for dynamic bash tab completion.

## What's included

The embedded plugin provides 11 slash commands and 4 specialized sub-agents:

**Commands:** `/research`, `/plan`, `/implement`, `/validate`, `/iterate`, `/commit`,
`/describe-pr`, `/handoff`, `/resume`, `/debug`, `/oneshot`

**Agents:** `codebase-locator`, `codebase-analyzer`, `codebase-pattern-finder`, `web-researcher`

## Install

```bash
go install github.com/jarrodchung/cc/cmd/cc@latest
```

Or build from source:

```bash
make build
# binary is at bin/cc
```

f## Usage

### Deploy the plugin to a project

```bash
# Install all commands and agents to the current project
cc install all

# Install only commands (no agents)
cc install commands

# Install to your user-level .claude/ directory
cc install all --scope user

# Preview what would be installed
cc install all --dry-run
```

### Manage deployments

```bash
# Check status of installed files
cc status

# Update to latest embedded versions (skips locally modified files)
cc update

# Force-update everything, overwriting local modifications
cc update --force

# Remove all installed files
cc uninstall
```

### Explore what's available

```bash
# List available commands
cc list commands

# List available agents
cc list agents

# List available plugins
cc list plugins
```

### Tab completion

```bash
# Bash — add to .bashrc
complete -C 'cc __completer' cc

# Zsh — generate and install
cc completion zsh > ~/.zsh/completions/_cc
```

## How it works

Plugin assets (markdown files defining commands and agents) are embedded into the binary at build
time via `go:embed`. When you run `cc install`, it copies these files into the target `.claude/`
directory and writes a `.cc-manifest.json` tracking checksums and versions. Subsequent `cc update`
and `cc status` commands use the manifest to detect drift, skip locally modified files, and apply
updates cleanly.
