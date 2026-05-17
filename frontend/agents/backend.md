# Backend API Documentation

## Base URL

```
http://localhost:8080/api
```

## Authentication

### POST /api/auth/login

Login and receive JWT token.
**Request:**

```json
{
  "login": "admin",
  "password": "admin123"
}
```

**Response (200):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "role": "admin"
}
```

**Response (401):**

```json
{
  "error": "invalid credentials"
}
```

---

## Schedule

### GET /api/schedule

Get schedule entries. Requires Bearer token.
**Query params:**

- `group_id` (optional) - filter by group
- `teacher_id` (optional) - filter by teacher
- `date` (optional) - filter by date (YYYY-MM-DD)
  **Response (200):**

```json
[
  {
    "id": 1,
    "group_id": 1,
    "date": "2024-01-15",
    "pairs": [
      {
        "id": 1,
        "schedule_id": 1,
        "subject_id": 1,
        "teacher_id": 1,
        "classroom_id": 1,
        "pair_number": 1,
        "subject": { "id": 1, "name": "Математика" },
        "teacher": { "id": 1, "name": "Иванов И.И." },
        "classroom": { "id": 1, "number": "101" }
      }
    ],
    "group": { "id": 1, "name": "ИУ8-21" }
  }
]
```

- ?date=2025-05-17 - конкретная дата
- ?start_date=2025-05-17&end_date=2025-05-23 - диапазон (неделя)

### POST /api/schedule/entries

Create schedule entry. Requires Bearer token. Role: teacher, admin.
**Request:**

```json
{
  "group_id": 1,
  "date": "2024-01-15",
  "subject_id": 1,
  "teacher_id": 1,
  "classroom_id": 1,
  "pair_number": 1,
  "force_save": false
}
```

**Response (201):**

```json
{
  "entry": {
    "id": 1,
    "schedule_id": 1,
    "subject_id": 1,
    "teacher_id": 1,
    "classroom_id": 1,
    "pair_number": 1
  }
}
```

**Response (409) - Conflict detected (admin without force_save):**

```json
{
  "error": "conflict detected",
  "details": [
    {
      "type": "classroom",
      "message": "Аудитория 101 уже занята в эту пару",
      "field": "classroom_id",
      "entry_id": 1
    },
    {
      "type": "teacher",
      "message": "Преподаватель Иванов уже занят в эту пару",
      "field": "teacher_id",
      "entry_id": 1
    }
  ]
}
```

### PUT /api/schedule/entries/:id

Update schedule entry. Requires Bearer token. Role: teacher, admin.
**Request:** Same as POST

### DELETE /api/schedule/entries/:id

Delete schedule entry. Requires Bearer token. Role: teacher, admin.
**Response (204):** No content

---

## Teachers

### GET /api/teachers

Get all teachers. Requires Bearer token.
**Response (200):**

```json
[
  {
    "id": 1,
    "name": "Administrator",
    "login": "admin",
    "role": "admin",
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

### POST /api/teachers

Create teacher. Requires Bearer token. Role: admin only.
**Request:**

```json
{
  "name": "Иванов И.И.",
  "login": "ivanov",
  "password": "password123",
  "role": "teacher"
}
```

**Response (201):**

```json
{
  "id": 2,
  "name": "Иванов И.И.",
  "login": "ivanov",
  "role": "teacher"
}
```

### PUT /api/teachers/:id

Update teacher. Requires Bearer token. Role: admin only.
**Request:**

```json
{
  "name": "Новое имя",
  "role": "teacher"
}
```

### DELETE /api/teachers/:id

Delete teacher. Requires Bearer token. Role: admin only.
**Response (204):** No content
**Response (409) - If teacher has schedule entries:**

```json
{
  "error": "cannot delete teacher: currently assigned to schedule entries"
}
```

---

## Subjects

### GET /api/subjects

Get all subjects. Requires Bearer token.
**Response (200):**

```json
[
  { "id": 1, "name": "Математика" },
  { "id": 2, "name": "Физика" }
]
```

### POST /api/subjects

Create subject. Requires Bearer token. Role: admin only.
**Request:**

```json
{ "name": "Математика" }
```

### PUT /api/subjects/:id

Update subject. Requires Bearer token. Role: admin only.

### DELETE /api/subjects/:id

Delete subject. Requires Bearer token. Role: admin only.
**Response (409) - If subject has schedule entries:**

```json
{
  "error": "cannot delete subject: currently assigned to schedule entries"
}
```

---

## Classrooms

### GET /api/classrooms

Get all classrooms. Requires Bearer token.
**Response (200):**

```json
[
  { "id": 1, "number": "101" },
  { "id": 2, "number": "202" }
]
```

### POST /api/classrooms

Create classroom. Requires Bearer token. Role: admin only.
**Request:**

```json
{ "number": "101" }
```

### PUT /api/classrooms/:id

Update classroom. Requires Bearer token. Role: admin only.

### DELETE /api/classrooms/:id

Delete classroom. Requires Bearer token. Role: admin only.
**Response (409) - If classroom has schedule entries:**

```json
{
  "error": "cannot delete classroom: currently assigned to schedule entries"
}
```

---

## Groups

### GET /api/groups

Get all groups. Requires Bearer token.
**Response (200):**

```json
[
  { "id": 1, "name": "ИУ8-21" },
  { "id": 2, "name": "ИУ8-22" }
]
```

### POST /api/groups

Create group. Requires Bearer token. Role: admin only.
**Request:**

```json
{ "name": "ИУ8-21" }
```

### PUT /api/groups/:id

Update group. Requires Bearer token. Role: admin only.

### DELETE /api/groups/:id

## Delete group. Requires Bearer token. Role: admin only.

## Authentication Header

All protected endpoints require:

```
Authorization: Bearer <jwt_token>
```

---

## Pair Numbers

Available pair numbers: 1-7
| Pair | Time |
|------|------|
| 1 | 08:00 - 09:35 |
| 2 | 09:45 - 11:20 |
| 3 | 12:20 - 13:55 |
| 4 | 14:05 - 15:40 |
| 5 | 16:00 - 17:35 |
| 6 | 17:45 - 19:20 |
| 7 | 19:30 - 21:05 |

---

## Default Admin Credentials

- Login: `admin`
- Password: `admin123`
