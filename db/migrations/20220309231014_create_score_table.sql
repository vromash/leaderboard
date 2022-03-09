-- +goose Up
-- +goose StatementBegin
CREATE TABLE "score"
(
    id      BIGSERIAL PRIMARY KEY,
    score   BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    CONSTRAINT "fk_user_id" FOREIGN KEY ("user_id") REFERENCES "user" ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "score"
-- +goose StatementEnd
