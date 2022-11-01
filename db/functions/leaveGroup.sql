CREATE OR REPLACE PROCEDURE leaveGroup(userId bigint, groupId bigint) language plpgsql AS $$
DECLARE numOfPeopleInGroup integer;
ownerId bigint;
BEGIN
DELETE FROM groups_users
WHERE user_id = userId
  AND group_id = groupId;
SELECT user_id_owner
FROM groups INTO ownerId
WHERE group_id = groupId;
IF ownerId = userId THEN CALL deleteGroup(groupId);
END IF;
CALL deleteGroupIfEmpty(groupId);
COMMIT;
END;
$$