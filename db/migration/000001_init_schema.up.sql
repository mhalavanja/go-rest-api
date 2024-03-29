CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "hashed_password" varchar(128) NOT NULL
);
CREATE TABLE "friends" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "user_id_friend" bigint NOT NULL,
  UNIQUE("user_id", "user_id_friend")
);
CREATE TABLE "groups" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "user_id_owner" bigint NOT NULL
);
CREATE TABLE "groups_users" (
  "id" bigserial PRIMARY KEY,
  "group_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  UNIQUE("group_id", "user_id")
);
CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "user_id" bigint UNIQUE NOT NULL,
  "refresh_token" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
ALTER TABLE "friends"
ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "friends"
ADD FOREIGN KEY ("user_id_friend") REFERENCES "users" ("id");
ALTER TABLE "groups"
ADD FOREIGN KEY ("user_id_owner") REFERENCES "users" ("id");
ALTER TABLE "groups_users"
ADD FOREIGN KEY ("group_id") REFERENCES "groups" ("id");
ALTER TABLE "groups_users"
ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "sessions"
ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
