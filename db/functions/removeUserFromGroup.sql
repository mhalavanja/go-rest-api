CREATE OR REPLACE PROCEDURE removeUserFromGroup(
        requesterUserId bigint,
        userId bigint,
        groupId bigint
    ) language plpgsql AS $$
DECLARE ownerId bigint;
BEGIN
SELECT user_id_owner
FROM groups INTO ownerId
WHERE group_id = groupId;
IF ownerId = requesterUserId THEN
DELETE FROM groups_users
WHERE user_id = userId
    AND group_id = groupId;
END IF;
COMMIT;
END;
$$