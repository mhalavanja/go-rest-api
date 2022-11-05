-- name: GetGroup :one
SELECT name AS group_name,
  username AS owner_username
FROM groups
  JOIN users ON users.user_id = groups.user_id_owner
WHERE groups.id = $1
  AND user_id_owner = $2;
-- name: GetGroups :many
SELECT name
FROM groups
WHERE user_id_owner = $1;
-- name: CreateGroup :exec
CALL createGroup(@group_name::text, @user_id::bigint);
-- name: UpdateGroupOwner :exec
UPDATE groups
SET user_id_owner = $1
WHERE id = $2;
-- name: UpdateGroupName :exec
UPDATE groups
SET name = $1
WHERE id = $2;
-- name: TryDeleteGroup :exec
CALL tryDeleteGroup(@group_id::bigint, @user_id::bigint);
-- name: JoinGroup :exec
INSERT INTO groups_users (group_id, user_id)
VALUES ($1, $2);
-- name: LeaveGroup :exec
CALL leaveGroup(@group_id::bigint, @user_id::bigint);
-- name: AddFriendToGroup :exec
INSERT INTO groups_users (group_id, user_id)
VALUES ($1, $2);
-- name: RemoveUserFromGroup :exec
CALL leaveGroup(@group_id::bigint, @user_id::bigint);