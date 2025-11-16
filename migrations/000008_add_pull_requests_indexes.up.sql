CREATE INDEX "idx_pull_requests_author_id"
    ON "pull_requests" USING BTREE ("author_id");

CREATE INDEX "idx_pull_requests_status"
    ON "pull_requests" USING BTREE ("status");
