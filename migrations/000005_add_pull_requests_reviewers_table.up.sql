CREATE TABLE "pull_request_reviewers" (
    "id" uuid PRIMARY KEY,
    "pull_request_id" text NOT NULL,
    "reviewer_id" text NOT NULL
);
