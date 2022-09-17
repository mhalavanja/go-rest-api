CREATE OR REPLACE PROCEDURE deleteUser(userId bigint)
language plpgsql
AS $$
DECLARE
numOfPeopleInGroup integer;
groupIdList bigint[];
BEGIN

DELETE FROM friends
WHERE user_id = userId
OR user_id_friend = userId;

DELETE FROM groups_users
WHERE user_id = userId;
RETURNING group_id
INTO groupIdList;

FOREACH groupId IN  ARRAY groupIdList
LOOP
  CALL deleteGroupIfEmpty(groupId);
END LOOP;

COMMIT;
END;$$
