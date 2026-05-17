-- +goose Up
-- +goose StatementBegin
INSERT INTO teachers (name, login, password_hash, role)
VALUES (
    'Administrator',
    'admin',
    'jLIjfQZ5yojbZGTqxg2pY0VROWQ=',
    'admin'
) ON CONFLICT (login) DO UPDATE SET password_hash = 'jLIjfQZ5yojbZGTqxg2pY0VROWQ=';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM teachers WHERE login = 'admin';
-- +goose StatementEnd
