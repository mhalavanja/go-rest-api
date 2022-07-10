-- name: CreateFriend :one
INSERT INTO friends (user_id, user_id_friend)
VALUES ($1, $2)
RETURNING *;
-- name: ListFriends :many
SELECT *
FROM friends
WHERE user_id = $1;
-- name: DeleteFriend :exec
DELETE FROM friends
WHERE user_id = $1
    AND user_id_friend = $2;