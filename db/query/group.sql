-- name: CreateGroup :one
INSERT INTO groups (name, user_id_owner)
VALUES ($1, $2)
RETURNING *;
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
-- name: IncNumOfPeople :one
UPDATE groups
SET num_of_people = num_of_people + 1
WHERE id = $1
RETURNING num_of_people;
-- name: DecNumOfPeople :one
UPDATE groups
SET num_of_people = num_of_people - 1
WHERE id = $1
RETURNING num_of_people;
-- name: JoinGroup :exec
INSERT INTO group_users (group_id, user_id)
VALUES ($1, $2);
-- name: LeaveGroup :exec
DELETE FROM group_users
WHERE group_id = $1
    AND user_id = $2;
-- name: AddUserToGroup :exec
INSERT INTO group_users (group_id, user_id)
VALUES ($1, $2);
-- name: DeleteUserFromGroup :exec
DELETE FROM group_users
WHERE group_id = $1
  AND user_id = $2;
-- name: AddUserAsAdmin :exec
UPDATE group_users
SET is_admin = true
WHERE group_id = $1
  AND user_id = $2;
-- name: RemoveUserAsAdmin :exec
UPDATE group_users
SET is_admin = false
WHERE group_id = $1
  AND user_id = $2;
