# hostfile

A cross-platform CLI tool for managing your hosts file with ease.

Designed for teams where technical staff guide non-technical members (PMs, SAs, FAEs) through hosts configuration — give them one command, they copy-paste and hit Enter.

## Features

- **Add / Remove** — manage IP-domain mappings, auto-merge same IP entries
- **Enable / Disable** — toggle entries without deleting them (per-IP or per-domain granularity)
- **Search / List** — query and display managed entries
- **Apply / Merge** — import hosts from external files (replace or merge)
- **Backup / Restore** — snapshot and restore your managed entries
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

### Homebrew (macOS / Linux)

```bash
brew install vulcanshen/tap/hostfile
```

### Scoop (Windows)

```powershell
scoop bucket add vulcanshen https://github.com/vulcanshen/scoop-bucket
scoop install hostfile
```

### Download Binary

Download the latest release from the [Releases page](https://github.com/vulcanshen/hostfile/releases).

## Commands

| Command | Description |
|---------|-------------|
| `hostfile init` | Take over the current hosts file — backs up as "origin", reformats all entries into managed block |
| `hostfile add <ip> <domain1> [domain2...]` | Add domains to an IP, auto-merge if the IP already exists |
| `hostfile remove <ip\|domain>` | Remove an IP (entire line) or a single domain |
| `hostfile search <ip\|domain>` | Search the managed block — IP returns domains, domain returns IP |
| `hostfile list` | List all entries in the managed block |
| `hostfile enable <ip\|domain>` | Re-enable a disabled entry |
| `hostfile disable <ip\|domain>` | Disable an entry without deleting it |
| `hostfile apply <file>` | Replace the managed block with content from a file |
| `hostfile merge <file>` | Merge content from a file into the managed block |
| `hostfile clean` | Clear all entries from the managed block |
| `hostfile backup <name>` | Backup the managed block to `~/.hostfile/<name>.hostfile` |
| `hostfile backup list` | List all backups |
| `hostfile backup restore <name>` | Restore the managed block from a backup |
| `hostfile backup delete <name>` | Delete a backup |
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

# List managed entries
hostfile list

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

# Backup
hostfile backup my-snapshot
hostfile backup list
hostfile backup restore my-snapshot
hostfile backup delete my-snapshot

# Clear everything
hostfile clean

# Restore original hosts file (before init)
hostfile backup restore origin

# Use a custom hosts file (for testing or dry-run)
hostfile list --hosts-file /tmp/test.hosts
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

## License

[GPL-3.0](LICENSE)
