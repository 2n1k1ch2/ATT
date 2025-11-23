CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    team_name TEXT NOT NULL REFERENCES teams(team_name),
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX users_team_active_idx ON users(team_name, is_active);
CREATE INDEX users_team_idx ON users(team_name);
