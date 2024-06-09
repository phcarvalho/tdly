-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS configs (
  alphabet TEXT NOT NULL,
  last_board_id INTEGER UNSIGNED NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE configs;
-- +goose StatementEnd
