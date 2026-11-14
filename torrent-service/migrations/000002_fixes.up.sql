UPDATE torrents SET description = '' WHERE description IS NULL;
ALTER TABLE torrents ALTER COLUMN description SET NOT NULL;
ALTER TABLE torrents ALTER COLUMN description SET DEFAULT '';

ALTER TABLE torrents
    ALTER COLUMN created_at TYPE TIMESTAMPTZ,
    ALTER COLUMN updated_at TYPE TIMESTAMPTZ;