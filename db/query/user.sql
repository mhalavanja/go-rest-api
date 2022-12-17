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
CALL updateUser(
    @id::bigint,
    @username::text,
    @email::text,
    @hashed_password::text
);
-- name: DeleteUser :exec
CALL deleteUser(@user_id::bigint);