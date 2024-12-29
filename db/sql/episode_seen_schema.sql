CREATE TABLE IF NOT EXISTS episode_seen (
    id TEXT NOT NULL PRIMARY KEY,
    episodes_seen INT, -- last seen episode on this site
    site_id TEXT,
    item_id TEXT
);
