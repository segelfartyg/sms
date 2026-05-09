CREATE TABLE IF NOT EXISTS datasources (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    type        TEXT        NOT NULL,
    description TEXT        NOT NULL,
    content     TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
