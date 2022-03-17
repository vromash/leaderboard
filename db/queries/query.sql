--- User ---

-- name: CreateUser :exec
INSERT INTO "user" ("name")
VALUES (@name);

--- Score ---

-- name: GetAllScores :many
SELECT sub.name, sub.score, sub.rank
FROM (
         SELECT score,
                u.name,
                ROW_NUMBER() OVER (ORDER BY "score" DESC) AS rank
         FROM "score"
                  LEFT JOIN "user" u ON u.id = score.user_id
     ) AS sub;

-- name: GetScoresInRange :many
SELECT sub.name, sub.score, sub.rank
FROM (
         SELECT score,
                u.name,
                ROW_NUMBER() OVER (ORDER BY "score" DESC) AS rank
         FROM "score"
                  LEFT JOIN "user" u ON u.id = score.user_id
         WHERE "updated_at" > @date_to
     ) AS sub
WHERE sub.rank BETWEEN @rank_from::bigint AND @rank_to::bigint;

-- name: GetScoreByPlayerName :many
SELECT sub.name, sub.score, sub.rank
FROM (
         SELECT user_id,
                score,
                u.name,
                ROW_NUMBER() OVER (ORDER BY "score" DESC) AS rank
         FROM "score"
                  LEFT JOIN "user" u ON u.id = score.user_id
     ) AS sub
WHERE sub.user_id = (SELECT "id" FROM "user" WHERE "name" = @user_name::varchar);

-- name: GetScoreByPlayerNameInTimeRange :many
SELECT sub.name, sub.score, sub.rank
FROM (
         SELECT user_id,
                score,
                u.name,
                ROW_NUMBER() OVER (ORDER BY "score" DESC) AS rank
         FROM "score"
                  LEFT JOIN "user" u ON u.id = score.user_id
         WHERE "updated_at" > @date_to
     ) AS sub
WHERE sub.user_id = (SELECT "id" FROM "user" WHERE "name" = @user_name::varchar);

-- name: GetRecordNumber :one
SELECT COUNT(id)
FROM "score";

-- name: GetRecordNumberInTimeRange :one
SELECT COUNT(id)
FROM "score"
WHERE "updated_at" > @updated_at::timestamp;

-- name: CreateScore :exec
INSERT INTO "score" ("score", "user_id")
VALUES (@score, (SELECT "id" FROM "user" WHERE "name" = @user_name::varchar));

-- name: UpdateScore :exec
UPDATE "score"
SET "score" = @score
WHERE "user_id" = (SELECT "id" FROM "user" WHERE "name" = @user_name::varchar);
