# hostfile

A cross-platform CLI tool for managing your hosts file with ease.

Designed for teams where technical staff guide non-technical members (PMs, SAs, FAEs) through hosts configuration — give them one command, they copy-paste and hit Enter.

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

### Homebrew (macOS / Linux)

```bash
brew install vulcanshen/tap/hostfile
```

### Scoop (Windows)

```powershell
scoop bucket add vulcanshen https://github.com/vulcanshen/scoop-bucket
scoop install hostfile
```

### Debian / Ubuntu (apt)

```bash
# Download the .deb package from the latest release
curl -LO https://github.com/vulcanshen/hostfile/releases/latest/download/hostfile_<version>_linux_amd64.deb

# Install
sudo dpkg -i hostfile_<version>_linux_amd64.deb
```

Replace `<version>` with the version number (e.g. `1.2.0`). For ARM64 systems, use `linux_arm64.deb` instead.

### RHEL / Fedora (rpm)

```bash
# Download the .rpm package from the latest release
curl -LO https://github.com/vulcanshen/hostfile/releases/latest/download/hostfile_<version>_linux_amd64.rpm

# Install
sudo rpm -i hostfile_<version>_linux_amd64.rpm
```

Replace `<version>` with the version number (e.g. `1.2.0`). For ARM64 systems, use `linux_arm64.rpm` instead.

### Download Binary

Download the archive for your platform from the [Releases page](https://github.com/vulcanshen/hostfile/releases), then extract and move to your PATH:

```bash
# Example for Linux amd64
curl -LO https://github.com/vulcanshen/hostfile/releases/latest/download/hostfile_<version>_linux_amd64.tar.gz
tar xzf hostfile_<version>_linux_amd64.tar.gz
sudo mv hostfile /usr/local/bin/
```

Available platforms: `linux`, `darwin`, `windows` × `amd64`, `arm64`.

## Commands

| Command | Description |
|---------|-------------|
| `hostfile init` | Take over the current hosts file — backs up as "origin", reformats all entries into managed block |
| `hostfile add <ip> <domain1> [domain2...]` | Add domains to an IP, auto-merge if the IP already exists |
| `hostfile remove <ip\|domain>` | Remove an IP (entire line) or a single domain |
| `hostfile search <ip\|domain>` | Search the managed block — IP returns domains, domain returns IP |
| `hostfile show` | Show all entries in the managed block |
| `hostfile show <name>` | Show the contents of a saved snapshot |
| `hostfile enable <ip\|domain>` | Re-enable a disabled entry |
| `hostfile disable <ip\|domain>` | Disable an entry without deleting it |
| `hostfile apply <file>` | Replace the managed block with content from a file |
| `hostfile merge <file>` | Merge content from a file into the managed block |
| `hostfile clean` | Clear all entries from the managed block |
| `hostfile save <name>` | Save the managed block to `~/.hostfile/<name>.hostfile` |
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

## License

[GPL-3.0](LICENSE)
