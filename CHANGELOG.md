# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

## [v1.4.1] - 2026-04-12

- Add stdin support for `apply` and `merge` — use `-` to read from pipe
- Add JSON input support for `apply` and `merge` — auto-detected
- Add input validation for `apply` and `merge` — rejects empty, invalid JSON, bad IPs, and unparseable hosts format before touching the hosts file

## [v1.4.0] - 2026-04-12

- Add colored and aligned output for `show` and `search` commands (auto-disabled when piped)
- Add `--json` flag to `show` command — outputs active entries as JSON (key: IP, value: domain array)
- Add `HOSTFILE__HOSTS_FILE` environment variable to override default hosts file path
- Add `DisplayPath()` — file paths now display `~` instead of absolute home directory
- Add duplicate domain detection in `add` command — shows "already exists" instead of false "added"
- Add demo GIF (`vhs` tape script included)
- Add uninstall scripts (`uninstall.sh`, `uninstall.ps1`)

## [v1.3.0] - 2026-04-12

- Add man pages (auto-generated from cobra commands, included in release archives and deb/rpm)
- Add tldr page (`docs/tldr.md`)
- Add uninstall scripts (`uninstall.sh`, `uninstall.ps1`)
- Add duplicate domain detection in `add` command — shows "already exists" instead of false "added"
- Update install.sh to install man pages on Unix systems

## [v1.2.3] - 2026-04-12

- Fix `writeBlock()` missing newline guard — prevents block marker from merging with preceding content
- Fix `sudo tee` stdout leakage on Unix — output no longer echoed to terminal
- Fix Windows privilege escalation — use PowerShell instead of `tee` (which doesn't exist on Windows)

## [v1.2.2] - 2026-04-12

### Bug Fixes

- Fix `load` command missing privilege escalation for non-raw snapshots
- Fix ignored error from `backup.Path()` in `save` and `init` commands
- Remove duplicate `containsDomain()` function, consolidate into `manager.ContainsDomain()`

### Improvements

- Add `scanner.Err()` check in prompt input handling
- Rename error messages from "backup" to "save" to match new command names
- Unify output message format across all commands

### Tests

- Add tests for `backup.CreateRaw`, `RestoreRaw`, `IsRaw`, `Exists`, `Path`
- Add tests for `AddForce` with disabled entries, edge cases for `Remove`, `Search`, `Enable`, `Disable`
- Add tests for `ContainsDomain` and `Merge` with disabled entries

## [v1.2.1] - 2026-04-12

- Add PowerShell one-liner installer for Windows (`install.ps1`)
- Add shell installer for macOS / Linux / Git Bash (`install.sh`)
- Add Traditional Chinese README (`README.zh-TW.md`)
- Add Windows administrator privileges note to installation docs
- Improve `init` to merge outside entries when managed block already exists
- Auto-generate timestamped origin name (`origin-YYYYMMDD_HHMMSS`) on repeated init
- Show full file path after `save` and `init`
- Reorganize installation docs into Quick Install and Package Managers sections

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
