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
                ROW_NUMBER() OVER (ORDER BY "score") AS rank
         FROM "score"
                  LEFT JOIN "user" u ON u.id = score.user_id
     ) AS sub;

-- name: GetScoresInRange :many
SELECT sub.name, sub.score, sub.rank
FROM (
         SELECT score,
                u.name,
                ROW_NUMBER() OVER (ORDER BY "score") AS rank
         FROM "score"
                  LEFT JOIN "user" u ON u.id = score.user_id
     ) AS sub
WHERE sub.rank BETWEEN @rank_from::bigint AND @rank_to::bigint;

-- name: GetScoreByPlayerName :one
SELECT sub.name, sub.score, sub.rank
FROM (
         SELECT score,
                u.name,
                ROW_NUMBER() OVER (ORDER BY "score") AS rank
         FROM "score"
                  LEFT JOIN "user" u ON u.id = score.user_id
     ) AS sub
WHERE sub.user_id = (SELECT "id" FROM "user" WHERE "name" = @user_name::varchar);

-- name: CreateScore :exec
INSERT INTO "score" ("score", "user_id")
VALUES (@score, (SELECT "id" FROM "user" WHERE "name" = @user_name));

-- name: UpdateScore :exec
UPDATE "score"
SET "score" = @score
WHERE "user_id" = (SELECT "id" FROM "user" WHERE "name" = @user_name);
