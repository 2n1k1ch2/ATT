-- TEAMS
INSERT INTO teams (team_name) VALUES
                                  ('backend'),
                                  ('frontend'),
                                  ('mobile');

-- USERS
INSERT INTO users (user_id, username, team_name, is_active) VALUES
                                                                ('u1', 'alice',  'backend', TRUE),
                                                                ('u2', 'bob',    'backend', TRUE),
                                                                ('u3', 'carol',  'backend', TRUE),
                                                                ('u4', 'dave',   'backend', TRUE),
                                                                ('u5', 'eric',   'backend', FALSE),

                                                                ('u6', 'frank',  'frontend', TRUE),
                                                                ('u7', 'george', 'frontend', TRUE),
                                                                ('u8', 'harry',  'frontend', TRUE),
                                                                ('u9', 'ivan',   'frontend', TRUE),
                                                                ('u10','jack',   'frontend', FALSE),

                                                                ('u11','kate',   'mobile',   TRUE),
                                                                ('u12','lisa',   'mobile',   TRUE),
                                                                ('u13','mike',   'mobile',   TRUE),
                                                                ('u14','nick',   'mobile',   TRUE),
                                                                ('u15','oliver', 'mobile',   FALSE);

-- PR
INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at, merged_at) VALUES
                                                                                                             ('pr1', 'Add billing feature', 'u1', 'OPEN', NOW() - INTERVAL '2 days', NULL),
                                                                                                             ('pr2', 'Fix login bug',       'u2', 'MERGED', NOW() - INTERVAL '7 days', NOW() - INTERVAL '6 days'),
                                                                                                             ('pr3', 'Refactor handlers',   'u6', 'OPEN', NOW() - INTERVAL '1 day', NULL);

-- REVIEWERS
INSERT INTO pr_reviewers (pull_request_id, user_id, assigned_at) VALUES
                                                                     ('pr1', 'u2', NOW() - INTERVAL '1 day'),
                                                                     ('pr1', 'u3', NOW() - INTERVAL '1 day'),

                                                                     ('pr2', 'u3', NOW() - INTERVAL '8 days'),

                                                                     ('pr3', 'u7', NOW() - INTERVAL '1 day');
