CREATE OR REPLACE PROCEDURE addFriendToGroup(
        userId bigint,
        groupId bigint,
        friendUsername varchar(30)
    ) AS $$
DECLARE friendId bigint;
BEGIN
SELECT users.id INTO friendId
FROM friends
    JOIN users ON users.id = friends.user_id_friend
WHERE user_id = userId
    AND username = friendUsername;
IF friendId IS NULL THEN RAISE EXCEPTION USING errcode = 'NOFRN',
MESSAGE = 'User is not friend with this person';
END IF;
INSERT INTO groups_users (group_id, user_id)
VALUES (groupId, friendId);
END;
$$ language plpgsql