-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS leaderboard (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    event_id INT NOT NULL,
    rank INT DEFAULT 0,
    score INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES event(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_leaderboard_user_event ON leaderboard (user_id, event_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS leaderboard;
-- +goose StatementEnd
