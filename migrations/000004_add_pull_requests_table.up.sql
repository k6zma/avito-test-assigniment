CREATE TABLE "pull_requests" (
    "pull_request_id" uuid PRIMARY KEY,
    "pull_request_name" text NOT NULL,
    "author_id" uuid NOT NULL,
    "status" pull_request_status NOT NULL,
    "created_at" timestamptz NOT NULL,
    "merged_at" timestamptz
);
