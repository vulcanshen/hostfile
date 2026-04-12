# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

## [v1.2.0] - 2026-04-12

- Rename `backup` → `save`, `backup restore` → `load`, `backup delete` → `delete` — all top-level commands
- Rename `list` → `show`, add `show <name>` to view saved snapshot contents
- New `list` command to list all saved snapshots
- Add CHANGELOG.md with full version history
- Release notes now sourced from CHANGELOG instead of auto-generated commit list
- Add `/update-docs` and `/release` Claude Code skills

## [v1.1.2] - 2026-04-12

- Add tab completion for backup restore and delete
- Make zsh completion setup copy-pasteable
- Add command reference table to README

## [v1.1.1] - 2026-04-12

- Fix IPv6 zone ID parsing (e.g. fe80::1%lo0)
- Add init command to take over existing hosts file
- Fix zsh completion instructions in README

## [v1.0.0] - 2026-04-11

- Initial release: cross-platform hosts file manager CLI
- Add / Remove / Search / List entries
- Enable / Disable entries (per-IP or per-domain)
- Apply / Merge from external files
- Backup / Restore managed block
- Clean managed block
- Shell completion (bash, zsh, fish, powershell)
- Homebrew, Scoop, deb, rpm packaging
