-- +goose Up
CREATE TABLE images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    filename TEXT NOT NULL,

    mime_type TEXT NOT NULL,

    extension TEXT NOT NULL,

    size BIGINT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE images;