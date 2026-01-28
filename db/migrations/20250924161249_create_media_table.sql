-- +goose Up
-- +goose StatementBegin
CREATE TABLE media (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    public_media_id TEXT NOT NULL UNIQUE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    filename TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    size BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE media;
-- +goose StatementEnd
