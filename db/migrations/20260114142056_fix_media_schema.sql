-- +goose Up
-- +goose StatementBegin
ALTER TABLE media
DROP COLUMN id;

ALTER TABLE media
ADD COLUMN id UUID PRIMARY KEY DEFAULT gen_random_uuid();

ALTER TABLE media
DROP COLUMN slug;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
