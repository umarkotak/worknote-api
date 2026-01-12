-- +migrate Up
CREATE TABLE job_application_logs (
  id SERIAL PRIMARY KEY,
  job_application_id INTEGER NOT NULL REFERENCES job_applications(id) ON DELETE CASCADE,
  process_name TEXT NOT NULL,
  note TEXT,
  audio_url TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_job_application_logs_job_application_id ON job_application_logs(job_application_id);

-- +migrate Down
DROP TABLE IF EXISTS job_application_logs;
