-- name: GetFriend :one
SELECT username,
    email
FROM friends
    JOIN users ON friends.user_id_friend = users.id
WHERE user_id = $1
    AND user_id_friend = $2;
-- name: GetFriends :many
SELECT users.id,
    username,
    email
FROM friends
    JOIN users ON friends.user_id_friend = users.id
WHERE user_id = $1;
-- name: DeleteFriend :exec
DELETE FROM friends
WHERE user_id = $1
    AND user_id_friend = $2;
-- name: AddFriend :exec
INSERT INTO friends (user_id, user_id_friend)
VALUES (
        $1,
        (
            SELECT id
            FROM users
            WHERE username = $2
        )
    );