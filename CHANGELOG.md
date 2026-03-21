# Changelog

All notable changes to this project will be documented in this file. Older entries are archived under `docs/releases/` once we ship tagged versions.

> **Maintenance**: Keep only the 10 most recent releases in reverse-chronological order. Purge older entries when adding new releases.

## [Unreleased]

## [0.1.0] - 2026-03-21

### Added

- **Project scaffolding**: Repository created with dual MIT/Apache-2.0 licensing.
- **Agent guide**: `AGENTS.md` with role catalog, session protocol, and commit attribution standard.
- **Role catalog**: 7 roles (`cxotech`, `devlead`, `devrev`, `infoarch`, `qa`, `releng`, `secrev`) in `config/agentic/roles/`.
- **Governance files**: `MAINTAINERS.md`, `REPOSITORY_SAFETY_PROTOCOLS.md`.
- **Build system**: Makefile with goneat integration, bootstrap, fmt, lint, test, build, release targets.
- **DX tooling**: `.goneat/` configuration (tools.yaml, assess.yaml) for formatting, linting, and quality gates.

### Quality

- `make fmt` verified for this release.
