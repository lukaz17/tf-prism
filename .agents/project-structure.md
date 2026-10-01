# Project Structure

## Guideline

- Keep `main()` functions thin: parse flags/config, wire dependencies, then delegate to `internal/` packages.
- Code that should never be imported by other projects belongs under `internal/`.
- Only put code under `pkg/` if it's genuinely meant to be a public, importable API - don't use it as a dumping ground.
- Group files within a package by responsibility, not by type (avoid generic buckets like `utils.go` or `helpers.go` when a more specific name fits).
- Keep test files (`_test.go`) alongside the code they test, in the same package or a `_test` package for black-box tests.

<!-- Project-specific / Guideline -->

## Layout

```
├── .agents/        # shared AI agent instructions (coding conventions, PR/task checklists, project structure)
├── .github/        # GitHub-specific config (workflows, issue/PR templates)
├── .vscode/        # VS Code editor/workspace settings
├── api/            # OpenAPI/Swagger specs, JSON schemas, protocol definitions
├── assets/         # repository assets (images, logos, etc)
├── cmd/            # main packages (one subdirectory per binary)
├── common/         # shared types, helper methods for whole project (project-level stdlib)
├── config/         # global application configuration
├── deployments/    # deployment configs and templates (docker-compose, kubernetes/helm, terraform)
├── docs/           # design and user documents (beyond godoc)
├── examples/       # examples for applications and/or public libraries
├── init/           # system init (systemd, upstart, sysv) and supervisor (runit, supervisord) configs
├── internal/       # private application/library code, not importable by other modules
├── pkg/            # public library code intended for external use (optional)
├── scripts/        # build/install/analysis scripts (keeps root Makefile small)
├── test/           # additional external test apps and test data
├── tui/            # terminal UI components
├── vendor/         # application dependencies (created by `go mod vendor`)
├── web/            # web app components: static assets, server-side templates, SPAs
├── .editorconfig   # editor formatting rules
├── .gitignore      # git ignore patterns
├── AGENTS.md       # instructions for AI coding agents
├── CLAUDE.md       # instructions specific for Claude coding agents
├── CONTRIBUTING.md # contribution guidelines
├── Dockerfile      # container build definition
├── Makefile        # build/test/lint task automation
├── go.mod          # Go module definition
├── go.sum          # Go module checksums
├── main.go         # default application entrypoint
```

<!-- Project-specific / Layout -->

This project add more packages to the standard layout:

```
├── db/             # data models and database access
├── diag/           # logging, instruments, progress tracking for long-running operations
├── engine/         # core application logic wiring CLI commands/controller to business logic
├── server/         # HTTP server for service API and WEB
```

## Project-specific

<!-- Project-specific -->
