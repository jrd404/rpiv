# Advanced Context Engineering for Coding Agents

1. what we learned applying context engineering to coding agents
2. the dimensions along which using these agents is a deeply technical craft
3. why I don't believe these approaches are generalizable
4. the number of times I've been repeatedly proven wrong about (3)

## Intentional Compaction

Intentional compaction (like `/compact` in Claude Code) is pausing work and starting over with fresh
context window. Use a prompt like:

> "Write everything we did so far to progress.md, ensure to note the approach we're taking, the
> steps we've done so far, and the current failure we're working on."

Git commit messages can also be used for intentional compaction:

> look at the head commit to understand what we're doing, read every file touched (no subagent)
> then - we're going to continue to iterate on {projectFeature}

This context window will be compacted into a `progress.md` file.

```
Session 1                                          Session 2
+----------------------+                           +----------------------+
| System Instructions  |                           | System Instructions  |
| CLAUDE.md            |                           | CLAUDE.md            |
| Claude Builtin Tools |                           | Claude Builtin Tools |
| MCP Tools            |                           | MCP Tools            |
+----------------------+                           | Read progress.md     |
| User message         |                           +----------------------+
| Read()               |                           |                      |
| Search()             |   Intentional Compaction  | +-----------------+  |
| Write()              |                           | | progress.md     |  |
+----------------------+  "Write everything we did | +-----------------+  |
| Assistant Message    |   so far to progress.md,  |                      |
+----------------------+   ensure to note the      | Read()               |
| User message         |   approach we're taking,  | Write()              |
| Great, now do ABC    |   the steps we've done    | Read()               |
| Read()               |   so far, and the current | Write()              |
| ...                  |   failure we're working   | ...                  |
+----------------------+   on"                     |                      |
| Assistant Message    |                           |                      |
| ...                  |                           |                      |
| Write()              |                           |                      |
+----------------------+                           |                      |
| User message         |                           |                      |
| Summarize progress   |                           |                      |
| to progress.md       |                           |                      |
+-----------+----------+                           +----------------------+
            |                                                ^
            |             +-----------------+                |
            +------------>| progress.md     |----------------+
                          +-----------------+
```
