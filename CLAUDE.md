# Claude Code guide

@AGENTS.md

The file above is the single source of truth for how to work in this
repository (layout, module pattern, commands, conventions, testing, CI).
Keep it up to date instead of duplicating its content here; this file exists
only so Claude Code picks up the guide automatically, and so other agents
that read `AGENTS.md` (OpenAI Codex, Cursor, Windsurf, Aider, GitHub
Copilot, Google Gemini CLI / Jules, Devin, Zed, and others) get the same
instructions Claude does.

This repo also ships Claude Code project skills under `.claude/skills/`,
each wrapping a workflow from `AGENTS.md` into a guided, invokable checklist:

- `scaffolding-crud-resources`: adding a new API resource.
- `creating-database-migrations`: schema changes via Atlas plus `make dao`.
- `adding-environment-variables`: keeping `Env`, `.env.example`, and CI in sync.
- `verifying-changes-before-commit`: lint, build, test, and generated-file checks.

Prefer invoking the relevant skill over running the underlying commands by
hand. `.claude/skills/` is committed on purpose in this repo even though a
machine-wide git ignore rule normally excludes it; see the note in
`.gitignore`.
