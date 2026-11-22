CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name VARCHAR(100) NOT NULL UNIQUE,

    parent_id UUID REFERENCES categories(id)
        ON DELETE SET NULL
        ON UPDATE CASCADE,

    CONSTRAINT categories_no_self_parent
        CHECK (id <> parent_id)
);

CREATE INDEX idx_categories_parent_id
    ON categories(parent_id);