<div align="center">

# falseFighter

### File integrity monitoring (FIM) for Linux, written in Go. Work in progress.

**Hash your files. Watch them. Get alerted when something changes.**

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![Status](https://img.shields.io/badge/status-work_in_progress-f4c430?style=flat-square)](#current-status)
[![Platform](https://img.shields.io/badge/platform-Linux-informational?style=flat-square)](#limitations)

**[Quick start](#quick-start) · [How it works](#how-it-works) · [Current status](#current-status) · [Repository layout](#repository-layout)**

</div>

> [!IMPORTANT]
> this is not production software. Use it on machines you own.

---

## What it does

falseFighter takes a snapshot of important files (like `/etc/passwd` and `/etc/shadow`), then tells you when any of them is created, modified, deleted, or has its permissions changed.

It uses only native Go for the core work: no shelling out to `sha256sum` or `find`.

## Quick start

> [!NOTE]
> Only the scanner is built so far. The commands below describe the target for v0.1. See [Current status](#current-status).

### Prerequisites

- Linux
- Go 1.22 or later
- Root access (needed to read files like `/etc/shadow`)

### Build and run

```bash
git clone https://github.com/msakhraoui/falseFighter && cd falseFighter
go build -o hids ./cmd/hids

sudo ./hids init        # build the baseline
sudo ./hids watch       # monitor in real time
sudo ./hids scan        # full re-scan, compare to baseline
sudo ./hids rebaseline  # accept the current state as trusted
```

## How it works

Three phases: **baseline, monitor, respond**.

```
 1. BASELINE            2. MONITOR                   3. RESPOND
 hids init              hids watch / hids scan       alert
 ┌──────────────┐       ┌──────────────────────┐     ┌────────────────┐
 │ walk paths   │       │ fsnotify event       │     │ JSON log line  │
 │ hash + stat  │──────▶│  → snapshot the file │────▶│ (path, old and │
 │ save to JSON │       │  → diff vs baseline  │     │  new hash)     │
 └──────────────┘       └──────────────────────┘     └────────────────┘
        │                         ▲
        ▼                         │
   baseline.json ─────────────────┘
```

- **Baseline:** for each file, store SHA-256, mode, uid, gid, size, mtime, inode.
- **Monitor:** real-time events from fsnotify, plus a periodic full scan as a safety net for missed events.
- **Respond:** every change is logged as a JSON line. The baseline is never updated automatically, only by `hids rebaseline`.

### What counts as a change

| Before | Now | Result |
| :--- | :--- | :--- |
| not in baseline | exists | created |
| in baseline | missing | deleted |
| in baseline | different hash | modified |
| in baseline | same hash, different mode/uid/gid | metadata changed |

## Current status

### Implemented

- `record`: the per-file metadata struct
- `scanner`: `Snapshot()` hashes and stats one file, `Walk()` scans a directory tree and skips unreadable files with a warning

### In progress

- `baseline`: save and load the baseline as JSON

### Planned

- `hids init` command
- `differ`: compare a snapshot against the baseline
- `watcher`: real-time monitoring with fsnotify
- `alert`: JSON alert log
- `parser`: field-level diff for `/etc/passwd` and `/etc/shadow` (new user, UID changed to 0, shell changed)
- Config file (watched paths, excludes)
- systemd service

## Repository layout

```
falseFighter/
├── cmd/hids/          # CLI entry point
├── internal/
│   ├── record/        # shared Record struct
│   ├── scanner/       # Snapshot() and Walk()
│   ├── baseline/      # load and save baseline JSON
│   ├── differ/        # planned
│   ├── watcher/       # planned
│   ├── parser/        # planned
│   ├── alert/         # planned
│   └── config/        # planned
├── go.mod
└── README.md
```

## Limitations

- Linux only (uses `syscall.Stat_t` and inotify).
- Run as root to see root-only files. As a normal user, unreadable files are skipped and logged.
- fsnotify is not recursive, so new subdirectories need to be added at runtime.
- The baseline file can be tampered with by anyone who gets root. Protecting it is on the roadmap.

## Contributing

Feedback is welcome. [Open an issue](https://github.com/msakhraoui/falseFighter/issues) to report a bug or suggest a feature.

## License

[GPL-3.0](LICENSE)
