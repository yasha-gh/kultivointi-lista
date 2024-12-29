 -- CREATE TABLE IF NOT EXISTS list_items (
 --     id TEXT NOT NULL PRIMARY KEY,
 --     title_id TEXT NOT NULL,
 --     -- titles: many to many
 --     type TEXT NOT NULL DEFAULT "", -- base, season
 --     broadcast_type NOT NULL DEFAULT "", -- TV, OVA, ONA, Theater Movie, TV Movie, OVA Movie, ONA Movie
 --     thumbnail_image_id TEXT NOT NULL DEFAULT "",
 --     ongoing BOOLEAN NOT NULL DEFAULT true,
 --     episodes_total INTEGER NOT NULL DEFAULT 0,
 --     episodes_seen INTEGER NOT NULL DEFAULT 0, -- Where to seen until
 --     parent_item_id TEXT, -- list_items id if type is season
 --     season_num INTEGER
    -- modified_at DATETIME,
    --  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP

    -- media: many to many medias
 --     -- seasons: Many to many: list_item type = season
 --
 --     -- CONSTRAINT fk_device_users_constraint FOREIGN KEY (id) REFERENCES device_users(user_id) ON DELETE CASCADE
 -- );

INSERT OR REPLACE INTO list_items (
	id,
	title_id,
	type,
	broadcast_type,
	thumbnail_image_id,
	ongoing,
	episodes_total,
	episodes_seen,
	parent_item_id,
	season_num,
	modified_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
