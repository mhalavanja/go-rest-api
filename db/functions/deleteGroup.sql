CREATE OR REPLACE PROCEDURE deleteGroup(groupId bigint)
language plpgsql
AS $$
BEGIN

DELETE FROM groups_users
WHERE group_id = groupId;

DELETE FROM groups
WHERE id = groupId;

END;$$
