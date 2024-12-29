CREATE TABLE IF NOT EXISTS media (
    id TEXT NOT NULL PRIMARY KEY,
    mime_type TEXT NOT NULL DEFAULT "",
    file_path TEXT NOT NULL DEFAULT "", -- local file path
    url TEXT NOT NULL DEFAULT "" -- remote file location
    -- CONSTRAINT fk_device_users_constraint FOREIGN KEY (id) REFERENCES device_users(user_id) ON DELETE CASCADE
);
