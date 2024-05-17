-- +goose Up
-- +goose StatementBegin
UPDATE links SET user_id=651963269 WHERE link_key='c046d12e';
UPDATE links SET user_id=651963269 WHERE link_key='ce25c5e2';
UPDATE links SET user_id=651963269 WHERE link_key='874288f0';
UPDATE links SET user_id=651963269 WHERE link_key='f4e0d7ad';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE links SET user_id=651963269 WHERE link_key=NULL;
UPDATE links SET user_id=651963269 WHERE link_key=NULL;
UPDATE links SET user_id=651963269 WHERE link_key=NULL;
UPDATE links SET user_id=651963269 WHERE link_key=NULL;
-- +goose StatementEnd
