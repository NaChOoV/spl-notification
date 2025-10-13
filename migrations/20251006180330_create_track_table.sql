-- +goose Up
CREATE TABLE IF NOT EXISTS track (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    chat_id VARCHAR(255) NOT NULL,
    external_id INTEGER NOT NULL,
    run VARCHAR(50) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    alias VARCHAR(255),
    last_entry TIMESTAMP,
    last_exit TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(chat_id, run)
);

CREATE INDEX idx_track_chat_id ON track(chat_id);
CREATE INDEX idx_track_external_id ON track(external_id);
CREATE INDEX idx_track_chat_id_external_id ON track(chat_id, external_id);

-- +goose Down
DROP INDEX IF EXISTS idx_track_chat_id_external_id;
DROP INDEX IF EXISTS idx_track_external_id;
DROP INDEX IF EXISTS idx_track_chat_id;

DROP TABLE IF EXISTS track;
