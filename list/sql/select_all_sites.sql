-- CREATE TABLE IF NOT EXISTS sites (
--     id TEXT NOT NULL PRIMARY KEY,
--     url TEXT NOT NULL, -- maybe computed field
--     domain_base TEXT NOT NULL DEFAULT "", -- domain name without top level, may contain subdomain
--     domain_top_level TEXT NOT NULL DEFAULT "", -- domain top level: .com .net .io etc.
--     domain_protocol TEXT NOT NULL DEFAULT "https", -- maybe allow http
--     episode_template TEXT NOT NULL DEFAULT "", -- Simple DSL for generating urls for episodes
--     main_page_template TEXT NOT NULL DEFAULT "" -- Simple DSL for generating for movie/serie
--     -- CONSTRAINT fk_device_users_constraint FOREIGN KEY (id) REFERENCES device_users(user_id) ON DELETE CASCADE
-- );

SELECT
id,
url,
domain_base,
domain_top_level,
domain_protocol,
episode_template,
main_page_template
FROM sites
