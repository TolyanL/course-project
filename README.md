# Electronic Schedule (Electronnoe Raspisanie)

A full-stack university schedule management system. Admins manage teachers, subjects, classrooms, groups, and the weekly timetable. Teachers view and edit their own schedule entries with conflict detection. The public can browse any group's schedule without authentication.

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend | Vue 3, TypeScript, Vite 8, Pinia, PrimeVue 4 |
| Backend | Go 1.26, Fiber v3, pgx v5, JWT |
| Database | PostgreSQL 17/18 |
| Infrastructure | Docker, Docker Compose |

## Project Structure

```
course-project/
├── frontend/              # Vue 3 SPA (details → frontend/README.md)
│   ├── src/
│   │   ├── api/           # Axios HTTP client + API modules
│   │   ├── components/    # Reusable UI (ScheduleTable, EntityTable, PairEditor, etc.)
│   │   ├── router/        # Vue Router with role-based guards
│   │   ├── stores/        # Pinia stores (auth, schedule)
│   │   ├── types/         # TypeScript interfaces & constants
│   │   └── views/         # Page components (Login, Schedule, Teacher, Admin/*)
│   └── ...
│
├── backend/               # Go REST API (details → backend/README.md)
│   ├── cmd/
│   │   ├── main.go        # Entry point
│   │   └── migrations/    # SQL schema (goose, embedded)
│   ├── internal/
│   │   ├── handlers/      # HTTP handlers (controllers)
│   │   ├── service/       # Business logic (JWT, orchestration)
│   │   ├── repository/    # Data access (raw SQL via pgx)
│   │   ├── database/      # PostgreSQL connection pool
│   │   ├── models/        # Domain structs & DTOs
│   │   ├── middleware/    # JWT auth + RBAC
│   │   ├── validators/    # Conflict detection
│   │   ├── constants/     # Class-pair time slots
│   │   └── config/        # Environment config loader
│   └── ...
│
├── .gitignore
├── .gitattributes
└── package.json           # Root (unused, for monorepo tooling)
```

## Architecture

```
┌─────────────────────┐         ┌─────────────────────┐
│      Frontend       │  HTTP   │      Backend        │
│                     │ ◄─────► │                     │
│  Vue 3 SPA          │  JSON   │  Fiber v3 REST API  │
│  Pinia / Router     │  REST   │  JWT Auth + RBAC    │
│  PrimeVue / Axios   │         │  pgx (raw SQL)      │
└─────────────────────┘         └────────┬────────────┘
                                         │
                                 ┌───────▼────────┐
                                 │   PostgreSQL    │
                                 │  6 tables       │
                                 │  3 indexes      │
                                 └────────────────┘
```

## Prerequisites

- **Node.js** ^20.19 || >=22.12
- **Go** 1.26+
- **PostgreSQL** 17+ (or Docker)
- **Docker & Docker Compose** (optional)

## Quick Start

### Option 1: Docker Compose

```bash
cd backend

# Copy and configure environment
cp .env.example .env

# Start PostgreSQL + backend
docker-compose up -d

# In another terminal — start frontend
cd ../frontend
npm install
npm run dev
```

### Option 2: Local Development

```bash
# 1. Start PostgreSQL and create database
createdb schedule_db

# 2. Configure backend
cd backend
cp .env.example .env
# Edit .env with your DB credentials and JWT_SECRET

# 3. Start backend (migrations run automatically)
go run ./cmd/main.go

# 4. Start frontend (in a separate terminal)
cd frontend
npm install
npm run dev
```

### Default Credentials

| Login | Password | Role |
|---|---|---|
| `admin` | `admin123` | admin |

## API Overview

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| POST | `/api/auth/login` | No | Login → JWT + role |
| GET | `/api/schedule` | No | Browse schedule (filterable) |
| GET | `/api/teachers` | No | List teachers |
| GET | `/api/subjects` | No | List subjects |
| GET | `/api/classrooms` | No | List classrooms |
| GET | `/api/groups` | No | List groups |
| POST/PUT/DELETE | `/api/*` | Admin | CRUD for teachers, subjects, classrooms, groups |
| POST/PUT/DELETE | `/api/schedule/entries` | Teacher/Admin | Schedule entry management with conflict detection |

Full API reference → [backend/README.md](backend/README.md#api-reference)

## User Roles

| Role | Capabilities |
|---|---|
| **Admin** | Full CRUD on all entities, force-save schedule conflicts |
| **Teacher** | View own schedule, create/edit own entries, override own conflicts |

## Key Features

- **Weekly schedule grid** — 7 days × 7 class periods with responsive layout
- **Conflict detection** — prevents double-booking of classrooms and teachers
- **Role-based access** — separate admin and teacher interfaces
- **JWT authentication** — token-based with auto-logout on expiry
- **Public schedule** — anyone can browse the timetable without logging in

## Available Scripts

### Frontend (`frontend/`)

| Script | Description |
|---|---|
| `npm run dev` | Start Vite dev server |
| `npm run build` | Type-check + production build |
| `npm run lint` | Run OxLint + ESLint |
| `npm run format` | Format with Prettier |

### Backend (`backend/`)

| Command | Description |
|---|---|
| `go run ./cmd/main.go` | Run server locally |
| `go test ./...` | Run all tests |
| `golangci-lint run` | Lint Go code |
| `go build -o backend ./cmd/main.go` | Build binary |

## License

This project is for educational purposes.
