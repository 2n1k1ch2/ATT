CREATE TABLE pr_reviewers (
    pull_request_id TEXT  NOT NULL REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(user_id),
    assigned_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY (pull_request_id, user_id)
);

CREATE INDEX reviewers_user_idx ON pr_reviewers(user_id);
CREATE INDEX reviewers_pr_idx ON pr_reviewers(pull_request_id);