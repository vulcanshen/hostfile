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

## Usage

```bash
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

# Use a custom hosts file (for testing or dry-run)
hostfile list --hosts-file /tmp/test.hosts
```

## Shell Completion

```bash
# Bash
hostfile completion bash > /etc/bash_completion.d/hostfile

# Zsh
hostfile completion zsh > "${fpath[1]}/_hostfile"

# Fish
hostfile completion fish > ~/.config/fish/completions/hostfile.fish

# PowerShell
hostfile completion powershell > hostfile.ps1
```

## License

[GPL-3.0](LICENSE)
