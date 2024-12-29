CREATE TABLE IF NOT EXISTS seasons (
    id TEXT NOT NULL PRIMARY KEY,
    title_id TEXT NOT NULL,
    serie_id TEXT NOT NULL,
    ongoing BOOLEAN NOT NULL DEFAULT true,
    episodes_total INTEGER NOT NULL DEFAULT 0,
    episodes_seen INTEGER NOT NULL DEFAULT 0, -- Where to seen until
    -- titles: many to many (Including title)
);
