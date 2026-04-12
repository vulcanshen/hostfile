# hostfile

クロスプラットフォーム hosts ファイル管理 CLI ツール。

チーム向けに設計 — 技術スタッフがコマンドを共有し、非技術メンバー（PM、SA、FAE）はコピー＆ペーストして Enter を押すだけ。

## 機能

- **Add / Remove** — IP とドメインのマッピング管理、同一 IP の自動マージ
- **Enable / Disable** — エントリを削除せずに有効/無効を切り替え（IP 単位またはドメイン単位）
- **Search / Show** — 現在の設定を検索・表示
- **Apply / Merge** — 外部ファイルからインポート（置換またはマージ）
- **Save / Load** — 設定スナップショットの保存と読み込み
- **Clean** — すべての設定を一括クリア
- IPv4 + IPv6 対応
- シェル補完（bash, zsh, fish, powershell）

## 仕組み

hostfile は hosts ファイル内の **managed block** のみを変更し、手書きの内容には一切触れません：

```
# あなたの元の内容 — hostfile はここを変更しません
127.0.0.1  localhost

#### hostfile >>>>>
192.168.1.100  web.company.local api.company.local
#[disable-ip] 192.168.1.200  minio.company.local
#[disable-domain] dockerhand.company.local 192.168.1.100
#### hostfile <<<<<
```

## インストール

### ワンライナーインストール

macOS / Linux / Git Bash：

```bash
curl -fsSL https://raw.githubusercontent.com/vulcanshen/hostfile/main/install.sh | sh
```

Windows（PowerShell）：

```powershell
irm https://raw.githubusercontent.com/vulcanshen/hostfile/main/install.ps1 | iex
```

アップデートは同じコマンドを再実行するだけです。アンインストール：

```bash
curl -fsSL https://raw.githubusercontent.com/vulcanshen/hostfile/main/uninstall.sh | sh
```

```powershell
irm https://raw.githubusercontent.com/vulcanshen/hostfile/main/uninstall.ps1 | iex
```

> **Windows の注意事項**：hostfile はシステムの hosts ファイルを変更するため、管理者権限が必要です。
> Windows 11 24H2 以降は `sudo` が組み込まれており、hostfile が自動的に使用します。
> それ以前のバージョンでは [gsudo](https://github.com/gerardog/gsudo) をインストールするか、PowerShell を管理者として実行してください。

### パッケージマネージャー

| プラットフォーム | コマンド |
|------------------|----------|
| Homebrew (macOS / Linux) | `brew install vulcanshen/tap/hostfile` |
| Scoop (Windows) | `scoop bucket add vulcanshen https://github.com/vulcanshen/scoop-bucket && scoop install hostfile` |
| Debian / Ubuntu | `sudo dpkg -i hostfile_<version>_linux_amd64.deb` |
| RHEL / Fedora | `sudo rpm -i hostfile_<version>_linux_amd64.rpm` |

`.deb` と `.rpm` パッケージは [Releases ページ](https://github.com/vulcanshen/hostfile/releases) からダウンロードできます。`<version>` をバージョン番号（例：`1.3.0`）に置き換えてください。ARM64 システムの場合は `linux_amd64` を `linux_arm64` に変更してください。

## コマンド一覧

| コマンド | 説明 |
|----------|------|
| `hostfile init` | 既存の hosts ファイルを引き継ぐ — 元ファイルを "origin" としてバックアップし、全エントリを managed block に変換 |
| `hostfile add <ip> <domain1> [domain2...]` | 指定 IP にドメインを追加、同一 IP は自動マージ |
| `hostfile remove <ip\|domain>` | IP（行全体）または単一ドメインを削除 |
| `hostfile search <ip\|domain>` | 検索 — IP を入力するとドメインを返し、ドメインを入力すると IP を返す |
| `hostfile show` | managed block の全エントリを表示 |
| `hostfile show <name>` | 保存されたスナップショットの内容を表示 |
| `hostfile enable <ip\|domain>` | 無効化されたエントリを有効化 |
| `hostfile disable <ip\|domain>` | エントリを削除せずに無効化 |
| `hostfile apply <file>` | 外部ファイルで managed block を置換 |
| `hostfile merge <file>` | 外部ファイルを managed block にマージ |
| `hostfile clean` | managed block の全エントリをクリア |
| `hostfile save <name>` | 現在の設定をスナップショットとして保存（`~/.hostfile/` に保存） |
| `hostfile list` | 保存済みスナップショットの一覧表示 |
| `hostfile load <name>` | スナップショットを読み込み |
| `hostfile delete <name>` | スナップショットを削除 |
| `hostfile version` | バージョン番号を表示 |

### グローバルフラグ

| フラグ | 説明 |
|--------|------|
| `--hosts-file <path>` | hosts ファイルのパスを指定（デフォルト：`/etc/hosts` または `C:\Windows\System32\drivers\etc\hosts`） |

## 使用例

```bash
# 初回セットアップ — 既存の hosts ファイルを引き継ぐ
hostfile init

# エントリを追加
hostfile add 192.168.1.100 web.local api.local

# 現在の設定を表示
hostfile show

# 検索
hostfile search web.local
hostfile search 192.168.1.100

# 無効化 / 有効化
hostfile disable web.local        # 単一ドメインを無効化
hostfile disable 192.168.1.100    # IP 全体を無効化
hostfile enable web.local

# 削除
hostfile remove web.local          # ドメインを削除
hostfile remove 192.168.1.100     # IP とその全ドメインを削除

# ファイルからインポート
hostfile apply hosts.txt           # managed block を置換
hostfile merge hosts.txt           # managed block にマージ

# 保存 / 読み込み
hostfile save my-snapshot
hostfile list
hostfile show my-snapshot
hostfile load my-snapshot
hostfile delete my-snapshot

# すべてクリア
hostfile clean

# 元の hosts ファイルを復元（init 前の状態）
hostfile load origin

# 別の hosts ファイルを指定（テスト用）
hostfile show --hosts-file /tmp/test.hosts
```

## シェル補完

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

## ライセンス

[GPL-3.0](LICENSE)
