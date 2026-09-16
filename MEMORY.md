# Memory

A running log of progress, decisions, and next steps for this repo, meant
for both humans and AI agents picking up work here. Read the most recent
entries before starting a session; add a new entry when you finish one.

Keep entries short and factual. This file is not a replacement for `git
log` or `AGENTS.md`: put durable architecture facts in `AGENTS.md`, not
here, and let this file decay naturally as items move from "next steps" to
done.

## Next steps / open items

- None currently tracked. Add items here as they come up, and remove them
  once done (move a one-line note to the progress log below instead of
  leaving stale checkboxes).

## Decisions

- **2026-09-16**: Skills are authored once under `.agents/skills/` and
  exposed to individual agents via symlinks (`.claude/skills`,
  `.codex/skills`), instead of duplicating SKILL.md files per tool.
  Rationale: one source of truth, no drift between what Claude Code and
  Codex see. Caveat: symlinks need `core.symlinks` support on checkout: on
  Windows without symlink support enabled, `.claude/skills`/`.codex/skills`
  will check out as broken plain-text files instead of directories. Fix
  forward if that ever bites a contributor rather than pre-emptively
  avoiding symlinks.
- **2026-09-16**: The `.claude/skills` symlink disappeared from disk once
  (working tree only; `.agents/skills` and `.codex/skills` were unaffected).
  Cause: it was the only one of the three that was both untracked and
  matched by this machine's global git ignore rule
  (`~/.config/git/ignore`: `**/.claude/skills/`), so something doing an
  ignored-file cleanup between sessions swept it. Fixed by
  `git add -f .claude/skills` so it stays in the index. If it disappears
  again, re-run `ln -s ../.agents/skills .claude/skills && git add -f
  .claude/skills` (do not edit this repo's own `.gitignore` to fix it).
- **2026-09-16**: New documentation (`AGENTS.md`, `CLAUDE.md`, `MEMORY.md`,
  skill files) avoids em dashes; existing docs (`README.md`, `SWAGGER.md`)
  are left as-is rather than reformatted for style.

## Progress log

- **2026-09-16**: Added `AGENTS.md` (agent orchestration / architecture
  guide, following the open AGENTS.md standard) and `CLAUDE.md` (imports
  `AGENTS.md` for Claude Code). Added four skills under `.agents/skills/`:
  `scaffolding-crud-resources`, `creating-database-migrations`,
  `adding-environment-variables`, `verifying-changes-before-commit`.
  Restructured skills from `.claude/skills/` into the canonical
  `.agents/skills/` location with `.claude/skills` and `.codex/skills` as
  symlinks to it. Added this file. No application code or existing
  documentation files were modified.
