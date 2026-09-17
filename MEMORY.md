# Memory

A log of durable progress, decisions, and next steps for this repo, meant
for both humans and AI agents picking up work here. Read it before
starting a session.

Update it only for milestones: a notable decision, a completed body of
work, or an open item worth remembering, not a line for every edit in a
session. If in doubt, don't add an entry; `git log` already covers routine
changes.

## Next steps / open items

- None currently tracked.

## Decisions

- **2026-09-17**: Skills are authored once under `.agents/skills/` and
  exposed to individual agents via symlinks (`.claude/skills`,
  `.codex/skills`), instead of duplicating SKILL.md files per tool. One
  source of truth, no drift between what Claude Code and Codex see.
  Caveat: needs `core.symlinks` support on checkout, so `.claude/skills`
  and `.codex/skills` can check out broken on Windows without it.
- **2026-09-17**: This repo intentionally has no `CLAUDE.md`, so Claude
  Code does not auto-load `AGENTS.md` at session start here (an agent
  needs to read it itself). Skills still work regardless, via the
  `.claude/skills` symlink.
- **2026-09-17**: New documentation in this repo avoids em dashes; existing
  docs are not reformatted to remove them.

## Progress log

- **2026-09-17**: Added AI-agent documentation and tooling: `AGENTS.md`
  (architecture, conventions, commands), and four Claude Code skills under
  `.agents/skills/` (`scaffolding-crud-resources`,
  `creating-database-migrations`, `adding-environment-variables`,
  `verifying-changes-before-commit`), exposed via `.claude/skills` and
  `.codex/skills` symlinks. No application code or existing documentation
  was modified.
- **2026-09-17**: Added `make link-skills` (`automate/scripts/link-agent-skills.sh`)
  to (re)create the agent skill symlinks on demand, e.g. after a fresh
  clone or if one goes missing; accepts extra agent names as arguments
  (e.g. `make link-skills cursor`) so a new agent's directory can be
  linked without editing the script.
