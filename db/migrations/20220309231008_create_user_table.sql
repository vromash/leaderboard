-- +goose Up
-- +goose StatementBegin
CREATE TABLE "user"
(
    id   BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "user"
-- +goose StatementEnd
