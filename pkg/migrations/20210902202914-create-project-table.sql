-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE project (
  id SERIAL PRIMARY KEY,
  del_project_id INT,
  git_branch TEXT,
  git_url TEXT,
  git_token TEXT,
  user_id INT,
  date TEXT,
  name TEXT,
  key TEXT,
  token TEXT,
  sonar_token TEXT,
  scan_counter INT NOT NULL DEFAULT 0
);
-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE project;