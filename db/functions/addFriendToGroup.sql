CREATE OR REPLACE PROCEDURE addFriendToGroup(userId bigint, groupId bigint, friendId bigint)
DECLARE isFriend bigint;
BEGIN
SELECT INTO isFriend
FROM friends
WHERE user_id = userId
    AND user_id_friend = friendId;
IF isFriend IS NULL THEN RAISE EXCEPTION USING errcode = 'NOFRN',
MESSAGE = 'User is not friend with this person';
END IF;
INSERT INTO groups_users (group_id, user_id)
VALUES (groupId, friendId);
END;
$$ language plpgsql