CREATE OR REPLACE FUNCTION getGroupUsers(userId bigint, groupId bigint) RETURNS table(
        user_id_ret bigint,
        username_ret varchar(30),
        email_ret varchar(30)
    ) AS $$
DECLARE isUserInGroup bigint;
BEGIN
SELECT INTO isUserInGroup id
FROM groups_users
WHERE group_id = groupId
    AND user_id = userId;
IF isUserInGroup IS NULL THEN RAISE EXCEPTION USING errcode = 'NOTIN',
MESSAGE = 'User is not in this group';
END IF;
RETURN query
SELECT users.id as user_id_ret,
    users.username as username_ret,
    users.email as email_ret
FROM users
    JOIN groups_users ON groups_users.user_id = users.id
WHERE group_id = groupId;
END;
$$ language plpgsql