-- name: GetAllScores :many
SELECT *
FROM "score"
         LEFT JOIN "user" u on u.id = score.user_id
ORDER BY "score";

-- name: CreateUser :exec
INSERT INTO "user" ("name")
VALUES (@name);

-- name: CreateScore :exec
INSERT INTO "score" ("score", "user_id")
VALUES (@score, @user_id);

-- name: UpdateScore :exec
UPDATE "score"
SET "score" = @score
WHERE "user_id" = @user_id;
