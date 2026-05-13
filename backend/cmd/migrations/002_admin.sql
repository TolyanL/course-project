-- +goose Up
-- +goose StatementBegin
INSERT INTO teachers (name, login, password_hash, role)
VALUES (
    'Administrator',
    'admin',
    '$2a$10$kKish/UoNreWtEqcMrPfvOW7EmKOJJ/uFanpVc3vE/RSQFJqQXLKu',
    'admin'
) ON CONFLICT (login) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM teachers WHERE login = 'admin';
-- +goose StatementEnd
