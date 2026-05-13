# Финальный план: «Электронное расписание»

## 1. Инфраструктура

### 1.1 Docker Compose

- PostgreSQL (с volume для данных)
- Backend (Golang/Fiber)
- Frontend (Vite/Vue)
- Общая сеть

**Критический момент:** Backend должен ждать готовности PostgreSQL (healthcheck / retry logic)

---

## 2. Backend (Golang/Fiber)

### 2.1 Структура проекта

```
backend/
├── cmd/main.go
├── internal/
│   ├── config/        # загрузка .env
│   ├── constants/     # times.go (времена пар)
│   ├── database/      # подключение + миграции
│   ├── models/        # structs для всех сущностей
│   ├── handlers/      # HTTP handlers
│   ├── middleware/    # JWT auth
│   └── validators/   # проверка накладок
├── migrations/
│   └── 001_init.sql
└── go.mod/go.sum
```

### 2.2 База данных (migrations/001_init.sql)

```sql
-- Таблицы:
teachers (id, name, login, password_hash, role, created_at)
subjects (id, name)
classrooms (id, number)
groups (id, name)
schedule (id, group_id, date, pairs JSONB)

-- Seed:
admin: login/password из ENV (с fallback на dev значения)
```

### 2.3 Константы пар (internal/constants/pairs.go)

```go
// Времена пар с учётом перерывов
Pair1Start = "08:00", Pair1End = "09:35"
Pair2Start = "09:45", Pair2End = "11:20"
Pair3Start = "12:29", Pair3End = "14:05"  // после большого перерыва
...
```

### 2.4 API эндпоинты

| Метод                                         | URL                       | Роль          | Описание                            |
| --------------------------------------------- | ------------------------- | ------------- | ----------------------------------- |
| GET                                           | /api/schedule             | all           | Фильтры: date, group_id, teacher_id |
| POST                                          | /api/schedule/entries     | teacher/admin | Добавить пару                       |
| PUT                                           | /api/schedule/entries/:id | teacher/admin | Редактировать                       |
| DELETE                                        | /api/schedule/entries/:id | teacher/admin | Удалить                             |
| GET                                           | /api/teachers             | admin         | Список                              |
| POST                                          | /api/teachers             | admin         | Создать                             |
| PUT                                           | /api/teachers/:id         | admin         | Редактировать                       |
| DELETE                                        | /api/teachers/:id         | admin         | Удалить                             |
| (аналогично для subjects, classrooms, groups) |                           |               |                                     |
| POST                                          | /api/auth/login           | -             | Вход                                |

### 2.5 Валидация (критично!)

При добавлении/изменении пары:

```
1. Проверить: аудитория свободна в этот день+пару
2. Проверить: преподаватель свободен в этот день+пару
3. Если конфликт → вернуть warning + позволить force save (для teacher)
4. Для admin: запретить при конфликте
```

**Критический момент:** Проверка через SQL с WHERE на существующие записи (не через JSONB запросы в Golang — медленно)

### 2.6 Авторизация

- JWT tokens (access + expires)
- Middleware проверяет роль для защищённых эндпоинтов
- Структура claims: {user_id, login, role}

---

## 3. Frontend (Vue.js)

### 3.1 Структура

```
src/
├── App.vue
├── main.ts
├── router/index.ts
├── stores/
│   ├── auth.ts      # Pinia: user, token, login/logout
│   └── schedule.ts  # загрузка данных
├── components/
│   ├── ScheduleTable.vue
│   ├── PairEditor.vue
│   ├── FilterBar.vue
│   └── ...
├── views/
│   ├── ScheduleView.vue     # /
│   ├── LoginView.vue        # /login
│   ├── TeacherView.vue      # /teacher
│   └── AdminView.vue        # /admin/*
└── api/
    └── client.ts            # axios instance
```

### 3.2 Роутинг

- `/` — ScheduleView (public)
- `/login` — LoginView
- `/teacher` — TeacherView (auth required, role=teacher)
- `/admin` — AdminView (auth required, role=admin)
- `/admin/teachers`, `/admin/subjects`, etc.

### 3.3 ScheduleView

- Фильтры: дата (datepicker), группа (select)
- Таблица пар с номерами (pair 1, pair 2...)
- Приоритет сортировки: по added_at из JSONB
- Валидация на фронте: показывать warning если пара в прошлом

**Критический момент:** Как получить номер пары? Из позиции в JSONB массиве + сортировка по added_at.

### 3.4 TeacherView

- Свои занятия (group_id, date, subject, classroom)
- Добавление/редактирование с выбором группы, даты, предмета, аудитории
- Warning при конфликте (с кнопкой "всё равно добавить")

### 3.5 AdminView

- CRUD для справочников
- Управление всеми записями schedule

---

## 4. Критические моменты (сводка)

| #   | Момент                                       | Решение                                            |
| --- | -------------------------------------------- | -------------------------------------------------- |
| 1   | Накладки по аудитории/преподавателю          | SQL CHECK при INSERT/UPDATE                        |
| 2   | Проверка конфликтов для teacher с force save | API возвращает warning, frontend показывает кнопку |
| 3   | Получение номера пары                        | Позиция в отсортированном по added_at массиве      |
| 4   | JWT в Vue                                    | Хранение в localStorage, axios interceptor         |
| 5   | Ожидание PostgreSQL                          | Healthcheck в docker-compose или retry в Golang    |
| 6   | Пароль админа                                | ENV файл с fallback для dev                        |

