# hostfile

[![GitHub Release](https://img.shields.io/github/v/release/vulcanshen/hostfile)](https://github.com/vulcanshen/hostfile/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/vulcanshen/hostfile)](https://go.dev/)
[![CI](https://img.shields.io/github/actions/workflow/status/vulcanshen/hostfile/release.yml?label=CI)](https://github.com/vulcanshen/hostfile/actions)
[![License](https://img.shields.io/github/license/vulcanshen/hostfile)](LICENSE)

[繁體中文](README.zh-TW.md) | [日本語](README.ja.md) | [한국어](README.ko.md)

![demo](docs/demo.gif)

A cross-platform CLI tool for managing your hosts file with ease.

Designed to be simple enough that anyone can use it — just copy-paste a command and hit Enter.

## Features

- **Add / Remove** — manage IP-domain mappings, auto-merge same IP entries
- **Enable / Disable** — toggle entries without deleting them (per-IP or per-domain granularity)
- **Search / Show** — query and display managed entries
- **Apply / Merge** — import hosts from external files (replace or merge)
- **Save / Load** — snapshot and restore your managed entries
- **Clean** — clear all managed entries in one command
- IPv4 + IPv6 support
- Shell completion (bash, zsh, fish, powershell)

## How It Works

hostfile only touches its own **managed block** inside your hosts file — it never modifies entries you wrote by hand:

```
# Your original content — hostfile won't touch this
127.0.0.1  localhost

#### hostfile >>>>>
192.168.1.100  web.company.local api.company.local
#[disable-ip] 192.168.1.200  minio.company.local
#[disable-domain] dockerhand.company.local 192.168.1.100
#### hostfile <<<<<
```

## Installation

### Quick Install

macOS / Linux / Git Bash:

```bash
curl -fsSL https://raw.githubusercontent.com/vulcanshen/hostfile/main/install.sh | sh
```

Windows (PowerShell):

```powershell
irm https://raw.githubusercontent.com/vulcanshen/hostfile/main/install.ps1 | iex
```

To update, run the same command again. To uninstall:

```bash
curl -fsSL https://raw.githubusercontent.com/vulcanshen/hostfile/main/uninstall.sh | sh
```

```powershell
irm https://raw.githubusercontent.com/vulcanshen/hostfile/main/uninstall.ps1 | iex
```

> **Windows Note**: hostfile modifies the system hosts file, which requires administrator privileges.
> On Windows 11 24H2+, `sudo` is built-in and hostfile will use it automatically.
> On older versions, either install [gsudo](https://github.com/gerardog/gsudo) or run PowerShell as Administrator.

### Package Managers

| Platform | Command |
|----------|---------|
| Homebrew (macOS / Linux) | `brew install vulcanshen/tap/hostfile` |
| Scoop (Windows) | `scoop bucket add vulcanshen https://github.com/vulcanshen/scoop-bucket && scoop install hostfile` |
| Debian / Ubuntu | `sudo dpkg -i hostfile_<version>_linux_amd64.deb` |
| RHEL / Fedora | `sudo rpm -i hostfile_<version>_linux_amd64.rpm` |

`.deb` and `.rpm` packages can be downloaded from the [Releases page](https://github.com/vulcanshen/hostfile/releases). Replace `<version>` with the version number (e.g. `1.2.0`). For ARM64 systems, use `linux_arm64` instead of `linux_amd64`.

## Commands

| Command | Description |
|---------|-------------|
| `hostfile init` | Take over the current hosts file — backs up as "origin", reformats all entries into managed block |
| `hostfile add <ip> <domain1> [domain2...]` | Add domains to an IP, auto-merge if the IP already exists |
| `hostfile remove <ip\|domain>` | Remove an IP (entire line) or a single domain |
| `hostfile search <ip\|domain>` | Search the managed block — IP returns domains, domain returns IP |
| `hostfile show` | Show all entries in the managed block (colored, aligned) |
| `hostfile show --json` | Output active entries as JSON |
| `hostfile show <name>` | Show the contents of a saved snapshot |
| `hostfile enable <ip\|domain>` | Re-enable a disabled entry |
| `hostfile disable <ip\|domain>` | Disable an entry without deleting it |
| `hostfile apply <file \| ->` | Replace the managed block with content from a file or stdin (supports JSON) |
| `hostfile merge <file \| ->` | Merge content from a file or stdin into the managed block (supports JSON) |
| `hostfile clean` | Clear all entries from the managed block |
| `hostfile save <name>` | Save the managed block as a snapshot to `~/.hostfile/<name>.hostfile` |
| `hostfile list` | List all saved snapshots |
| `hostfile load <name>` | Load a saved snapshot into the managed block |
| `hostfile delete <name>` | Delete a saved snapshot |
| `hostfile version` | Print the version number |

### Global Flags

| Flag | Description |
|------|-------------|
| `--hosts-file <path>` | Path to hosts file (default: `/etc/hosts` or `C:\Windows\System32\drivers\etc\hosts`) |

## Usage Examples

```bash
# First time setup — take over your existing hosts file
hostfile init

# Add entries
hostfile add 192.168.1.100 web.local api.local

# Show managed entries
hostfile show
hostfile show --json            # JSON output (active only)

# Search
hostfile search web.local
hostfile search 192.168.1.100

# Disable / Enable
hostfile disable web.local        # disable a single domain
hostfile disable 192.168.1.100    # disable an entire IP line
hostfile enable web.local

# Remove
hostfile remove web.local          # remove a domain
hostfile remove 192.168.1.100     # remove an IP and all its domains

# Import from file
hostfile apply hosts.txt           # replace managed block
hostfile merge hosts.txt           # merge into managed block

# Import from JSON
hostfile apply config.json         # auto-detects JSON format
hostfile show --json | hostfile apply -  # pipe between instances

# Save / Load
hostfile save my-snapshot
hostfile list
hostfile show my-snapshot
hostfile load my-snapshot
hostfile delete my-snapshot

# Clear everything
hostfile clean

# Restore original hosts file (before init)
hostfile load origin

# Use a custom hosts file (for testing or dry-run)
hostfile show --hosts-file /tmp/test.hosts
```

## Shell Completion

```bash
# Zsh
mkdir -p ~/.zsh/completions
hostfile completion zsh > ~/.zsh/completions/_hostfile
echo 'fpath=(~/.zsh/completions $fpath)' >> ~/.zshrc
echo 'autoload -Uz compinit && compinit' >> ~/.zshrc
source ~/.zshrc

# Bash
hostfile completion bash > /etc/bash_completion.d/hostfile

# Fish
hostfile completion fish > ~/.config/fish/completions/hostfile.fish

# PowerShell
hostfile completion powershell > hostfile.ps1
```

## Advanced

| Environment Variable | Description |
|---------------------|-------------|
| `HOSTFILE__HOSTS_FILE` | Override the default hosts file path. When set, all commands use this path instead of `/etc/hosts`. |

## License

[GPL-3.0](LICENSE)
