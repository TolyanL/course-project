# Backend Agent Plan: «Электронное расписание»

## Цель

Реализовать полноценный backend на Golang/Fiber с PostgreSQL для системы электронного расписания.

---

## Фаза 1: Инфраструктура

### 1.1 Docker Compose

- [ ] PostgreSQL с volume для данных
- [ ] Backend сервис (Golang/Fiber)
- [ ] Общая сеть между сервисами
- [ ] Healthcheck для PostgreSQL
- [ ] Retry logic в Golang для ожидания готовности БД

---

## Фаза 2: Структура проекта

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
│   └── validators/    # проверка накладок
├── migrations/
│   └── 001_init.sql
└── go.mod/go.sum
```

### Задачи:

- [ ] Инициализировать Go модуль
- [ ] Создать структуру папок
- [ ] Подключить зависимости: Fiber, pgx, jwt, godotenv

---

## Фаза 3: База данных

### 3.1 migrations/001_init.sql

Создать таблицы:

- [ ] `teachers` (id, name, login, password_hash, role, created_at)
- [ ] `subjects` (id, name)
- [ ] `classrooms` (id, number)
- [ ] `groups` (id, name)
- [ ] `schedule` (id, group_id, date, pairs JSONB)
- [ ] Добавить seed для admin (login/password из ENV с fallback на dev)

### 3.2 SQL функции для валидации

- [ ] Функция проверки занятости аудитории
- [ ] Функция проверки занятости преподавателя

---

## Фаза 4: Константы (internal/constants/pairs.go)

- [ ] Определить времена пар с перерывами:
  - Pair1: 08:00 - 09:35
  - Pair2: 09:45 - 11:20
  - Pair3: 12:30 - 14:05
  - Pair4: 14:15 - 15:50
  - Pair5: 16:00 - 17:35
  - Pair6: 17:45 - 19:20
  - Pair7: 19:30 - 21:05

---

## Фаза 5: API эндпоинты

### 5.1 Auth

| Метод | URL             | Описание                      |
| ----- | --------------- | ----------------------------- |
| POST  | /api/auth/login | Вход (login + password → JWT) |

### 5.2 Schedule

| Метод  | URL                       | Роль          | Описание                            |
| ------ | ------------------------- | ------------- | ----------------------------------- |
| GET    | /api/schedule             | all           | Фильтры: date, group_id, teacher_id |
| POST   | /api/schedule/entries     | teacher/admin | Добавить пару                       |
| PUT    | /api/schedule/entries/:id | teacher/admin | Редактировать                       |
| DELETE | /api/schedule/entries/:id | teacher/admin | Удалить                             |

### 5.3 CRUD справочников

| Сущность   | GET             | POST            | PUT                 | DELETE              |
| ---------- | --------------- | --------------- | ------------------- | ------------------- |
| teachers   | /api/teachers   | /api/teachers   | /api/teachers/:id   | /api/teachers/:id   |
| subjects   | /api/subjects   | /api/subjects   | /api/subjects/:id   | /api/subjects/:id   |
| classrooms | /api/classrooms | /api/classrooms | /api/classrooms/:id | /api/classrooms/:id |
| groups     | /api/groups     | /api/groups     | /api/groups/:id     | /api/groups/:id     |

---

## Фаза 6: Валидация накладок (КРИТИЧНО!)

### При добавлении/изменении пары:

1. [ ] Проверить: аудитория свободна в этот день+пару (SQL запрос)
2. [ ] Проверить: преподаватель свободен в этот день+пару (SQL запрос)
3. [ ] Если конфликт → вернуть warning с деталями
4. [ ] Для teacher: позволить force save
5. [ ] Для admin: запретить при конфликте

**Важно:** Проверка через SQL WHERE, не через JSONB запросы в Golang!

---

## Фаза 7: Авторизация (JWT)

- [ ] Структура claims: {user_id, login, role}
- [ ] Access token с expires
- [ ] Middleware для проверки токена
- [ ] Защита эндпоинтов по ролям
- [ ] Password hashing (bcrypt)

---

## Фаза 8: Middleware

- [ ] CORS
- [ ] JSON body limit
- [ ] Request logging
- [ ] Auth middleware
- [ ] Role check middleware

---

## Проверка

- [ ] Написать unit тесты для validators
- [ ] Проверить все эндпоинты через curl/httpie
- [ ] Проверить валидацию накладок
- [ ] Проверить JWT авторизацию

---

## Критические моменты

| #   | Момент                                       | Решение                                            |
| --- | -------------------------------------------- | -------------------------------------------------- |
| 1   | Накладки по аудитории/преподавателю          | SQL CHECK при INSERT/UPDATE + валидация в handlers |
| 2   | Проверка конфликтов для teacher с force save | API возвращает warning, frontend сам решит         |
| 3   | Ожидание PostgreSQL                          | Healthcheck + retry в коде                         |
| 4   | Пароль админа                                | ENV файл с fallback для dev                        |

---

## Dependencies (go.mod)

```
github.com/gofiber/fiber/v2
github.com/jackc/pgx/v5
github.com/golang-jwt/jwt/v5
github.com/joho/godotenv
golang.org/x/crypto (bcrypt)
```

