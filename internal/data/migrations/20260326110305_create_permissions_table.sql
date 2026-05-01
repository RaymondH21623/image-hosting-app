-- +goose Up
-- +goose StatementBegin
CREATE TABLE permissions (
    id bigserial PRIMARY KEY,
    code text NOT NULL
);

CREATE TABLE users_permissions (
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    permission_id bigint NOT NULL REFERENCES permissions ON DELETE CASCADE,
    PRIMARY KEY (user_id, permission_id)
);

INSERT INTO permissions (code) VALUES
    ('media:read'),
    ('media:write');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users_permissions;
DROP TABLE permissions;
-- +goose StatementEnd
