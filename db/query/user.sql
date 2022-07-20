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
-- name: JoinGroup :exec
INSERT INTO group_users (group_id, user_id)
VALUES ($1, $2);
-- name: LeaveGroup :exec
DELETE FROM group_users
WHERE group_id = $1
    AND user_id = $2;