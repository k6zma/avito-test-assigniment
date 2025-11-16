ALTER TABLE "pull_request_reviewers"
    DROP CONSTRAINT IF EXISTS "pr_reviewers_pull_request_id_fkey";

ALTER TABLE "pull_request_reviewers"
    DROP CONSTRAINT IF EXISTS "pr_reviewers_reviewer_id_fkey";
