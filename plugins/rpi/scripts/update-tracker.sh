#!/usr/bin/env bash
#
# PostToolUse hook for Write tool — auto-updates docs/.tracker.json
# when artifacts are written to docs/research/, docs/plans/, or docs/handoffs/.
#
# Receives JSON on stdin with tool_input.file_path from Claude Code.

set -euo pipefail

# Read hook payload from stdin
payload=$(cat)

# Extract the file path that was written
file_path=$(printf '%s' "$payload" | python3 -c "
import sys, json
data = json.load(sys.stdin)
print(data.get('tool_input', {}).get('file_path', ''))
" 2>/dev/null || echo "")

[ -z "$file_path" ] && exit 0

# Determine project root (look for .git upward from the written file)
dir=$(dirname "$file_path")
project_root=""
while [ "$dir" != "/" ]; do
    if [ -d "$dir/.git" ]; then
        project_root="$dir"
        break
    fi
    dir=$(dirname "$dir")
done
[ -z "$project_root" ] && exit 0

# Get path relative to project root
rel_path="${file_path#"$project_root"/}"

# Only process docs/research/*, docs/plans/*, docs/handoffs/*
case "$rel_path" in
    docs/research/*.md)  artifact_type="research" ;;
    docs/plans/*.md)     artifact_type="plan" ;;
    docs/handoffs/*.md)  artifact_type="handoff" ;;
    docs/handoffs/*/*.md) artifact_type="handoff" ;;
    *)                   exit 0 ;;
esac

tracker="$project_root/docs/.tracker.json"

# Extract a slug from the filename (strip date prefix and extension)
basename_no_ext=$(basename "$file_path" .md)
# Remove YYYY-MM-DD_HH-MM-SS_ prefix if present
slug=$(echo "$basename_no_ext" | sed -E 's/^[0-9]{4}-[0-9]{2}-[0-9]{2}_[0-9]{2}-[0-9]{2}-[0-9]{2}_//')
# Remove ticket prefix if present (e.g., PROJ-1234_)
slug=$(echo "$slug" | sed -E 's/^[A-Z]+-[0-9]+_//')

[ -z "$slug" ] && exit 0

now=$(date -u +"%Y-%m-%dT%H:%M:%S+00:00")

# Create or update tracker using Python for reliable JSON handling
python3 << PYEOF
import json, os, sys

tracker_path = "$tracker"
slug = "$slug"
artifact_type = "$artifact_type"
rel_path = "$rel_path"
now = "$now"

status_map = {
    "research": "research-complete",
    "plan": "plan-complete",
    "handoff": "paused",
}

new_status = status_map.get(artifact_type, "")
if not new_status:
    sys.exit(0)

# Load or create tracker
if os.path.exists(tracker_path):
    with open(tracker_path, "r") as f:
        tracker = json.load(f)
else:
    os.makedirs(os.path.dirname(tracker_path), exist_ok=True)
    tracker = {"version": 1, "tasks": {}}

# Create task entry if missing
if slug not in tracker["tasks"]:
    tracker["tasks"][slug] = {
        "title": slug.replace("-", " ").title(),
        "status": new_status,
        "phase": None,
        "branch": None,
        "created": now,
        "updated": now,
        "artifacts": {
            "research": None,
            "plan": None,
            "handoffs": [],
        },
        "notes": None,
    }
else:
    tracker["tasks"][slug]["status"] = new_status
    tracker["tasks"][slug]["updated"] = now

task = tracker["tasks"][slug]

# Update artifact reference
if artifact_type == "research":
    task["artifacts"]["research"] = rel_path
elif artifact_type == "plan":
    task["artifacts"]["plan"] = rel_path
elif artifact_type == "handoff":
    if rel_path not in task["artifacts"].get("handoffs", []):
        task["artifacts"].setdefault("handoffs", []).append(rel_path)

with open(tracker_path, "w") as f:
    json.dump(tracker, f, indent=2)
    f.write("\n")
PYEOF
