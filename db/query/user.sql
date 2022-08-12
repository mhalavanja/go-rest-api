-- name: GetUser :one
SELECT * FROM users 
WHERE id = $1;
-- name: CreateUser :one
INSERT INTO users (username, email, password)
VALUES ($1, $2, $3)
RETURNING *;
-- name: UpdateUsername :exec
UPDATE users
SET username = $1
WHERE id = $2;
-- name: UpdateEmail :exec
UPDATE users
SET email = $1
WHERE id = $2;
-- name: UpdatePassword :exec
UPDATE users
SET password = $1
WHERE id = $2;
-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
