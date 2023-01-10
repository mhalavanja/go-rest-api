CREATE OR REPLACE PROCEDURE removeUserFromGroup(
        userId bigint,
        groupId bigint,
        friendId bigint
    ) language plpgsql AS $$
DECLARE ownerId bigint;
BEGIN
SELECT user_id_owner
FROM groups INTO ownerId
WHERE id = groupId;
IF ownerId != userId THEN RAISE EXCEPTION USING errcode = 'NOOWN',
MESSAGE = 'User is not owner of this group';
END IF;
DELETE FROM groups_users
WHERE user_id = friendId
    AND group_id = groupId;
COMMIT;
END;
$$