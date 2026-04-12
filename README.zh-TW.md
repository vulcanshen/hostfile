# hostfile

跨平台 hosts 檔管理 CLI 工具。

專為團隊設計 — 技術人員給一行指令，非技術成員（PM、SA、FAE）複製貼上按 Enter 就搞定。

## 功能

- **Add / Remove** — 管理 IP 與 domain 的對應，自動合併相同 IP
- **Enable / Disable** — 停用 / 啟用設定，不需要刪除（可針對整個 IP 或單一 domain）
- **Search / Show** — 查詢與顯示目前的設定
- **Apply / Merge** — 從外部檔案匯入（取代或合併）
- **Save / Load** — 儲存與載入設定快照
- **Clean** — 一鍵清空所有設定
- 支援 IPv4 + IPv6
- Shell 自動補全（bash, zsh, fish, powershell）

## 運作原理

hostfile 只會修改 hosts 檔中的 **managed block**，不會動到你手寫的內容：

```
# 你原本的內容 — hostfile 不會動這裡
127.0.0.1  localhost

#### hostfile >>>>>
192.168.1.100  web.company.local api.company.local
#[disable-ip] 192.168.1.200  minio.company.local
#[disable-domain] dockerhand.company.local 192.168.1.100
#### hostfile <<<<<
```

## 安裝

### 一鍵安裝（macOS / Linux / Git Bash）

```bash
curl -fsSL https://raw.githubusercontent.com/vulcanshen/hostfile/main/install.sh | sh
```

自動偵測作業系統和架構，下載最新版並安裝到 `~/.local/bin`（root 則安裝到 `/usr/local/bin`）。在 Windows 的 Git Bash 上會安裝到 `~/bin`。

### macOS / Linux（Homebrew）

```bash
brew install vulcanshen/tap/hostfile
```

### Windows（PowerShell 一鍵安裝）

打開 PowerShell，貼上：

```powershell
irm https://raw.githubusercontent.com/vulcanshen/hostfile/main/install.ps1 | iex
```

自動下載最新版、解壓到 `%LOCALAPPDATA%\hostfile`、加入 PATH。安裝完重開終端機即可使用。

要更新的話，再執行一次同樣的指令就好。

> **注意**：hostfile 會修改系統 hosts 檔，需要管理員權限。
> Windows 11 24H2 以上版本內建 `sudo`，hostfile 會自動使用。
> 舊版 Windows 請安裝 [gsudo](https://github.com/gerardog/gsudo)，或以系統管理員身分開啟 PowerShell。

### Windows（Scoop）

如果你已經在用 [Scoop](https://scoop.sh)：

```powershell
scoop bucket add vulcanshen https://github.com/vulcanshen/scoop-bucket
scoop install hostfile
```

### Debian / Ubuntu

```bash
# 下載 .deb 套件
curl -LO https://github.com/vulcanshen/hostfile/releases/latest/download/hostfile_<version>_linux_amd64.deb

# 安裝
sudo dpkg -i hostfile_<version>_linux_amd64.deb
```

將 `<version>` 替換為版本號（例如 `1.2.0`）。ARM64 系統請用 `linux_arm64.deb`。

### RHEL / Fedora

```bash
# 下載 .rpm 套件
curl -LO https://github.com/vulcanshen/hostfile/releases/latest/download/hostfile_<version>_linux_amd64.rpm

# 安裝
sudo rpm -i hostfile_<version>_linux_amd64.rpm
```

將 `<version>` 替換為版本號（例如 `1.2.0`）。ARM64 系統請用 `linux_arm64.rpm`。

### 直接下載

從 [Releases 頁面](https://github.com/vulcanshen/hostfile/releases) 下載對應平台的壓縮檔，解壓後放到 PATH 中：

```bash
# 以 Linux amd64 為例
curl -LO https://github.com/vulcanshen/hostfile/releases/latest/download/hostfile_<version>_linux_amd64.tar.gz
tar xzf hostfile_<version>_linux_amd64.tar.gz
sudo mv hostfile /usr/local/bin/
```

支援平台：`linux`、`darwin`、`windows` × `amd64`、`arm64`

## 指令一覽

| 指令 | 說明 |
|------|------|
| `hostfile init` | 接管現有 hosts 檔 — 備份原始檔為 "origin"，將所有設定轉為 managed block |
| `hostfile add <ip> <domain1> [domain2...]` | 新增 domain 到指定 IP，同 IP 自動合併 |
| `hostfile remove <ip\|domain>` | 移除一個 IP（整行）或單一 domain |
| `hostfile search <ip\|domain>` | 搜尋 — 輸入 IP 回傳 domain，輸入 domain 回傳 IP |
| `hostfile show` | 顯示目前 managed block 的所有設定 |
| `hostfile show <name>` | 顯示某個快照的內容 |
| `hostfile enable <ip\|domain>` | 啟用被停用的設定 |
| `hostfile disable <ip\|domain>` | 停用設定（不刪除） |
| `hostfile apply <file>` | 用外部檔案取代 managed block |
| `hostfile merge <file>` | 將外部檔案合併進 managed block |
| `hostfile clean` | 清空 managed block 的所有設定 |
| `hostfile save <name>` | 儲存目前的設定為快照（存放於 `~/.hostfile/`） |
| `hostfile list` | 列出所有已儲存的快照 |
| `hostfile load <name>` | 載入快照 |
| `hostfile delete <name>` | 刪除快照 |
| `hostfile version` | 顯示版本號 |

### 全域參數

| 參數 | 說明 |
|------|------|
| `--hosts-file <path>` | 指定 hosts 檔路徑（預設：`/etc/hosts` 或 `C:\Windows\System32\drivers\etc\hosts`） |

## 使用範例

```bash
# 首次使用 — 接管現有的 hosts 檔
hostfile init

# 新增設定
hostfile add 192.168.1.100 web.local api.local

# 查看目前設定
hostfile show

# 搜尋
hostfile search web.local
hostfile search 192.168.1.100

# 停用 / 啟用
hostfile disable web.local        # 停用單一 domain
hostfile disable 192.168.1.100    # 停用整個 IP
hostfile enable web.local

# 移除
hostfile remove web.local          # 移除一個 domain
hostfile remove 192.168.1.100     # 移除整個 IP 和其下所有 domain

# 從檔案匯入
hostfile apply hosts.txt           # 取代 managed block
hostfile merge hosts.txt           # 合併進 managed block

# 儲存 / 載入
hostfile save my-snapshot
hostfile list
hostfile show my-snapshot
hostfile load my-snapshot
hostfile delete my-snapshot

# 清空所有設定
hostfile clean

# 還原原始 hosts 檔（init 之前的狀態）
hostfile load origin

# 指定其他 hosts 檔（測試用）
hostfile show --hosts-file /tmp/test.hosts
```

## Shell 自動補全

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

## 授權

[GPL-3.0](LICENSE)
