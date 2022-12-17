CREATE OR REPLACE PROCEDURE updateUser(
        userId bigint,
        usernameVar varchar(30),
        emailVar varchar(30),
        hashedPassword varchar(128)
    ) language plpgsql AS $$ BEGIN IF length(hashedPassword) > 0 THEN
UPDATE users
SET username = usernameVar,
    email = emailVar,
    hashed_password = hashedPassword
WHERE id = userId;
ELSE
UPDATE users
SET username = usernameVar,
    email = emailVar
WHERE id = userId;
END IF;
COMMIT;
END;
$$