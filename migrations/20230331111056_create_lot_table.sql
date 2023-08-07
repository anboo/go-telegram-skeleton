-- +goose Up
-- +goose StatementBegin
CREATE TABLE lots (
    id UUID NOT NULL,
    external_id VARCHAR(255) NOT NULL,
    title VARCHAR(255) DEFAULT NULL,
    description TEXT,
    created_at timestamp,
    PRIMARY KEY(id)
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE UNIQUE INDEX uniq_lots_external_id ON lots (external_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX uniq_lots_external_id;
DROP TABLE lots;
-- +goose StatementEnd
