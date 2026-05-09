CREATE TABLE IF NOT EXISTS pages (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    title      TEXT        NOT NULL,
    slug       TEXT        NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS boxes (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    page_id    UUID        NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    type       TEXT        NOT NULL,
    content    JSONB       NOT NULL DEFAULT '{}',
    position   INTEGER     NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS boxes_page_id_idx       ON boxes(page_id);
CREATE INDEX IF NOT EXISTS boxes_page_position_idx ON boxes(page_id, position);
