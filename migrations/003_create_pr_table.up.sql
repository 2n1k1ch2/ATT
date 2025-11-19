CREATE TABLE pull_requests (
    pull_request_id UUID PRIMARY KEY,
    pull_request_name VARCHAR(100) NOT NULL,
    author_id UUID NOT NULL REFERENCES users(user_id),
    status TEXT NOT NULL CHECK(status IN ('OPEN', 'MERGED')),
    created_at TIMESTAMP,
    merged_at TIMESTAMP NULL
);

CREATE INDEX pr_author_idx ON pull_requests(author_id);
