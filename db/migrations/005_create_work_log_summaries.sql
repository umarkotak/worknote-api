-- +migrate Up
CREATE TABLE work_log_summaries (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  month VARCHAR(7) NOT NULL,  -- Format: YYYY-MM
  summary TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE (user_id, month)
);

CREATE INDEX idx_work_log_summaries_user_id ON work_log_summaries(user_id);
CREATE INDEX idx_work_log_summaries_month ON work_log_summaries(month);

-- +migrate Down
DROP TABLE IF EXISTS work_log_summaries;
