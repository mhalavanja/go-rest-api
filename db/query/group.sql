-- name: CreateGroup :exec
CALL createGroup(@groupName, @userId);
-- name: UpdateGroupOwner :exec
UPDATE groups
SET user_id_owner = $1
WHERE id = $2;
-- name: UpdateGroupName :exec
UPDATE groups
SET name = $1
WHERE id = $2;
-- name: TryDeleteGroup :exec
CALL tryDeleteGroup(sqlc.arg(groupId), sqlc.arg(userId));
-- name: JoinGroup :exec
INSERT INTO groups_users (group_id, user_id)
VALUES ($1, $2);
-- name: LeaveGroup :exec
CALL leaveGroup(sqlc.arg(groupId), sqlc.arg(userId));
-- name: AddUserToGroup :exec
INSERT INTO groups_users (group_id, user_id)
VALUES ($1, $2);
-- name: RemoveUserFromGroup :exec
CALL leaveGroup(sqlc.arg(groupId), sqlc.arg(userId));
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
