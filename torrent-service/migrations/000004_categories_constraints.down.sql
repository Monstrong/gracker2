DROP INDEX IF EXISTS idx_torrents_category_id;

ALTER TABLE torrents
    DROP CONSTRAINT IF EXISTS fk_torrents_category;

DROP INDEX IF EXISTS idx_categories_parent_name;
DROP INDEX IF EXISTS idx_categories_root_name;

ALTER TABLE categories
    ADD CONSTRAINT categories_name_key UNIQUE (name);