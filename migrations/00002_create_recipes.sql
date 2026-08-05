-- +goose Up
CREATE TABLE recipes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    author_id UUID NOT NULL,

    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,

    ingredients TEXT NOT NULL,
    instructions TEXT NOT NULL,

    prep_time_min INTEGER NOT NULL
        CHECK (prep_time_min >= 0),

    difficulty TEXT NOT NULL
        CHECK (difficulty IN ('easy', 'medium', 'hard'))

    servings INTEGER NOT NULL
        CHECK (servings > 0),

    rating_avg NUMERIC(3,1) NOT NULL DEFAULT 0
        CHECK (rating_avg BETWEEN 0 AND 5),

    rating_count INTEGER NOT NULL DEFAULT 0
        CHECK (rating_count >= 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_recipes_author
        FOREIGN KEY (author_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_recipes_author_id
    ON recipes(author_id);

CREATE INDEX idx_recipes_title
    ON recipes(title);

-- +goose Down
DROP TABLE recipes;