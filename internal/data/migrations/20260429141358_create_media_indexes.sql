-- +goose Up
-- +goose StatementBegin
CREATE INDEX media_title_idx ON media (title);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX media_title_idx;
-- +goose StatementEnd
