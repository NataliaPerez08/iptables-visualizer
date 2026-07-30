# IPTables Visualizer

A web-based firewall policy management tool with support for **iptables** and **nftables** backends. Manage, visualize, validate, compile, and deploy Linux firewall rules through a clean web interface.

> Previously known as *Firewall Policy Manager*.

## Features

- **Graphical rule management** — Create and organize firewall policies with multiple rules
- **Dual driver support** — Compile rules to either `iptables` or `nftables` syntax
- **Validation** — Address, port, protocol, and action validation before deployment
- **Dry-run mode** — Preview the exact commands that will be executed without touching the firewall
- **One-click apply** — Deploy policies directly to the live firewall
- **Rollback** — Automatic rollback on failure to revert to the previous state
- **Audit logging** — Full audit trail of all user actions (login, CRUD, apply, etc.)
- **Role-based access control** — Admin, editor, and viewer roles
- **Single binary** — Frontend is embedded in the Go binary via `go:embed`

## Tech Stack

| Layer       | Technology                                                      |
|-------------|-----------------------------------------------------------------|
| Backend     | Go — [Chi router](https://github.com/go-chi/chi), JWT, bcrypt  |
| Database    | SQLite3 (WAL mode, single connection)                            |
| Frontend    | Vanilla JavaScript, HTML5, CSS3 (no framework)                   |
| API         | RESTful JSON over HTTP                                           |
| Firewall    | iptables / nftables (executed via `os/exec`)                     |

## Quick Start

### Prerequisites

- Go 1.21+
- Linux with `iptables` or `nftables` installed (optional for dry-run)
- `gcc` (required by `mattn/go-sqlite3` — CGO)

### Setup

```bash
# Clone and enter the project
git clone https://github.com/anomalyco/iptables-visualizer.git
cd iptables-visualizer

# Initialize Go module
go mod init github.com/anomalyco/iptables-visualizer
go mod tidy

# Build
cd backend
go build -o server .
```

### Configuration

All configuration is via environment variables:

| Variable                  | Default                   | Description                         |
|---------------------------|---------------------------|-------------------------------------|
| `SERVER_HOST`             | `0.0.0.0`                | Bind address                        |
| `SERVER_PORT`             | `8080`                    | HTTP port                           |
| `DB_PATH`                 | `./data/firewall.db`      | SQLite database path                |
| `JWT_SECRET`              | `change-me-in-production` | JWT signing secret                  |
| `JWT_EXPIRATION`          | `24h`                     | Token lifetime                      |
| `FIREWALL_DRIVER`         | `iptables`                | Default driver (`iptables`/`nftables`) |
| `DRY_RUN_ONLY`            | `true`                    | If true, no real firewall commands   |

### Run

```bash
# Dry-run only (safe — no firewall changes)
export DRY_RUN_ONLY=true
./backend/server

# Live mode (requires root/sudo for firewall access)
export DRY_RUN_ONLY=false
sudo ./backend/server
```

Open [http://localhost:8080](http://localhost:8080).

### Default Credentials

- **Username:** `admin`
- **Password:** `admin123`

> Change the password immediately in production.

## API

All endpoints are under `/api/v1`:

### Authentication
| Method | Path             | Description      |
|--------|------------------|------------------|
| POST   | `/auth/login`    | Login            |
| GET    | `/auth/me`       | Current user     |

### Users (admin only)
| Method | Path          | Description  |
|--------|---------------|--------------|
| POST   | `/users`      | Create user  |
| GET    | `/users`      | List users   |

### Policies
| Method | Path                    | Description            |
|--------|-------------------------|------------------------|
| GET    | `/policies`             | List policies          |
| POST   | `/policies`             | Create policy          |
| GET    | `/policies/{id}`        | Get policy             |
| PUT    | `/policies/{id}`        | Update policy          |
| DELETE | `/policies/{id}`        | Delete policy          |
| POST   | `/policies/{id}/validate` | Validate rules       |
| POST   | `/policies/{id}/dry-run`  | Preview commands      |
| POST   | `/policies/{id}/apply`    | Apply to firewall     |

### Audit (admin only)
| Method | Path       | Description   |
|--------|------------|---------------|
| GET    | `/audit`   | Query logs    |

## Project Structure

```
├── backend/
│   ├── main.go                    # Entry point, migrations, SPA handler
│   ├── migrations/                # SQL migration files
│   ├── internal/
│   │   ├── api/
│   │   │   ├── router.go          # Route definitions
│   │   │   ├── handlers/          # HTTP handlers
│   │   │   └── middleware/        # Auth & audit middleware
│   │   ├── config/                # Environment config
│   │   ├── deployment/            # Firewall executor & dry-run
│   │   ├── drivers/               # iptables & nftables compilers
│   │   ├── engine/                # Policy compiler & validator
│   │   ├── models/                # Data structures
│   │   └── repository/            # SQLite repositories
│   └── web/                       # Frontend SPA
│       ├── index.html
│       ├── css/
│       └── js/
│           ├── api.js             # HTTP client
│           ├── app.js             # SPA router
│           ├── components/        # Reusable UI components
│           └── views/             # Page views
└── .gitignore
```

## Development

The frontend is embedded in the Go binary at build time using `//go:embed web`. Edit the files under `backend/web/` and rebuild.

```bash
cd backend
go build -o server .
./server
```

## License

MIT
