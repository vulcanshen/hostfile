# hostfile

[![GitHub Release](https://img.shields.io/github/v/release/vulcanshen/hostfile)](https://github.com/vulcanshen/hostfile/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/vulcanshen/hostfile)](https://go.dev/)
[![CI](https://img.shields.io/github/actions/workflow/status/vulcanshen/hostfile/ci.yml?label=CI)](https://github.com/vulcanshen/hostfile/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/vulcanshen/hostfile)](https://goreportcard.com/report/github.com/vulcanshen/hostfile)
[![License](https://img.shields.io/github/license/vulcanshen/hostfile)](LICENSE)

[English](README.md) | [繁體中文](README.zh-TW.md) | [日本語](README.ja.md)

![demo](docs/demo.gif)

크로스 플랫폼 hosts 파일 관리 CLI 도구.

누구나 사용할 수 있을 만큼 간단합니다 — 명령어를 복사하고 붙여넣고 Enter만 누르면 됩니다.

## 기능

- **Add / Remove** — IP와 도메인 매핑 관리, 동일 IP 자동 병합
- **Enable / Disable** — 항목을 삭제하지 않고 활성화/비활성화 전환 (IP 단위 또는 도메인 단위)
- **Search / Show** — 현재 설정 검색 및 표시, 컬러 및 정렬 출력 지원
- **Apply / Merge** — 파일 또는 stdin에서 가져오기, JSON 및 hosts 형식 자동 감지 및 검증
- **JSON I/O** — `show --json`으로 내보내기, `apply -` / `merge -`로 파이프라인에서 가져오기
- **Save / Load** — 설정 스냅샷 저장 및 불러오기
- **Clean** — 모든 설정을 한 번에 초기화
- Managed block 격리 — 수동으로 작성한 내용은 절대 건드리지 않음
- 자동 권한 상승 (sudo / doas / gsudo)
- IPv4 + IPv6 지원 (zone ID 포함)
- 셸 자동 완성 (bash, zsh, fish, powershell)
- 크로스 플랫폼: macOS, Linux, Windows

## 작동 방식

hostfile은 hosts 파일 내의 **managed block**만 수정하며, 수동으로 작성한 내용은 절대 건드리지 않습니다:

```
# 원래 내용 — hostfile은 여기를 수정하지 않습니다
127.0.0.1  localhost

#### hostfile >>>>>
192.168.1.100  web.company.local api.company.local
#[disable-ip] 192.168.1.200  minio.company.local
#[disable-domain] dockerhand.company.local 192.168.1.100
#### hostfile <<<<<
```

## 설치

### 원라인 설치

macOS / Linux / Git Bash:

```bash
curl -fsSL https://raw.githubusercontent.com/vulcanshen/hostfile/main/install.sh | sh
```

Windows (PowerShell):

```powershell
irm https://raw.githubusercontent.com/vulcanshen/hostfile/main/install.ps1 | iex
```

업데이트는 같은 명령어를 다시 실행하면 됩니다. 제거:

```bash
curl -fsSL https://raw.githubusercontent.com/vulcanshen/hostfile/main/uninstall.sh | sh
```

```powershell
irm https://raw.githubusercontent.com/vulcanshen/hostfile/main/uninstall.ps1 | iex
```

> **Windows 참고**: hostfile은 시스템 hosts 파일을 수정하므로 관리자 권한이 필요합니다.
> Windows 11 24H2 이상에서는 `sudo`가 내장되어 있으며 hostfile이 자동으로 사용합니다.
> 이전 버전에서는 [gsudo](https://github.com/gerardog/gsudo)를 설치하거나 PowerShell을 관리자로 실행하세요.
>
> gsudo 설치:
> ```powershell
> # PowerShell 원라인
> irm https://raw.githubusercontent.com/gerardog/gsudo/master/installgsudo.ps1 | iex
>
> # 또는 Scoop
> scoop install gsudo
> ```

### 패키지 관리자

| 플랫폼 | 명령어 |
|--------|--------|
| Homebrew (macOS / Linux) | `brew install vulcanshen/tap/hostfile` |
| Scoop (Windows) | `scoop bucket add vulcanshen https://github.com/vulcanshen/scoop-bucket && scoop install hostfile` |
| Debian / Ubuntu | `sudo dpkg -i hostfile_<version>_linux_amd64.deb` |
| RHEL / Fedora | `sudo rpm -i hostfile_<version>_linux_amd64.rpm` |

`.deb` 및 `.rpm` 패키지는 [Releases 페이지](https://github.com/vulcanshen/hostfile/releases)에서 다운로드할 수 있습니다. `<version>`을 버전 번호(예: `1.3.0`)로 바꾸세요. ARM64 시스템은 `linux_amd64` 대신 `linux_arm64`를 사용하세요.

## 명령어 목록

| 명령어 | 설명 |
|--------|------|
| `init` | 기존 hosts 파일 인수 — 원본을 "origin"으로 백업하고 모든 항목을 managed block으로 변환 |
| `add <ip> <domain1> [domain2...]` | 지정 IP에 도메인 추가, 동일 IP 자동 병합 |
| `remove <ip\|domain>` | IP(전체 행) 또는 단일 도메인 제거 |
| `search <ip\|domain>` | 퍼지 검색 — 대소문자 구분 없는 부분 문자열 매칭, 검색 결과 하이라이트 표시 |
| `show` | managed block의 모든 항목 표시 (컬러 및 정렬) |
| `show --json` | 활성 항목을 JSON으로 출력 |
| `show <name>` | 저장된 스냅샷의 내용 표시 |
| `enable <ip\|domain>` | 비활성화된 항목 활성화 |
| `disable <ip\|domain>` | 항목을 삭제하지 않고 비활성화 |
| `apply <file \| ->` | 파일 또는 stdin으로 managed block 대체 (JSON 지원) |
| `merge <file \| ->` | 파일 또는 stdin을 managed block에 병합 (JSON 지원) |
| `clean` | managed block의 모든 항목 초기화 |
| `save <name>` | 현재 설정을 스냅샷으로 저장 (`~/.hostfile/`에 저장) |
| `list` | 저장된 스냅샷 목록 표시 |
| `load <name>` | 스냅샷 불러오기 |
| `delete <name>` | 스냅샷 삭제 |
| `version` | 버전 번호 표시 |

### 글로벌 플래그

| 플래그 | 설명 |
|--------|------|
| `--hosts-file <path>` | hosts 파일 경로 지정 (기본값: `/etc/hosts` 또는 `C:\Windows\System32\drivers\etc\hosts`) |

### Show 플래그

| 플래그 | 설명 |
|--------|------|
| `--json` | 활성 항목을 JSON으로 출력 (`{"ip": ["domain1", "domain2"]}`) |
| `--all` | managed block 외부 항목도 포함 |

### Search 플래그

| 플래그 | 설명 |
|--------|------|
| `--all` | managed block 외부 항목도 포함 |

## 사용 예시

```bash
# 최초 설정 — 기존 hosts 파일 인수
hostfile init

# 항목 추가
hostfile add 192.168.1.100 web.local api.local

# 현재 설정 표시
hostfile show
hostfile show --json            # JSON 출력 (활성만)

# 검색
hostfile search web.local
hostfile search 192.168.1.100

# 비활성화 / 활성화
hostfile disable web.local        # 단일 도메인 비활성화
hostfile disable 192.168.1.100    # IP 전체 비활성화
hostfile enable web.local

# 제거
hostfile remove web.local          # 도메인 제거
hostfile remove 192.168.1.100     # IP와 모든 도메인 제거

# 파일에서 가져오기
hostfile apply hosts.txt           # managed block 대체
hostfile merge hosts.txt           # managed block에 병합

# JSON에서 가져오기
hostfile apply config.json         # JSON 형식 자동 감지
hostfile show --json | hostfile apply -  # 파이프라인으로 전송

# 저장 / 불러오기
hostfile save my-snapshot
hostfile list
hostfile show my-snapshot
hostfile load my-snapshot
hostfile delete my-snapshot

# 모두 초기화
hostfile clean

# 원본 hosts 파일 복원 (init 이전 상태)
hostfile load origin

# 다른 hosts 파일 지정 (테스트용)
hostfile show --hosts-file /tmp/test.hosts
```

## 실제 활용 시나리오

### 현장 엔지니어의 네트워크 환경 전환

FAE와 컨설턴트는 여러 장소를 오갑니다 — 회사, 고객사, 재택.
환경마다 다른 내부 도메인이 있습니다. Windows에서 DNS 설정을 변경하는 것은
깊이 묻혀 있고 항목도 2개로 제한됩니다. 전환할 때마다 다시 설정하는 것은
실수하기 쉽고 번거롭습니다.

hosts 파일이 가장 범용적인 해결책입니다. hostfile을 사용하면:

고객사로 출발

```bash
hostfile save company     # 사무실 설정을 스냅샷
hostfile clean            # 모두 초기화
```

사무실로 복귀

```bash
hostfile load company     # 한 명령어로 복원
```

### 깨끗한 머신 유지

내부 IP를 네트워크 설정에 하드코딩하는 대신,
hostfile로 통합 관리하세요. 필요한 것은 활성화, 필요 없는 것은 비활성화.
수동으로 작성한 설정은 절대 변경되지 않습니다.

### 팀 간 설정 공유

```bash
hostfile show --json | ssh teammate hostfile apply -
```

한 줄 명령으로 hosts 설정을 다른 머신에 동기화.

## 셸 자동 완성

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

## 고급 설정

| 환경 변수 | 설명 |
|-----------|------|
| `HOSTFILE__HOSTS_FILE` | 기본 hosts 파일 경로를 오버라이드. 설정하면 모든 명령어가 `/etc/hosts` 대신 이 경로를 사용합니다. |

## 라이선스

[GPL-3.0](LICENSE)
