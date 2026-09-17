---
name: adding-environment-variables
description: Adds a new environment variable to this Go Gin skeleton's config, keeping the Env struct, .env.example, and CI in sync. Use when a change needs a new config value, secret, or feature flag read from the environment.
---

# Adding an environment variable

Every environment variable in this project flows through one struct, so it
must be added in more than one place or it silently reads as empty.

## Steps

1. Add a field to the `Env` struct in `lib/config/env.go`, with a
   `mapstructure` tag matching the variable name exactly:

   ```go
   MyNewSetting string `mapstructure:"MY_NEW_SETTING"`
   ```

   Add a `validate:"required"` (or other validator tag) if the app should
   fail fast at startup when the value is missing, matching the pattern
   already used for `DBUsername`, `JwtAccessSecret`, and similar fields.

2. Add the variable to `.env.example`, with a comment explaining what it is
   for and any accepted values, placed near related variables (for example
   alongside `IDEMPOTENCY_STORE` for a cache setting, or the JWT block for
   an auth setting).

3. If local development needs a real value, add it to your own `.env` (not
   committed).

4. If CircleCI needs the variable (most runtime config used only inside the
   app does not), add it to the env file written in `.circleci/config.yml`'s
   `Initialize the environment variable file` step, sourced from a CircleCI
   context/env var.

5. Use `env.MyNewSetting` from the `config.Env` struct wherever the value is
   needed; do not read `os.Getenv` directly elsewhere in the app.

6. Run the verifying-changes-before-commit skill before treating the change
   as done.
