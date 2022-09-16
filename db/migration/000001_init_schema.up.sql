CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar(30) UNIQUE NOT NULL,
  "email" varchar(30) UNIQUE NOT NULL,
  "hashed_password" varchar(128) NOT NULL
);

CREATE TABLE "friends" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "user_id_friend" bigint NOT NULL
);

CREATE TABLE "groups" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(60) UNIQUE NOT NULL,
  "user_id_owner" bigint NOT NULL
);

CREATE TABLE "groups_users" (
  "id" bigserial PRIMARY KEY,
  "group_id" bigint UNIQUE NOT NULL,
  "user_id" bigint NOT NULL,
  "is_admin" bool NOT NULL DEFAULT false
);

ALTER TABLE "friends" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "friends" ADD FOREIGN KEY ("user_id_friend") REFERENCES "users" ("id");

ALTER TABLE "groups" ADD FOREIGN KEY ("user_id_owner") REFERENCES "users" ("id");

ALTER TABLE "groups_users" ADD FOREIGN KEY ("group_id") REFERENCES "groups" ("id");

ALTER TABLE "groups_users" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
