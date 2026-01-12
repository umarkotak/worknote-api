-- +migrate Up
CREATE TABLE work_logs (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  date DATE NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE (user_id, date)
);

CREATE INDEX idx_work_logs_user_id ON work_logs(user_id);
CREATE INDEX idx_work_logs_date ON work_logs(date);

-- +migrate Down
DROP TABLE IF EXISTS work_logs;
