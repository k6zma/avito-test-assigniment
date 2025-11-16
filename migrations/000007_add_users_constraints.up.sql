ALTER TABLE "users"
    ADD CONSTRAINT "users_team_name_fkey"
        FOREIGN KEY ("team_name") REFERENCES "teams" ("team_name");
