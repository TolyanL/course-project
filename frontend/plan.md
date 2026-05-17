# Frontend Agent Plan: «Электронное расписание»

## Цель
Реализовать frontend на Vue 3 / TypeScript / Vite для системы электронного расписания.

---

## Технологии

- Vue 3 + Composition API (`<script setup>`)
- TypeScript
- Vite
- Pinia (state management)
- Vue Router
- Axios (HTTP client)

---

## Фаза 1: Структура проекта

```
src/
├── App.vue
├── main.ts
├── router/index.ts
├── stores/
│   ├── auth.ts        # role, token, login/logout
│   └── schedule.ts    # загрузка расписания
├── components/
│   ├── ScheduleTable.vue    # Таблица пар
│   ├── PairEditor.vue       # Форма добавления/редактирования пары
│   ├── FilterBar.vue         # Фильтры (дата, группа)
│   ├── ConflictModal.vue    # Модалка конфликта
│   └── EntityTable.vue      # Универсальная таблица для админки
├── views/
│   ├── ScheduleView.vue      # / - публичное расписание
│   ├── LoginView.vue         # /login
│   ├── TeacherView.vue       # /teacher
│   └── AdminView.vue         # /admin/*
├── api/
│   └── client.ts             # axios instance с interceptors
├── types/
│   └── index.ts              # TypeScript interfaces
└── composables/
    └── useApi.ts             # хуки для работы с API
```

---

## Фаза 2: Типы (types/index.ts)

```typescript
// Auth
interface AuthResponse {
  token: string;
  role: 'admin' | 'teacher';
}

// Справочники
interface Teacher {
  id: number;
  name: string;
  login: string;
  role: 'admin' | 'teacher';
  created_at: string;
}

interface Subject {
  id: number;
  name: string;
}

interface Classroom {
  id: number;
  number: string;
}

interface Group {
  id: number;
  name: string;
}

// Schedule
interface Pair {
  id: number;
  schedule_id: number;
  subject_id: number;
  teacher_id: number;
  classroom_id: number;
  pair_number: number; // 1-7
  subject: Subject;
  teacher: Teacher;
  classroom: Classroom;
}

interface ScheduleEntry {
  id: number;
  group_id: number;
  date: string; // YYYY-MM-DD
  pairs: Pair[];
  group: Group;
}

// API responses
interface ConflictDetail {
  type: 'classroom' | 'teacher';
  message: string;
  field: string;
  entry_id: number;
}

interface ConflictResponse {
  error: 'conflict detected';
  details: ConflictDetail[];
}

interface ApiError {
  error: string;
}
```

---

## Фаза 3: API клиент (api/client.ts)

Base URL: `http://localhost:8080/api`

- [ ] Базовый axios instance
- [ ] Interceptor для добавления `Authorization: Bearer <token>`
- [ ] Interceptor для обработки ошибок (401 → logout, etc.)
- [ ] Типизированные методы для каждого ресурса

```typescript
// Эндпоинты согласно backend API:
POST   /auth/login         → AuthResponse
GET    /schedule           → ScheduleEntry[] (filters: group_id, teacher_id, date)
POST   /schedule/entries   → { entry: Pair }
PUT    /schedule/entries/:id
DELETE /schedule/entries/:id
GET    /teachers           → Teacher[]
POST   /teachers           → Teacher
PUT    /teachers/:id
DELETE /teachers/:id
GET    /subjects           → Subject[]
POST   /subjects           → Subject
PUT    /subjects/:id
DELETE /subjects/:id
GET    /classrooms         → Classroom[]
POST   /classrooms          → Classroom
PUT    /classrooms/:id
DELETE /classrooms/:id
GET    /groups             → Group[]
POST   /groups              → Group
PUT    /groups/:id
DELETE /groups/:id
```

---

## Фаза 4: Stores

### 4.1 auth.ts (Pinia)
- [ ] State: `role: 'admin' | 'teacher' | null`, `token: string | null`
- [ ] Действие: `login(login, password)` → сохраняет token и role
- [ ] Действие: `logout()` → очищает state
- [ ] Persist: token в localStorage (восстановление при загрузке)
- [ ] Getter: `isAdmin`, `isTeacher`, `isAuthenticated`

### 4.2 schedule.ts (Pinia)
- [ ] State: `entries: ScheduleEntry[]`, `loading`, `error`
- [ ] Действие: `fetchSchedule(filters?)` → GET /schedule
- [ ] Действие: `addEntry(data, force?)` → POST /schedule/entries
- [ ] Действие: `updateEntry(id, data, force?)` → PUT /schedule/entries/:id
- [ ] Действие: `deleteEntry(id)` → DELETE /schedule/entries/:id

---

## Фаза 5: Роутинг (router/index.ts)

| Путь | Компонент | Доступ | Описание |
|------|-----------|--------|----------|
| `/` | ScheduleView | public | Публичное расписание |
| `/login` | LoginView | guest | Вход |
| `/teacher` | TeacherView | teacher/admin | Панель преподавателя |
| `/admin` | AdminView | admin | Главная админки |
| `/admin/teachers` | AdminTeachersView | admin | Управление преподавателями |
| `/admin/subjects` | AdminSubjectsView | admin | Управление предметами |
| `/admin/classrooms` | AdminClassroomsView | admin | Управление аудиториями |
| `/admin/groups` | AdminGroupsView | admin | Управление группами |

- [ ] Navigation guards:
  - Проверка авторизации для protected routes
  - Проверка роли для admin routes
- [ ] Redirect на /login для неавторизованных
- [ ] Redirect на / для unauthorized access

---

## Фаза 6: Компоненты

### 6.1 FilterBar.vue
- [ ] Datepicker для выбора даты
- [ ] Select для группы
- [ ] Кнопка "Показать"
- [ ] Emits: `filter-change { date, group_id }`

### 6.2 ScheduleTable.vue
- [ ] Отображение ScheduleEntry
- [ ] Пары отсортированы по pair_number (1-7)
- [ ] Колонки: номер пары, время, предмет, аудитория, преподаватель
- [ ] Время пары по константе (pair_number → время)
- [ ] Пустое состояние если нет пар на этот день

**Константы пар:**
| № | Время |
|---|-------|
| 1 | 08:00 - 09:35 |
| 2 | 09:45 - 11:20 |
| 3 | 12:20 - 13:55 |
| 4 | 14:05 - 15:40 |
| 5 | 16:00 - 17:35 |
| 6 | 17:45 - 19:20 |
| 7 | 19:30 - 21:05 |

### 6.3 PairEditor.vue
- [ ] Поля: группа (select), дата (datepicker), предмет (select), аудитория (select), номер пары (select 1-7)
- [ ] Валидация: все поля обязательны
- [ ] Отображение ConflictModal при 409 ответе
- [ ] Emits: `save(data)`, `cancel`
- [ ] Props: `isEdit: boolean`, `entry?: Pair`

### 6.4 ConflictModal.vue
- [ ] Показ списка конфликтов (details[])
- [ ] Типы конфликтов: classroom, teacher
- [ ] Кнопка "Отмена" — закрыть модалку
- [ ] Кнопка "Сохранить принудительно" — только для teacher
- [ ] Emits: `force-save`, `cancel`

### 6.5 EntityTable.vue (универсальный для админки)
- [ ] Generic компонент для таблиц
- [ ] Props: `columns`, `data`, `onEdit`, `onDelete`
- [ ] Кнопки: добавить, редактировать, удалить

### 6.6 EntityForm.vue (универсальная форма)
- [ ] Props: `fields`, `initialData`
- [ ] emits: `submit`, `cancel`

---

## Фаза 7: Views

### 7.1 LoginView.vue
- [ ] Форма: login, password
- [ ] Обработка ошибки 401
- [ ] После успеха: redirect на / или /teacher или /admin в зависимости от роли

### 7.2 ScheduleView.vue
- [ ] FilterBar вверху
- [ ] ScheduleTable с результатами
- [ ] По умолчанию: сегодняшняя дата

### 7.3 TeacherView.vue
- [ ] ScheduleTable (только свои пары — filter by teacher_id)
- [ ] Кнопка "Добавить пару"
- [ ] PairEditor в модалке
- [ ] Кнопки редактирования/удаления своих пар
- [ ] Получить teacher_id из auth store

### 7.4 AdminView.vue (layout)
- [ ] Sidebar с навигацией (/admin/teachers, /admin/subjects, etc.)
- [ ] RouterView для дочерних страниц
- [ ] Кнопка "Выйти"

### 7.5 AdminTeachersView.vue
- [ ] EntityTable со списком преподавателей
- [ ] Кнопка "Добавить"
- [ ] Модалка с формой: имя, логин, пароль (при создании), роль
- [ ] Обработка ошибки 409 (нельзя удалить если есть записи)

### 7.6 AdminSubjectsView.vue
- [ ] EntityTable со списком предметов
- [ ] Кнопка "Добавить"
- [ ] Форма: название
- [ ] Обработка ошибки 409

### 7.7 AdminClassroomsView.vue
- [ ] EntityTable со списком аудиторий
- [ ] Форма: номер аудитории

### 7.8 AdminGroupsView.vue
- [ ] EntityTable со списком групп
- [ ] Форма: название группы

---

## Фаза 8: App.vue

- [ ] RouterView
- [ ] Header с логином и кнопкой выхода (если авторизован)
- [ ] Toast/уведомления для ошибок

---

## Фаза 9: Интеграция

- [ ] Подключить axios к backend (CORS уже настроен на бэке)
- [ ] Проверить flow: login → получить token → использовать в запросах
- [ ] Проверить force save для teacher
- [ ] Проверить защиту admin routes

---

## Проверка

- [ ] Typecheck: `npm run type-check`
- [ ] Lint: `npm run lint`
- [ ] Build: `npm run build`
- [ ] Manual тест: login, просмотр расписания, добавление пары

---

## Критические моменты

| # | Момент | Решение |
|---|--------|---------|
| 1 | JWT хранение | localStorage + axios interceptor |
| 2 | Защита роутов | Navigation guards (auth + role check) |
| 3 | Номер пары | Уже есть в pair_number (1-7), время из константы |
| 4 | Force save | Кнопка в ConflictModal, параметр `force_save: true` |
| 5 | 409 Conflict | Показывать детали конфликта, кнопка принудительного сохранения (только для teacher) |

---

## Dependencies

```json
{
  "dependencies": {
    "pinia": "^3.0.4",
    "vue": "^3.5.32",
    "vue-router": "^5.0.4",
    "axios": "^1.7.x"
  }
}
```

**Установить:** `npm install axios`