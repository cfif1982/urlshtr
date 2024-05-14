-- +goose Up
-- +goose StatementBegin
INSERT INTO links (link_key, link_url) VALUES
('c046d12e', 'https://practicum.yandex.ru'),
('ce25c5e2', 'https://testsite.ru'),
('874288f0', 'https://helloworld.ru/hi'),
('f4e0d7ad', 'https://onemoresite.com/qjkdyr')
ON CONFLICT (link_url) DO UPDATE SET 
link_key = EXCLUDED.link_key,
link_url = EXCLUDED.link_url;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM links WHERE link_key in (
  'c046d12e', 'ce25c5e2', '874288f0', 'f4e0d7ad'
);
-- +goose StatementEnd
