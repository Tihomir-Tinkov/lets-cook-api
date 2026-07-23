-- +goose Up
CREATE TABLE recipe_images (
    recipe_id UUID NOT NULL,
    image_id UUID NOT NULL,

    display_order INTEGER NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_recipe_images_recipe
        FOREIGN KEY (recipe_id)
        REFERENCES recipes(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_recipe_images_image
        FOREIGN KEY (image_id)
        REFERENCES images(id)
        ON DELETE CASCADE,

    PRIMARY KEY (recipe_id, image_id),

    UNIQUE (recipe_id, display_order)
);

CREATE INDEX idx_recipe_images_recipe_id
    ON recipe_images(recipe_id);

-- +goose Down
DROP TABLE recipe_images;