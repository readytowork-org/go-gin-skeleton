# Claude Code guide

@AGENTS.md

The file above is the single source of truth for how to work in this
repository (layout, module pattern, commands, conventions, testing, CI).
Keep it up to date instead of duplicating its content here; this file exists
only so Claude Code picks up the guide automatically, and so other agents
that read `AGENTS.md` (OpenAI Codex, Cursor, Windsurf, Aider, GitHub
Copilot, Google Gemini CLI / Jules, Devin, Zed, and others) get the same
instructions Claude does.

See `MEMORY.md` for a running log of progress, decisions, and next steps.
Read it at the start of a session and update it as you work.

## Skills

The canonical skill definitions live under `.agents/skills/`, one workflow
from `AGENTS.md` per skill:

- `scaffolding-crud-resources`: adding a new API resource.
- `creating-database-migrations`: schema changes via Atlas plus `make dao`.
- `adding-environment-variables`: keeping `Env`, `.env.example`, and CI in sync.
- `verifying-changes-before-commit`: lint, build, test, and generated-file checks.

`.claude/skills` and `.codex/skills` are symlinks to `.agents/skills`, so
Claude Code and Codex (and any other agent that reads one of those
locations) discover the same skills from one source. Edit a skill only
under `.agents/skills/<name>/SKILL.md`; do not edit through the symlinks
or duplicate content into `.claude/` or `.codex/` directly.

Prefer invoking the relevant skill over running the underlying commands by
hand.
