-- name: GetGroup :one
SELECT groups.id,
  name AS group_name,
  username AS owner_username
FROM groups
  JOIN users ON users.id = groups.user_id_owner
  JOIN groups_users ON groups.id = groups_users.group_id
WHERE groups.id = $1
  AND groups_users.user_id = $2;
-- name: GetGroups :many
SELECT groups.id,
  groups.name
FROM groups
  JOIN groups_users ON groups.id = groups_users.group_id
WHERE user_id = $1;
-- name: CreateGroup :one
SELECT createGroup(@user_id::bigint, @group_name::text);
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
DELETE FROM groups_users
WHERE user_id = $1
  AND group_id = $2;
-- CALL leaveGroup(@group_id::bigint, @user_id::bigint);
-- name: AddFriendToGroup :exec
CALL addFriendToGroup(
  @user_id::bigint,
  @group_id::bigint,
  @friend_username::varchar(30)
);
-- name: RemoveUserFromGroup :exec
CALL removeUserFromGroup(
  @user_id::bigint,
  @group_id::bigint,
  @friend_id::bigint
);
-- name: GetGroupUsers :many
SELECT user_id_ret::bigint as user_id,
  username_ret::text as username,
  email_ret::text as email
FROM getGroupUsers(@user_id::bigint, @group_id::bigint);