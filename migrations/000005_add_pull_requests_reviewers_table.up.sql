CREATE TABLE "pull_request_reviewers" (
    "id" uuid PRIMARY KEY,
    "pull_request_id" uuid NOT NULL,
    "reviewer_id" uuid NOT NULL
);
