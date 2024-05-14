-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS links(
	link_key TEXT,
	link_url TEXT UNIQUE NOT NULL
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS links;
-- +goose StatementEnd
