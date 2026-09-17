#!/bin/bash

# Re-creates the per-agent skill symlinks so every coding agent reads the
# same skills from one canonical location: .agents/skills/. Safe to
# re-run any time, e.g. after a fresh clone or if a symlink ever goes
# missing.
#
# Links .claude and .codex by default. Pass any additional agent
# directory names as arguments to link those too, for example:
#   make link-skills cursor .windsurf
# (the leading dot is optional; it's added automatically if missing).

set -e

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$REPO_ROOT"

CANONICAL="./.agents/skills"
DEFAULT_AGENT_DIRS=(".claude" ".codex")

if [ ! -d "$CANONICAL" ]; then
  echo "No $CANONICAL directory found. Create your skills there first."
  exit 1
fi

AGENT_DIRS=("${DEFAULT_AGENT_DIRS[@]}")
for extra in "$@"; do
  case "$extra" in
    .*) AGENT_DIRS+=("$extra") ;;
    *)  AGENT_DIRS+=(".$extra") ;;
  esac
done

for dir in "${AGENT_DIRS[@]}"; do
  link="$dir/skills"
  mkdir -p "$dir"

  if [ -L "$link" ] && [ "$(readlink "$link")" = "../.agents/skills" ]; then
    echo "OK      $link already links to ../.agents/skills"
    continue
  fi

  if [ -e "$link" ] && [ ! -L "$link" ]; then
    echo "SKIP    $link exists and is not a symlink; remove it manually first if you want it replaced."
    continue
  fi

  rm -f "$link"
  ln -s "../.agents/skills" "$link"
  echo "LINKED  $link -> ../.agents/skills"
done

echo ""
echo "If a link doesn't show up in 'git status', it may be covered by a"
echo "machine-local git ignore rule; force-add it with: git add -f <path>"
