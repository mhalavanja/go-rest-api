CREATE OR REPLACE PROCEDURE leaveGroup(userId bigint, groupId bigint)
language plpgsql
AS $$
DECLARE
  numOfPeopleInGroup integer;
BEGIN

DELETE FROM groups_users
WHERE user_id = userId
AND group_id = groupId;

SELECT COUNT(user_id) FROM groups_users
INTO numOfPeopleInGroup
WHERE group_id = groupId;

IF numOfPeopleInGroup = 0 THEN
  DELETE FROM groups
  WHERE id = groupId;
END IF;

COMMIT;
END;$$
