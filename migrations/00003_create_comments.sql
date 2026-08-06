-- +goose Up
CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    recipe_id UUID NOT NULL,
    author_id UUID NOT NULL,

    content TEXT NOT NULL,

    rating INTEGER NOT NULL
        CHECK (rating BETWEEN 1 AND 5),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_comments_recipe
        FOREIGN KEY (recipe_id)
        REFERENCES recipes(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_comments_author
        FOREIGN KEY (author_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT unique_recipe_author_comment
        UNIQUE (recipe_id, author_id)
);

CREATE INDEX idx_comments_recipe_id
    ON comments(recipe_id);

CREATE INDEX idx_comments_author_id
    ON comments(author_id);

-- +goose Down
DROP TABLE comments;