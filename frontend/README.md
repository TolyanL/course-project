# Frontend — Electronic Schedule (Electronnoe Raspisanie)

A Vue 3 single-page application for managing and viewing university class schedules. Supports two user roles: **Admin** and **Teacher**.

## Tech Stack

| Layer            | Technology                                |
| ---------------- | ----------------------------------------- |
| Framework        | Vue 3 (Composition API, `<script setup>`) |
| Language         | TypeScript                                |
| Build Tool       | Vite 8                                    |
| State Management | Pinia 3                                   |
| Routing          | Vue Router 5                              |
| UI Library       | PrimeVue 4 (Aura theme) + PrimeIcons      |
| CSS Utilities    | Bootstrap 5 (CSS only)                    |
| HTTP Client      | Axios                                     |
| Linting          | ESLint + OxLint + Prettier                |
| Type Checking    | vue-tsc                                   |

## Project Structure

```
frontend/
├── index.html                      # SPA entry point
├── vite.config.ts                  # Vite config with @ alias
├── package.json
├── tsconfig.json
│
└── src/
    ├── main.ts                     # App bootstrap: Pinia, Router, PrimeVue
    ├── App.vue                     # Root shell: header (auth), router-view
    │
    ├── api/
    │   └── client.ts               # Axios instance + all API modules
    │
    ├── router/
    │   └── index.ts                # Route definitions + navigation guards
    │
    ├── stores/
    │   ├── auth.ts                 # Authentication state (token, role, user)
    │   └── schedule.ts             # Schedule entries CRUD state
    │
    ├── types/
    │   └── index.ts                # All TypeScript interfaces & constants
    │
    ├── views/
    │   ├── LoginView.vue           # /login — login form
    │   ├── ScheduleView.vue        # / — public weekly schedule
    │   ├── TeacherView.vue         # /teacher — teacher panel
    │   ├── AdminView.vue           # /admin — admin layout (sidebar/tabs)
    │   ├── AdminScheduleView.vue   # /admin/schedule — manage schedule
    │   ├── AdminTeachersView.vue   # /admin/teachers — CRUD teachers
    │   ├── AdminSubjectsView.vue   # /admin/subjects — CRUD subjects
    │   ├── AdminClassroomsView.vue # /admin/classrooms — CRUD classrooms
    │   └── AdminGroupsView.vue     # /admin/groups — CRUD groups
    │
    └── components/
        ├── ScheduleTable.vue       # Weekly schedule grid renderer
        ├── FilterBar.vue           # Week & group filter dropdowns
        ├── PairEditor.vue          # Create/edit class pair dialog
        ├── ConflictModal.vue       # Conflict resolution modal
        ├── EntityTable.vue         # Generic CRUD data table
        └── EntityForm.vue          # Generic CRUD form dialog
```

## Architecture

### Routing & Access Control

Routes use `meta.requiresAuth`, `meta.guest`, and `meta.role` fields. A global `beforeEach` guard enforces:

| Path                | Component             | Auth       | Role      |
| ------------------- | --------------------- | ---------- | --------- |
| `/`                 | `ScheduleView`        | No         | Any       |
| `/login`            | `LoginView`           | Guest-only | —         |
| `/teacher`          | `TeacherView`         | Yes        | `teacher` |
| `/admin`            | `AdminView`           | Yes        | `admin`   |
| `/admin/schedule`   | `AdminScheduleView`   | Yes        | `admin`   |
| `/admin/teachers`   | `AdminTeachersView`   | Yes        | `admin`   |
| `/admin/subjects`   | `AdminSubjectsView`   | Yes        | `admin`   |
| `/admin/classrooms` | `AdminClassroomsView` | Yes        | `admin`   |
| `/admin/groups`     | `AdminGroupsView`     | Yes        | `admin`   |

All admin child routes are **lazy-loaded** via dynamic `import()`.

### Authentication

- JWT-based with `localStorage` persistence
- Token decoded client-side to extract `user_id` and `login`
- Axios interceptor attaches `Authorization: Bearer <token>` to every request
- 401 responses trigger automatic logout and redirect to `/login`

### State Management (Pinia)

| Store      | Purpose                                                                                                               |
| ---------- | --------------------------------------------------------------------------------------------------------------------- |
| `auth`     | Token, role, teacherId, name. Actions: `login()`, `logout()`, `initFromStorage()`                                     |
| `schedule` | Schedule entries with loading/error state. Actions: `fetchSchedule()`, `addEntry()`, `updateEntry()`, `deleteEntry()` |

Reference data (groups, subjects, classrooms, teachers) is fetched directly by views, not stored in Pinia.

### API Layer (`api/client.ts`)

Single Axios instance with base URL `http://localhost:8080/api`. API modules organized by domain:

| Module          | Endpoints                                                                                              |
| --------------- | ------------------------------------------------------------------------------------------------------ |
| `authApi`       | `POST /auth/login`                                                                                     |
| `scheduleApi`   | `GET /schedule`, `POST /schedule/entries`, `PUT /schedule/entries/:id`, `DELETE /schedule/entries/:id` |
| `teachersApi`   | `GET /teachers`, `POST /teachers`, `PUT /teachers/:id`, `DELETE /teachers/:id`                         |
| `subjectsApi`   | `GET /subjects`, `POST /subjects`, `PUT /subjects/:id`, `DELETE /subjects/:id`                         |
| `classroomsApi` | `GET /classrooms`, `POST /classrooms`, `PUT /classrooms/:id`, `DELETE /classrooms/:id`                 |
| `groupsApi`     | `GET /groups`, `POST /groups`, `PUT /groups/:id`, `DELETE /groups/:id`                                 |

Conflict responses (HTTP 409) return structured `ConflictDetail[]` with conflict type, message, and affected fields.

### Component Hierarchy

```
App.vue
├── Menubar (header, shown when authenticated)
├── LoginView
├── ScheduleView
│   ├── FilterBar
│   └── ScheduleTable
├── TeacherView
│   ├── ScheduleTable
│   └── PairEditor
│       └── ConflictModal
└── AdminView (layout shell with sidebar/tabs)
    ├── AdminScheduleView
    │   ├── ScheduleTable
    │   └── PairEditor
    │       └── ConflictModal
    ├── AdminTeachersView
    │   ├── EntityTable
    │   └── EntityForm
    ├── AdminSubjectsView
    │   ├── EntityTable
    │   └── EntityForm
    ├── AdminClassroomsView
    │   ├── EntityTable
    │   └── EntityForm
    └── AdminGroupsView
        ├── EntityTable
        └── EntityForm
```

### Key Components

- **`ScheduleTable`** — Renders a weekly grid (7 days x 7 pair slots). Supports multi-week display, today highlighting, and conditional action buttons based on role and ownership.
- **`PairEditor`** — Smart form that auto-detects available pair slots for a given teacher+date. Supports `force_save` for teachers to override conflicts.
- **`EntityTable` / `EntityForm`** — Generic CRUD components reused across all 4 admin entity views (teachers, subjects, classrooms, groups).

## Getting Started

```bash
npm install
npm run dev
```

## Available Scripts

| Script               | Description                       |
| -------------------- | --------------------------------- |
| `npm run dev`        | Start Vite dev server             |
| `npm run build`      | Type-check + production build     |
| `npm run preview`    | Preview production build          |
| `npm run type-check` | Run `vue-tsc` type checking       |
| `npm run lint`       | Run OxLint + ESLint with auto-fix |
| `npm run format`     | Format with Prettier              |
