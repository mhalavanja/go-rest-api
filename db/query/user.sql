-- name: GetUser :one
SELECT username,
    email
FROM users
WHERE id = $1;
-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1;
-- name: CreateUser :one
INSERT INTO users (username, email, hashed_password)
VALUES ($1, $2, $3)
RETURNING *;
-- name: UpdateUser :exec
UPDATE users
SET username = $2,
    email = $3,
    hashed_password = $4
WHERE id = $1;
-- name: DeleteUser :exec
CALL deleteUser(@user_id::bigint);