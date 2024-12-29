CREATE TABLE IF NOT EXISTS list_item_titles (
    id TEXT NOT NULL PRIMARY KEY,
    title TEXT NOT NULL,
    lang TEXT NOT NULL DEFAULT "", -- common: zh (Chinese), zh_romaji, en, jp, jp_romaji 
    item_id TEXT NOT NULL, -- list item or season id
    primary_title BOOLEAN NOT NULL DEFAULT false
);
