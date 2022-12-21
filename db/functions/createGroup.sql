CREATE OR REPLACE FUNCTION createGroup(userId bigint, groupName varchar(60)) RETURNS bigint AS $$
DECLARE groupId bigint;
BEGIN
INSERT INTO groups (name, user_id_owner)
VALUES (groupName, userId)
RETURNING id INTO groupId;
INSERT INTO groups_users (group_id, user_id)
VALUES (groupId, userId);
RETURN groupId;
END;
$$ language plpgsql