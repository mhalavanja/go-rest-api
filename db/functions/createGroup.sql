CREATE OR REPLACE PROCEDURE createGroup(userId bigint, groupName varchar(60)) language plpgsql AS $$
DECLARE groupId bigint;
BEGIN
INSERT INTO groups (name, user_id_owner)
VALUES (groupName, userId)
RETURNING id INTO groupId;
INSERT INTO groups_users (group_id, user_id, is_admin)
VALUES (groupId, userId, true);
COMMIT;
END;
$$