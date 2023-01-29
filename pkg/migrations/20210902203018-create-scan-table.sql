-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE scan (
    id SERIAL PRIMARY KEY,
    project_id INT,
    git_commit_hash TEXT,
    pipeline_id INT,
    num_bug TEXT DEFAULT '0',
    num_vulnerability TEXT DEFAULT '0',
    num_debt TEXT DEFAULT '0',
    num_code_smell TEXT DEFAULT '0',
    num_file TEXT DEFAULT '0',
    num_duplicate_line TEXT DEFAULT '0',
    line_code TEXT DEFAULT '0',
    line_comment TEXT DEFAULT '0'
    );

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE scan;
