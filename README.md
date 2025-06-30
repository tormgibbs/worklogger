# WorkLogger

WorkLogger is a Go-based CLI tool for tracking work sessions, integrating with GitHub for commit logging, and visualizing data via a TUI or web interface. It uses SQLite to store tasks, sessions, and commits, with CSV export for analysis.

## Features
- Track tasks with `start`, `pause`, `resume`, and `stop` commands.
- Sync Git commits to sessions with automatic hook (`setup-hook`) or manual sync (`sync`).
- Authenticate via GitHub OAuth or local credentials (`signup`, `login`, `logout`).
- View logs in an interactive TUI (`log`) or web interface (`studio`).
- Export session data to CSV (`export`).
- View productivity stats (`summary`).

## Installation
1. Clone the repo:
   ```bash
   git clone https://github.com/tormgibbs/worklogger
   cd worklogger
   ```
2. Build the frontend and binary:
   ```bash
   make build
   ```
3. Move the binary to PATH (optional):
   ```bash
   mv worklogger /usr/local/bin/
   ```

## Prerequisites
- Go 1.21+
- Git
- SQLite (`mattn/go-sqlite3`)
- pnpm (for frontend build)
- System keyring (for OAuth credentials)

## Usage
Initialize the environment:
```bash
worklogger init
```
Follow prompts to set up SQLite and GitHub OAuth.

Key commands:
- Start a task: `worklogger start --task "Write code"`
- Pause/resume: `worklogger pause`, `worklogger resume`
- Stop: `worklogger stop`
- View logs: `worklogger log`
- Sync commits: `worklogger sync --new --desc "Fix bug"`
- Auto-log commits: `worklogger setup-hook`
- Export data: `worklogger export --csv --out data.csv`
- View stats: `worklogger summary`
- Web interface: `worklogger studio` (opens `http://localhost:8080`)
- Auth: `worklogger signup --github`, `worklogger login --local`, `worklogger logout`

## Development
- **Build**: `make build` (builds Vite frontend and Go binary).
- **Database**: Manage migrations with `make db/migrations/up` or `make db/migrations/reset`.
- **Dependencies**: Uses `cobra`, `bubbletea` (TUI), `httprouter` (server), and Vite (frontend).
