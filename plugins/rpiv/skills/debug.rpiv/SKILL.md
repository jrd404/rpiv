---
name: debug.rpiv
description: Debug issues by investigating logs, process state, and git history
disable-model-invocation: true
---

# Debug

You are tasked with helping debug issues during manual testing or implementation. This command
allows you to investigate problems by examining logs, process state, and git history without editing
files. Think of this as a way to bootstrap a debugging session without using the primary window's
context.

## Initial Response

When invoked WITH a context file (plan, ticket, or issue description):

    I'll help debug issues with [file name]. Let me understand the current state.

    What specific problem are you encountering?
    - What were you trying to test/implement?
    - What went wrong?
    - Any error messages?

    I'll investigate logs, process state, and git history to help figure out what's happening.

When invoked WITHOUT parameters:

    I'll help debug your current issue.

    Please describe what's going wrong:
    - What are you working on?
    - What specific problem occurred?
    - When did it last work?

    I can investigate logs, process state, and recent changes to help identify the issue.

## Process Steps

### Step 1: Understand the Problem

After the user describes the issue:

1. **Read any provided context** (plan or ticket file):
   - Understand what they're implementing/testing
   - Note which phase or step they're on
   - Identify expected vs actual behavior

2. **Quick state check**:
   - Current git branch and recent commits
   - Any uncommitted changes
   - When the issue started occurring

### Step 2: Discover the Environment

Before investigating, discover what's available in this project:

1. **Identify the project type and tooling**:
   - Check for `package.json`, `go.mod`, `Cargo.toml`, `pyproject.toml`, `Makefile`,
     `docker-compose.yml`, etc.
   - Determine what services/processes should be running

2. **Find log locations**:
   - Check common log directories: `logs/`, `log/`, `tmp/`, `.log` files in project root
   - Check if the project defines log paths in config files (`.env`, config directories)
   - Check system logs if relevant: `journalctl`, `/var/log/`, `~/Library/Logs/`
   - Check for Docker/container logs: `docker compose logs`
   - Ask the user if you can't find logs automatically

3. **Identify running services**:
   - Check for running processes related to the project: `ps aux | grep [project-relevant-terms]`
   - Check for listening ports if relevant: `lsof -i -P | grep LISTEN`
   - Check Docker containers if applicable: `docker compose ps`

### Step 3: Investigate the Issue

Spawn parallel investigation tasks based on what you discovered:

**Task 1 - Check Recent Logs:** Find and analyze the most recent logs for errors:

1. Look in discovered log locations
2. Search for errors, warnings, or stack traces around the problem timeframe
3. Look for patterns or repeated errors
4. Check stderr output from running services Return: Key errors/warnings with timestamps

**Task 2 - Process and Service State:** Check the current state of relevant services:

1. Verify expected processes are running
2. Check for crashed or zombie processes
3. Look for port conflicts
4. Check resource usage (memory, disk) if relevant Return: Service health findings

**Task 3 - Git and File State:** Understand what changed recently:

1. Check git status and current branch
2. Look at recent commits: `git log --oneline -10`
3. Check uncommitted changes: `git diff`
4. Verify expected files exist
5. Look for any file permission issues Return: Git state and any file issues

### Step 4: Present Findings

Based on the investigation, present a focused debug report:

    ## Debug Report

    ### What's Wrong
    [Clear statement of the issue based on evidence]

    ### Evidence Found

    **From Logs:**
    - [Error/warning with timestamp and source file]
    - [Pattern or repeated issue]

    **From Services/Processes:**
    - [Running/stopped/crashed status]
    - [Port or resource issues]

    **From Git/Files:**
    - [Recent changes that might be related]
    - [File state issues]

    ### Root Cause
    [Most likely explanation based on evidence]

    ### Next Steps

    1. **Try This First**:
       [Specific command or action]

    2. **If That Doesn't Work**:
       - [Alternative approach]
       - [Another diagnostic step]

    ### Can't Access?
    Some issues might be outside my reach:
    - Browser console errors (F12 in browser)
    - External service state (APIs, databases on remote hosts)
    - System-level issues (kernel, hardware)

    Would you like me to investigate something specific further?

## Important Notes

- **Focus on investigation, not fixing** - This command is for understanding problems, not editing
  code
- **Always require problem description** - Can't debug without knowing what's wrong
- **Read files completely** - No limit/offset when reading context
- **Discover before assuming** - Don't assume log paths or service names; find them
- **Guide back to user** - Some issues (browser console, external services) are outside reach
- **No file editing** - Pure investigation only

## Quick Reference

**Find Logs**: # Project-specific logs find . -name "\*.log" -mmin -60 2>/dev/null | head -10 ls -lt
logs/ 2>/dev/null | head -5

    # Docker logs
    docker compose logs --tail=50 2>/dev/null

    # System logs (macOS)
    log show --predicate 'process == "your-process"' --last 30m 2>/dev/null

**Service Check**: # Find project-related processes ps aux | grep -i [project-name]

    # Check listening ports
    lsof -i -P | grep LISTEN

    # Docker services
    docker compose ps 2>/dev/null

**Git State**: git status git log --oneline -10 git diff git diff --stat HEAD~5

Remember: This command helps you investigate without burning the primary window's context. Perfect
for when you hit an issue during manual testing and need to dig into logs, processes, or git state.
