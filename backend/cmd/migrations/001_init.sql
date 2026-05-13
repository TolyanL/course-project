-- +goose Up
-- +goose StatementBegin
CREATE TABLE teachers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    login VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'teacher' CHECK (role IN ('admin', 'teacher')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE subjects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE classrooms (
    id SERIAL PRIMARY KEY,
    number VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE schedules (
    id SERIAL PRIMARY KEY,
    group_id INTEGER REFERENCES groups(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    UNIQUE(group_id, date)
);

CREATE TABLE schedule_entries (
    id SERIAL PRIMARY KEY,
    schedule_id INTEGER REFERENCES schedules(id) ON DELETE CASCADE,
    subject_id INTEGER REFERENCES subjects(id) ON DELETE RESTRICT,
    teacher_id INTEGER REFERENCES teachers(id) ON DELETE RESTRICT,
    classroom_id INTEGER REFERENCES classrooms(id) ON DELETE RESTRICT,
    pair_number INTEGER NOT NULL CHECK (pair_number BETWEEN 1 AND 7),
    UNIQUE(schedule_id, pair_number)
);

CREATE INDEX idx_schedule_entries_schedule ON schedule_entries(schedule_id);
CREATE INDEX idx_schedule_entries_teacher ON schedule_entries(teacher_id);
CREATE INDEX idx_schedule_entries_classroom ON schedule_entries(classroom_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS schedule_entries CASCADE;
DROP TABLE IF EXISTS schedules CASCADE;
DROP TABLE IF EXISTS groups CASCADE;
DROP TABLE IF EXISTS classrooms CASCADE;
DROP TABLE IF EXISTS subjects CASCADE;
DROP TABLE IF EXISTS teachers CASCADE;
-- +goose StatementEnd
