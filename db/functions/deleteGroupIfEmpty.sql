CREATE OR REPLACE PROCEDURE deleteGroupIfEmpty(groupId bigint) language plpgsql AS $$
DECLARE numOfPeopleInGroup integer;
BEGIN
SELECT COUNT(user_id)
FROM groups_users INTO numOfPeopleInGroup
WHERE group_id = groupId;
IF numOfPeopleInGroup = 0 THEN
DELETE FROM groups
WHERE id = groupId;
END IF;
COMMIT;
END;
$$