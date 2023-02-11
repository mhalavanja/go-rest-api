CREATE OR REPLACE FUNCTION deleteExpiredSessions() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN
DELETE FROM sessions
WHERE expires_at < NOW();
RETURN NEW;
END;
$$;
CREATE OR REPLACE TRIGGER deleteExpiredSessionsTrigger BEFORE
INSERT ON sessions EXECUTE PROCEDURE deleteExpiredSessions();