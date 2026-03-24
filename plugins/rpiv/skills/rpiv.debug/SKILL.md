---
name: rpiv.debug
description: Debug issues by investigating logs, process state, and git history
disable-model-invocation: true
---

# Debug

Investigate problems by examining logs, process state, and git history. This is a read-only
investigation — do NOT edit files.

## Process

### Step 1: Understand the Problem

If invoked with a context file (plan, ticket), read it FULLY. Ask the user what specific problem
they are encountering, what they were trying to do, and any error messages.

### Step 2: Discover the Environment

- Identify project type from config files (`package.json`, `go.mod`, `Makefile`, `docker-compose.yml`, etc.)
- Find log locations (project logs, Docker logs, system logs)
- Identify running services and listening ports

### Step 3: Investigate

Spawn parallel investigation tasks:

1. **Logs** — find and analyze recent errors, warnings, stack traces around the problem timeframe.
2. **Process/Service State** — verify expected processes are running, check for crashes, port
   conflicts, resource issues.
3. **Git/File State** — check git status, recent commits, uncommitted changes, file permission
   issues.

### Step 4: Present Findings

```
## Debug Report

### What's Wrong
[Clear statement based on evidence]

### Evidence Found
[From logs, services, git — with timestamps and sources]

### Root Cause
[Most likely explanation]

### Next Steps
1. [Specific action to try first]
2. [Alternative if that doesn't work]
```

Acknowledge anything outside your reach (browser console, external services, system-level issues).

## Important

- **Investigation only** — do not edit code.
- **Discover before assuming** — find log paths and service names, don't guess.
- Read context files FULLY (no limit/offset).
