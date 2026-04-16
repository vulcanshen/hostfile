# hostfile

[![GitHub Release](https://img.shields.io/github/v/release/vulcanshen/hostfile)](https://github.com/vulcanshen/hostfile/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/vulcanshen/hostfile)](https://go.dev/)
[![CI](https://img.shields.io/github/actions/workflow/status/vulcanshen/hostfile/ci.yml?label=CI)](https://github.com/vulcanshen/hostfile/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/vulcanshen/hostfile)](https://goreportcard.com/report/github.com/vulcanshen/hostfile)
[![License](https://img.shields.io/github/license/vulcanshen/hostfile)](LICENSE)

[English](README.md) | [繁體中文](README.zh-TW.md) | [한국어](README.ko.md)

![demo](docs/demo.gif)

クロスプラットフォーム hosts ファイル管理 CLI ツール。

誰でも使えるほどシンプル — コマンドをコピーして貼り付け、Enter を押すだけ。

## 機能

- **Add / Remove** — IP とドメインのマッピング管理、同一 IP の自動マージ
- **Enable / Disable** — エントリを削除せずに有効/無効を切り替え（IP 単位またはドメイン単位）
- **Search / Show** — 現在の設定を検索・表示、カラー＆整列出力対応
- **Apply / Merge** — ファイルまたは stdin からインポート、JSON と hosts 形式を自動検出＆バリデーション
- **JSON I/O** — `show --json` でエクスポート、`apply -` / `merge -` でパイプラインからインポート
- **Save / Load** — 設定スナップショットの保存と読み込み
- **Clean** — すべての設定を一括クリア
- Managed block 分離 — 手書きの内容には一切触れません
- 自動権限昇格（sudo / doas / gsudo）
- IPv4 + IPv6 対応（ゾーン ID 含む）
- シェル補完（bash, zsh, fish, powershell）
- クロスプラットフォーム：macOS、Linux、Windows

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
>
> gsudo のインストール：
> ```powershell
> # PowerShell ワンライナー
> irm https://raw.githubusercontent.com/gerardog/gsudo/master/installgsudo.ps1 | iex
>
> # または Scoop
> scoop install gsudo
> ```

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
| `init` | 既存の hosts ファイルを引き継ぐ — 元ファイルを "origin" としてバックアップし、全エントリを managed block に変換 |
| `add <ip> <domain1> [domain2...]` | 指定 IP にドメインを追加、同一 IP は自動マージ |
| `remove <ip\|domain>` | IP（行全体）または単一ドメインを削除 |
| `search <ip\|domain>` | あいまい検索 — 大文字小文字を区別しない部分文字列マッチング、検索結果をハイライト表示 |
| `show` | managed block の全エントリを表示（カラー＆整列） |
| `show --json` | アクティブなエントリを JSON で出力 |
| `show <name>` | 保存されたスナップショットの内容を表示 |
| `enable <ip\|domain>` | 無効化されたエントリを有効化 |
| `disable <ip\|domain>` | エントリを削除せずに無効化 |
| `apply <file \| ->` | ファイルまたは stdin で managed block を置換（JSON 対応） |
| `merge <file \| ->` | ファイルまたは stdin を managed block にマージ（JSON 対応） |
| `clean` | managed block の全エントリをクリア |
| `save <name>` | 現在の設定をスナップショットとして保存（`~/.hostfile/` に保存） |
| `list` | 保存済みスナップショットの一覧表示 |
| `load <name>` | スナップショットを読み込み |
| `delete <name>` | スナップショットを削除 |
| `version` | バージョン番号を表示 |

### グローバルフラグ

| フラグ | 説明 |
|--------|------|
| `--hosts-file <path>` | hosts ファイルのパスを指定（デフォルト：`/etc/hosts` または `C:\Windows\System32\drivers\etc\hosts`） |

### Show フラグ

| フラグ | 説明 |
|--------|------|
| `--json` | アクティブなエントリを JSON で出力（`{"ip": ["domain1", "domain2"]}`） |
| `--all` | managed block 外のエントリも含める |

### Search フラグ

| フラグ | 説明 |
|--------|------|
| `--all` | managed block 外のエントリも含める |

## 使用例

```bash
# 初回セットアップ — 既存の hosts ファイルを引き継ぐ
hostfile init

# エントリを追加
hostfile add 192.168.1.100 web.local api.local

# 現在の設定を表示
hostfile show
hostfile show --json            # JSON 出力（アクティブのみ）

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

# JSON からインポート
hostfile apply config.json         # JSON 形式を自動検出
hostfile show --json | hostfile apply -  # パイプラインで転送

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

## 実際の活用シーン

### フィールドエンジニアのネットワーク環境切り替え

FAE やコンサルタントは複数の拠点を行き来します — 自社オフィス、顧客先、自宅。
環境ごとに異なる内部ドメインがあります。Windows で DNS 設定を変更するのは
深い階層に埋もれていて、しかも2つしか設定できません。切り替えるたびにやり直すのは
ミスしやすく面倒です。

hosts ファイルが最もポータブルな解決策です。hostfile なら：

顧客先へ出発

```bash
hostfile save company     # オフィスの設定をスナップショット
hostfile clean            # すべてクリア
```

オフィスに戻ったら

```bash
hostfile load company     # 一つのコマンドで復元
```

### マシンをクリーンに保つ

内部 IP をネットワーク設定にハードコーディングする代わりに、
hostfile で一元管理。必要なものを有効化、不要なものを無効化。
手書きの設定は一切変更されません。

### チームで設定を共有

```bash
hostfile show --json | ssh teammate hostfile apply -
```

ワンライナーで hosts 設定を別のマシンに同期。

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

## 詳細設定

| 環境変数 | 説明 |
|----------|------|
| `HOSTFILE__HOSTS_FILE` | デフォルトの hosts ファイルパスをオーバーライド。設定すると、すべてのコマンドが `/etc/hosts` の代わりにこのパスを使用します。 |

## ライセンス

[GPL-3.0](LICENSE)
