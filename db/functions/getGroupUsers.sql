CREATE OR REPLACE FUNCTION getGroupUsers(userId bigint, groupId bigint) RETURNS table(user_id bigint, usenrame text, email text) AS $$
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
SELECT users.id,
    username,
    email
FROM users
    JOIN groups_users ON groups_users.user_id = users.id
WHERE group_id = groupId;
END;
$$ language plpgsql