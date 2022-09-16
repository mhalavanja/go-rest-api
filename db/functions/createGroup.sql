CREATE OR REPLACE PROCEDURE createGroup(userId bigint, groupName varchar(60))
language plpgsql
AS $$
DECLARE
lastId bigint;
BEGIN

INSERT INTO groups (name, user_id_owner)
VALUES (groupName, userId)
RETURNING id INTO lastId;

INSERT INTO groups_users (group_id, user_id, is_admin)
VALUES (lastId, userId, true);

COMMIT;
END;$$

