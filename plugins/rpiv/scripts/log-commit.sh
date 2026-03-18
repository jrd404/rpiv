#!/usr/bin/env bash
#
# PostToolUse hook for Bash tool — detects git commits and appends a timestamped
# entry to the most recent handoff document in docs/handoffs/.
#
# Receives JSON on stdin with tool_input.command from Claude Code.

set -euo pipefail

payload=$(cat)

# Extract the command that was run
command=$(printf '%s' "$payload" | python3 -c "
import sys, json
data = json.load(sys.stdin)
print(data.get('tool_input', {}).get('command', ''))
" 2>/dev/null || echo "")

[ -z "$command" ] && exit 0

# Only proceed if this was a git commit command
case "$command" in
    *git\ commit*) ;;
    *) exit 0 ;;
esac

# Check exit code — only log successful commits
exit_code=$(printf '%s' "$payload" | python3 -c "
import sys, json
data = json.load(sys.stdin)
r = data.get('tool_response', {})
print(r.get('exit_code', r.get('exitCode', 1)))
" 2>/dev/null || echo "1")

[ "$exit_code" != "0" ] && exit 0

# Determine project root
cwd=$(printf '%s' "$payload" | python3 -c "
import sys, json
data = json.load(sys.stdin)
print(data.get('cwd', ''))
" 2>/dev/null || echo "")

[ -z "$cwd" ] && exit 0

project_root=""
dir="$cwd"
while [ "$dir" != "/" ]; do
    if [ -d "$dir/.git" ]; then
        project_root="$dir"
        break
    fi
    dir=$(dirname "$dir")
done
[ -z "$project_root" ] && exit 0

# Gather commit metadata
cd "$project_root"
commit_hash=$(git rev-parse HEAD 2>/dev/null || echo "unknown")
commit_short=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
repo=$(basename "$project_root")
commit_msg=$(git log -1 --pretty=format:"%s" 2>/dev/null || echo "")
commit_author=$(git log -1 --pretty=format:"%an" 2>/dev/null || echo "")
timestamp=$(date -u +"%Y-%m-%dT%H:%M:%S+00:00")
timestamp_human=$(date +"%Y-%m-%d %H:%M:%S %Z")

# Find the most recent handoff doc
handoffs_dir="$project_root/docs/handoffs"
[ ! -d "$handoffs_dir" ] && exit 0

latest_handoff=$(find "$handoffs_dir" -name '*.md' -type f 2>/dev/null | sort -r | head -1)
[ -z "$latest_handoff" ] && exit 0

# Append commit entry to the handoff doc
cat >> "$latest_handoff" << EOF

---

### Commit Log Entry — ${timestamp_human}

| Field      | Value                        |
|------------|------------------------------|
| Timestamp  | \`${timestamp}\`             |
| Repository | \`${repo}\`                  |
| Branch     | \`${branch}\`                |
| Commit     | \`${commit_short}\` (\`${commit_hash}\`) |
| Author     | ${commit_author}             |
| Message    | ${commit_msg}                |
EOF
