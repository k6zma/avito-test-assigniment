CREATE UNIQUE INDEX "idx_pr_reviewers_unique"
    ON "pull_request_reviewers" USING BTREE ("pull_request_id", "reviewer_id");

CREATE INDEX "idx_pr_reviewers_reviewer"
    ON "pull_request_reviewers" USING BTREE ("reviewer_id");
