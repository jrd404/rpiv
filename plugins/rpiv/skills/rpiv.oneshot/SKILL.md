---
name: rpiv.oneshot
description: Quick research task without formal artifacts
disable-model-invocation: true
---

# Oneshot

You are tasked with performing a quick, focused research or investigation task. Unlike full research
workflows that produce formal documents and plans, this command is for quick questions that need a
thorough but concise answer delivered directly in the conversation.

## When to use this

- Quick architectural questions ("How does X work in this codebase?")
- Understanding a specific module, pattern, or integration
- Investigating a specific behavior or code path
- Answering "what would it take to do X?" without producing a formal plan
- Gathering context before deciding whether a full research/plan cycle is needed

## Process

### 1. Understand the question

- Read the user's question carefully
- Identify what specifically they need to know
- Determine the scope: is this about one file, one module, or cross-cutting?

### 2. Investigate

- Search the codebase for relevant code, configs, and documentation
- Read the key files thoroughly (do NOT skim or use limit/offset)
- Follow references and imports to understand the full picture
- Check for tests that demonstrate expected behavior
- Look at git history if understanding "why" something is the way it is matters

### 3. Synthesize and respond

- Answer the question directly and concisely
- Include relevant file paths so the user can dig deeper
- Use `file:line` references for specific code locations
- Call out any surprises, gotchas, or important caveats
- If the question turns out to be bigger than expected, say so and suggest a full research cycle

## Guidelines

- **No formal artifacts** - Don't create documents, plans, or research files. Answer in the
  conversation.
- **Be thorough but concise** - Read deeply, respond briefly. The user wants the answer, not a tour
  of everything you read.
- **Include file paths** - Always reference specific files so the user can verify and explore
  further.
- **Know when to stop** - If the investigation reveals the question needs a full research/plan
  cycle, tell the user rather than trying to answer an unbounded question in one shot.
- **Parallel investigation** - When multiple independent areas need to be checked, investigate them
  in parallel for speed.

