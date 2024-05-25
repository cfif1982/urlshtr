-- +goose Up
-- +goose StatementBegin
ALTER TABLE links ADD deleted_flag BOOLEAN;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE links DROP COLUMN deleted_flag;
-- +goose StatementEnd
