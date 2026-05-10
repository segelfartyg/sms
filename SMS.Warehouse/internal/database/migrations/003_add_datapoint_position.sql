ALTER TABLE datapoints ADD COLUMN IF NOT EXISTS position INTEGER NOT NULL DEFAULT 0;

UPDATE datapoints d
SET position = sub.rn - 1
FROM (
    SELECT id, ROW_NUMBER() OVER (PARTITION BY datasource_id ORDER BY created_at) - 1 AS rn
    FROM datapoints
) sub
WHERE d.id = sub.id;
