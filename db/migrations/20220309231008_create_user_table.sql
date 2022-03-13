-- +goose Up
-- +goose StatementBegin
CREATE TABLE "user"
(
    id   BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);
CREATE INDEX "user_name" ON "user" ("name");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX "user_name";
DROP TABLE "user"
-- +goose StatementEnd
