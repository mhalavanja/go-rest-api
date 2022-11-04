CREATE OR REPLACE PROCEDURE tryDeleteGroup(groupId bigint, userId bigint) language plpgsql AS $$
DECLARE ownerId integer;
BEGIN
SELECT user_id_owner
FROM groups INTO ownerId
WHERE id = groupId;
IF userId != ownerId THEN RAISE EXCEPTION USING errcode = 'NOOWN',
message = 'User is not the owner of the group';
END IF;
CALL deleteGroup(groupId);
COMMIT;
END;
$$