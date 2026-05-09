CREATE TABLE IF NOT EXISTS datapoints (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    datasource_id UUID        NOT NULL REFERENCES datasources(id) ON DELETE CASCADE,
    tag           TEXT        NOT NULL,
    content       TEXT        NOT NULL DEFAULT '',
    description   TEXT        NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS datapoints_datasource_id_idx ON datapoints(datasource_id);
