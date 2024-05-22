-- +goose Up
-- +goose StatementBegin
ALTER TABLE links ADD user_id INT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE links DROP COLUMN user_id;
-- +goose StatementEnd
