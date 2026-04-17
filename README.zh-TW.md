# hostfile

[![GitHub Release](https://img.shields.io/github/v/release/vulcanshen/hostfile)](https://github.com/vulcanshen/hostfile/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/vulcanshen/hostfile)](https://go.dev/)
[![CI](https://img.shields.io/github/actions/workflow/status/vulcanshen/hostfile/ci.yml?label=CI)](https://github.com/vulcanshen/hostfile/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/vulcanshen/hostfile)](https://goreportcard.com/report/github.com/vulcanshen/hostfile)
[![License](https://img.shields.io/github/license/vulcanshen/hostfile)](LICENSE)

[English](README.md) | [日本語](README.ja.md) | [한국어](README.ko.md)

![demo](docs/demo.gif)

跨平台 hosts 檔管理 CLI 工具。

簡單到任何人都能用 — 複製指令、貼上、按 Enter，搞定。

## 功能

- **Add / Remove** — 管理 IP 與 domain 的對應，自動合併相同 IP
- **Enable / Disable** — 停用 / 啟用設定，不需要刪除（可針對整個 IP 或單一 domain）
- **Search / Show** — 查詢與顯示目前的設定，支援著色對齊輸出
- **Apply / Merge** — 從檔案或 stdin 匯入，自動偵測 JSON 和 hosts 格式，含格式驗證
- **JSON I/O** — `show --json` 匯出、`apply -` / `merge -` 透過 pipeline 匯入
- **Save / Load** — 儲存與載入設定快照
- **Clean** — 一鍵清空所有設定
- Managed block 隔離 — 不會動到你手寫的內容
- 自動提權（sudo / doas / gsudo）
- 支援 IPv4 + IPv6（含 zone ID）
- Shell 自動補全（bash, zsh, fish, powershell）
- 跨平台：macOS、Linux、Windows

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

### 一鍵安裝

macOS / Linux / Git Bash：

```bash
curl -fsSL https://raw.githubusercontent.com/vulcanshen/hostfile/main/install.sh | sh
```

Windows（PowerShell）：

```powershell
irm https://raw.githubusercontent.com/vulcanshen/hostfile/main/install.ps1 | iex
```

要更新的話，再執行一次同樣的指令就好。解除安裝：

```bash
curl -fsSL https://raw.githubusercontent.com/vulcanshen/hostfile/main/uninstall.sh | sh
```

```powershell
irm https://raw.githubusercontent.com/vulcanshen/hostfile/main/uninstall.ps1 | iex
```

> **Windows 注意事項**：hostfile 會修改系統 hosts 檔，需要管理員權限。
> Windows 11 24H2 以上版本內建 `sudo`，hostfile 會自動使用。
> 舊版 Windows 請安裝 [gsudo](https://github.com/gerardog/gsudo)，或以系統管理員身分開啟 PowerShell。
>
> 安裝 gsudo：
> ```powershell
> # PowerShell 一鍵安裝
> irm https://raw.githubusercontent.com/gerardog/gsudo/master/installgsudo.ps1 | iex
>
> # 或透過 Scoop
> scoop install gsudo
> ```

### 套件管理器

| 平台 | 指令 |
|------|------|
| Homebrew (macOS / Linux) | `brew install vulcanshen/tap/hostfile` |
| Scoop (Windows) | `scoop bucket add vulcanshen https://github.com/vulcanshen/scoop-bucket && scoop install hostfile` |
| Debian / Ubuntu | `sudo dpkg -i hostfile_<version>_linux_amd64.deb` |
| RHEL / Fedora | `sudo rpm -i hostfile_<version>_linux_amd64.rpm` |

`.deb` 和 `.rpm` 套件可從 [Releases 頁面](https://github.com/vulcanshen/hostfile/releases) 下載。將 `<version>` 替換為版本號（例如 `1.2.0`）。ARM64 系統請將 `linux_amd64` 改為 `linux_arm64`。

## 指令一覽

| 指令 | 說明 |
|------|------|
| `init` | 接管現有 hosts 檔 — 備份原始檔為 "origin"，將所有設定轉為 managed block |
| `add <ip> <domain1> [domain2...]` | 新增 domain 到指定 IP，同 IP 自動合併 |
| `remove <ip\|domain>` | 移除一個 IP（整行）或單一 domain |
| `search <ip\|domain>` | 模糊搜尋 — 不區分大小寫的子字串匹配，搜尋結果高亮顯示 |
| `show` | 顯示目前 managed block 的所有設定（著色對齊） |
| `show --json` | 以 JSON 格式輸出 active 設定 |
| `show <name>` | 顯示某個快照的內容 |
| `enable <ip\|domain>` | 啟用被停用的設定 |
| `disable <ip\|domain>` | 停用設定（不刪除） |
| `apply <file \| ->` | 用檔案或 stdin 取代 managed block（支援 JSON） |
| `merge <file \| ->` | 將檔案或 stdin 合併進 managed block（支援 JSON） |
| `clean` | 清空 managed block 的所有設定 |
| `save <name>` | 儲存目前的設定為快照（存放於 `~/.hostfile/`） |
| `list` | 列出所有已儲存的快照 |
| `load <name>` | 載入快照 |
| `delete <name>` | 刪除快照 |
| `open` | 用預設編輯器開啟 hosts 檔（`$EDITOR`，預設 `vi`；Windows 使用 `notepad`） |
| `version` | 顯示版本號 |

### 全域參數

| 參數 | 說明 |
|------|------|
| `--hosts-file <path>` | 指定 hosts 檔路徑（預設：`/etc/hosts` 或 `C:\Windows\System32\drivers\etc\hosts`） |

### Show 參數

| 參數 | 說明 |
|------|------|
| `--json` | 以 JSON 格式輸出 active 設定（`{"ip": ["domain1", "domain2"]}`） |
| `--all` | 包含 managed block 以外的設定 |

### Search 參數

| 參數 | 說明 |
|------|------|
| `--all` | 包含 managed block 以外的設定 |

## 使用範例

```bash
# 首次使用 — 接管現有的 hosts 檔
hostfile init

# 新增設定
hostfile add 192.168.1.100 web.local api.local

# 查看目前設定
hostfile show
hostfile show --json            # JSON 輸出（僅 active）

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

# 從 JSON 匯入
hostfile apply config.json         # 自動偵測 JSON 格式
hostfile show --json | hostfile apply -  # pipeline 傳輸

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

## 實際應用場景

### 外勤工程師切換網路環境

FAE 和顧問經常在不同地點之間往返 — 公司辦公室、客戶現場、在家。
每個環境都有不同的內部 domain。在 Windows 上改 DNS 設定藏得很深，
而且只能設兩組。每次切換都要重來一遍，既容易出錯又麻煩。

hosts 檔是最通用的解決方案。用 hostfile：

出發去客戶現場

```bash
hostfile save company     # 快照你的辦公室設定
hostfile clean            # 清空所有設定
```

回到辦公室

```bash
hostfile load company     # 一個指令還原
```

### 保持機器乾淨

不要把內部 IP 硬寫進網路設定，
用 hostfile 統一管理。需要的啟用，不需要的停用。
你手寫的設定永遠不會被動到。

### 團隊共享設定

```bash
hostfile show --json | ssh teammate hostfile apply -
```

一行指令同步 hosts 設定到另一台機器。

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

## 進階設定

| 環境變數 | 說明 |
|---------|------|
| `HOSTFILE__HOSTS_FILE` | 覆蓋預設的 hosts 檔路徑。設定後，所有指令都會使用此路徑取代 `/etc/hosts`。 |
| `EDITOR` | `open` 指令使用的編輯器。macOS/Linux 預設 `vi`，Windows 預設 `notepad`。 |

## 授權

[GPL-3.0](LICENSE)
