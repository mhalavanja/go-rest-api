-- name: CreateGroup :exec
CALL createGroup($1, $2);
-- name: UpdateGroupOwner :exec
UPDATE groups
SET user_id_owner = $1
WHERE id = $2;
-- name: UpdateGroupName :exec
UPDATE groups
SET name = $1
WHERE id = $2;
-- name: DeleteGroup :exec
DELETE FROM groups
WHERE id = $1;
-- name: JoinGroup :exec
INSERT INTO groups_users (group_id, user_id)
VALUES ($1, $2);
-- name: LeaveGroup :exec
CALL leaveGroup($1, $2);
-- name: AddUserToGroup :exec
INSERT INTO groups_users (group_id, user_id)
VALUES ($1, $2);
-- name: DeleteUserFromGroup :exec
DELETE FROM groups_users
WHERE group_id = $1
  AND user_id = $2;
-- name: AddUserAsAdmin :exec
UPDATE groups_users
SET is_admin = true
WHERE group_id = $1
  AND user_id = $2;
-- name: RemoveUserAsAdmin :exec
UPDATE groups_users
SET is_admin = false
WHERE group_id = $1
  AND user_id = $2;
