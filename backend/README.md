# Backend — Electronic Schedule (Electronnoe Raspisanie)

A Go REST API server for managing university class schedules. Built with Fiber and PostgreSQL, using raw SQL queries via pgx.

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.26 |
| HTTP Framework | Fiber v3 (Fasthttp) |
| Database | PostgreSQL 17/18 |
| DB Driver | pgx v5 (connection pool) |
| Migrations | Goose v3 (embedded, auto-run on startup) |
| Auth | JWT (golang-jwt v5, HS256) |
| Config | godotenv (.env files) |
| Linting | golangci-lint |
| Containerization | Docker + Docker Compose |

## Project Structure

```
backend/
├── cmd/
│   ├── main.go                          # Entry point: config, DB, migrations, routes, server
│   └── migrations/
│       ├── 001_init.sql                 # Schema creation (6 tables, 3 indexes)
│       └── 002_admin.sql                # Seed admin user (admin/admin123)
│
├── internal/
│   ├── config/
│   │   └── config.go                    # Env-based config loader
│   ├── constants/
│   │   └── pairs.go                     # 7 university class-pair time slots
│   ├── database/
│   │   └── database.go                  # PostgreSQL connection pool wrapper
│   ├── handlers/
│   │   └── handlers.go                  # HTTP request handlers (controllers)
│   ├── middleware/
│   │   └── auth.go                      # JWT auth + RBAC role middleware
│   ├── models/
│   │   └── models.go                    # Domain structs & DTOs
│   ├── repository/
│   │   └── repository.go                # Data access layer (raw SQL via pgx)
│   ├── service/
│   │   └── service.go                   # Business logic layer
│   └── validators/
│       ├── conflict.go                  # Schedule conflict detection
│       └── conflict_test.go             # Unit tests with mocks
│
├── tests/
│   └── api_test.sh                      # End-to-end shell-based API test
│
├── .env.example                         # Environment variable template
├── .golangci.yml                        # Linter configuration
├── Dockerfile                           # Multi-stage build (golang -> alpine)
├── docker-compose.yml                   # Production: postgres + backend
├── docker-compose.dev.yml               # Dev: postgres only
├── go.mod
└── go.sum
```

## Architecture

### Layered Architecture (3-Tier)

```
┌───────────────────────────────────────────────────┐
│              handlers (Controllers)                │
│   HTTP request/response, JSON binding,             │
│   status codes, input validation                   │
│   internal/handlers/handlers.go                    │
├───────────────────────────────────────────────────┤
│              service (Business Logic)              │
│   JWT creation/parsing, schedule assembly,         │
│   orchestration, delegates to repository           │
│   internal/service/service.go                      │
├───────────────────────────────────────────────────┤
│              repository (Data Access)              │
│   All SQL queries, pgx pool operations,            │
│   CRUD, conflict detection queries                 │
│   internal/repository/repository.go                │
├───────────────────────────────────────────────────┤
│              database (Infrastructure)             │
│   pgxpool.Pool wrapper, DSN management             │
│   internal/database/database.go                    │
└───────────────────────────────────────────────────┘
```

**Supporting packages:**

| Package | Purpose |
|---|---|
| `config` | Loads configuration from environment variables / `.env` |
| `middleware` | JWT authentication and role-based access control |
| `models` | Shared domain structs and request/response DTOs |
| `constants` | University class-pair time definitions (7 pairs/day) |
| `validators` | Business rule validation (schedule conflict detection) |

### Dependency Chain

```
main → handlers → service → repository → database
                                           ↑
                                    config (env vars)
```

Manual wiring in `main()` — no DI container.

## API Reference

### Authentication

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| POST | `/api/auth/login` | No | Login, returns JWT + role |

### Schedule

| Method | Endpoint | Auth | Role | Description |
|---|---|---|---|---|
| GET | `/api/schedule` | No | Any | Get schedule (filter by `group_id`, `teacher_id`, `date`, `start_date`, `end_date`) |
| POST | `/api/schedule/entries` | Yes | teacher, admin | Create schedule entry (with conflict detection) |
| PUT | `/api/schedule/entries/:id` | Yes | teacher, admin | Update schedule entry |
| DELETE | `/api/schedule/entries/:id` | Yes | teacher, admin | Delete schedule entry |

### Teachers

| Method | Endpoint | Auth | Role | Description |
|---|---|---|---|---|
| GET | `/api/teachers` | No | Any | List all teachers |
| POST | `/api/teachers` | Yes | admin | Create teacher |
| PUT | `/api/teachers/:id` | Yes | admin | Update teacher |
| DELETE | `/api/teachers/:id` | Yes | admin | Delete teacher (409 if FK references exist) |

### Subjects, Classrooms, Groups

Same CRUD pattern as Teachers. GET is public, POST/PUT/DELETE require admin role.

| Resource | Endpoints |
|---|---|
| `/api/subjects` | GET, POST, PUT /:id, DELETE /:id |
| `/api/classrooms` | GET, POST, PUT /:id, DELETE /:id |
| `/api/groups` | GET, POST, PUT /:id, DELETE /:id |

### Conflict Detection

Schedule entry creation/update runs conflict checks:

1. **Classroom conflict** — same classroom already booked for that pair on that date
2. **Teacher conflict** — same teacher already booked for that pair on that date

Response behavior:
- **Admin** without `force_save` → HTTP 409 with conflict details
- **Teacher** → saves anyway, includes warnings in response
- **Admin** with `force_save: true` → saves with warnings

## Database Schema

```
teachers        (id, name, login UNIQUE, password_hash, role CHECK admin/teacher, created_at)
subjects        (id, name)
classrooms      (id, number UNIQUE)
groups          (id, name UNIQUE)
schedules       (id, group_id FK→groups, date, UNIQUE(group_id, date))
schedule_entries(id, schedule_id FK→schedules, subject_id FK→subjects,
                 teacher_id FK→teachers, classroom_id FK→classrooms,
                 pair_number CHECK 1-7, UNIQUE(schedule_id, pair_number))
```

**Indexes:** `schedule_entries(schedule_id)`, `schedule_entries(teacher_id)`, `schedule_entries(classroom_id)`

**Referential actions:** `schedules.group_id` → CASCADE on delete; subject/teacher/classroom FKs → RESTRICT on delete.

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | `postgres` | Database password |
| `DB_NAME` | `schedule_db` | Database name |
| `JWT_SECRET` | — | Secret key for JWT signing |
| `SERVER_PORT` | `8080` | HTTP server port |

## Getting Started

### With Docker Compose (recommended)

```bash
# Production
docker-compose up -d

# Development (PostgreSQL only)
docker-compose -f docker-compose.dev.yml up -d
```

### Local Development

```bash
# 1. Copy and configure environment
cp .env.example .env

# 2. Ensure PostgreSQL is running and the database exists
createdb schedule_db

# 3. Run the server (migrations run automatically on startup)
go run ./cmd/main.go
```

### Running Tests

```bash
# Unit tests (conflict validator)
go test ./internal/validators/ -v

# End-to-end API test (requires running server + DB)
bash tests/api_test.sh
```

### Linting

```bash
golangci-lint run
```

## Available Make Commands

| Command | Description |
|---|---|
| `docker-compose up -d` | Start production stack |
| `docker-compose down` | Stop production stack |
| `go run ./cmd/main.go` | Run server locally |
| `go test ./...` | Run all Go tests |
| `go build -o backend ./cmd/main.go` | Build binary |

## Startup Sequence

1. Load config from `.env` via `config.Load()`
2. Connect to PostgreSQL with retry (30s, 1s interval)
3. Run embedded SQL migrations via goose
4. Wire dependency chain: `repository` → `service` → `handlers`
5. Create Fiber app with middleware stack (recover → logger → CORS)
6. Register all routes under `/api`
7. Start listening on configured port
8. Graceful shutdown on SIGINT/SIGTERM
