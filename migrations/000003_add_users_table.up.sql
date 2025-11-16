CREATE TABLE "users" (
    "user_id" uuid PRIMARY KEY,
    "username" text,
    "team_name" text NOT NULL,
    "is_active" boolean NOT NULL
);