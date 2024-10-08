-- +goose Up
-- +goose StatementBegin
CREATE TABLE audio (
    id INTEGER PRIMARY KEY,
    user INTEGER NOT NULL,
    audio_type TEXT NOT NULL,
    audio_service TEXT NOT NULL,
    prompt TEXT NOT NULL,
    audio_filename TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user) REFERENCES user(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE audio;
-- +goose StatementEnd
