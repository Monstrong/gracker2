-- Убираем глобальную уникальность названий категорий
ALTER TABLE categories
    DROP CONSTRAINT categories_name_key;

-- Уникальные названия среди корневых категорий
CREATE UNIQUE INDEX idx_categories_root_name
    ON categories (name)
    WHERE parent_id IS NULL;

-- Уникальные названия среди дочерних категорий одного родителя
CREATE UNIQUE INDEX idx_categories_parent_name
    ON categories (parent_id, name)
    WHERE parent_id IS NOT NULL;

-- Внешний ключ от торрентов к категориям
ALTER TABLE torrents
    ADD CONSTRAINT fk_torrents_category
    FOREIGN KEY (category_id)
    REFERENCES categories (id)
    ON DELETE RESTRICT
    ON UPDATE CASCADE;

CREATE INDEX idx_torrents_category_id
    ON torrents (category_id);