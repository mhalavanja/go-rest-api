CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "username" varchar UNIQUE NOT NULL,
    "password_hash" varchar NOT NULL,
    "is_online" bool DEFAULT false
);
CREATE TABLE "friends" (
    "id" bigserial PRIMARY KEY,
    "user_id1" bigint NOT NULL,
    "user_id2" bigint NOT NULL
);
CREATE TABLE "groups" (
    "id" bigserial PRIMARY KEY,
    "name" varchar UNIQUE NOT NULL,
    "owner" bigint NOT NULL,
    "admin_ids" bigint [],
    "user_ids" bigint []
);
COMMENT ON COLUMN "groups"."admin_ids" IS 'Svaki admin_id je FK na users';
COMMENT ON COLUMN "groups"."user_ids" IS 'Svaki user_id je FK na users';
ALTER TABLE "friends"
ADD FOREIGN KEY ("user_id1") REFERENCES "users" ("id");
ALTER TABLE "friends"
ADD FOREIGN KEY ("user_id2") REFERENCES "users" ("id");
ALTER TABLE "groups"
ADD FOREIGN KEY ("owner") REFERENCES "users" ("id");