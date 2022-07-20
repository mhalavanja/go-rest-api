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