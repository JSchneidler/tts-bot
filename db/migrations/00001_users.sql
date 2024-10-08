-- +goose Up
-- +goose StatementBegin
CREATE TABLE user (
    id INTEGER PRIMARY KEY,
    discord_id TEXT NOT NULL,
    discord_name TEXT NOT NULL,
    character_usage INTEGER DEFAULT 0,
    tts_count INTEGER DEFAULT 0,
    sfx_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user;
-- +goose StatementEnd
