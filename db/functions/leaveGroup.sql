CREATE OR REPLACE PROCEDURE leaveGroup(userId bigint, groupId bigint)
language plpgsql
AS $$
DECLARE
  numOfPeopleInGroup integer;
BEGIN

DELETE FROM groups_users
WHERE user_id = userId
AND group_id = groupId;

CALL deleteGroupIfEmpty(groupId);

COMMIT;
END;$$
