CREATE INDEX "idx_users_team_name" ON "users" USING BTREE ("team_name");

CREATE INDEX "idx_users_team_name_is_active" ON "users" USING BTREE ("team_name", "is_active");
