CREATE TABLE IF NOT EXISTS links(
    id BIGSERIAL PRIMARY KEY,
    original_url VARCHAR,
    short_url VARCHAR,
    created_at timestamp with time zone default CURRENT_TIMESTAMP not null
);