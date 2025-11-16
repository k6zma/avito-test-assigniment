ALTER TABLE "pull_request_reviewers"
    ADD CONSTRAINT "pr_reviewers_pull_request_id_fkey"
        FOREIGN KEY ("pull_request_id") REFERENCES "pull_requests" ("pull_request_id");

ALTER TABLE "pull_request_reviewers"
    ADD CONSTRAINT "pr_reviewers_reviewer_id_fkey"
        FOREIGN KEY ("reviewer_id") REFERENCES "users" ("user_id");
