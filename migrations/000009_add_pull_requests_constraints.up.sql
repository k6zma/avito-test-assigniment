ALTER TABLE "pull_requests"
    ADD CONSTRAINT "pull_requests_author_id_fkey"
        FOREIGN KEY ("author_id") REFERENCES "users" ("user_id");
