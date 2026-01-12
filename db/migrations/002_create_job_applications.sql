-- +migrate Up
CREATE TABLE job_applications (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  company_name TEXT NOT NULL,
  job_title TEXT NOT NULL,
  job_url TEXT,
  salary_range TEXT,
  email TEXT,
  notes TEXT,
  state TEXT NOT NULL DEFAULT 'todo' CHECK (state IN ('todo', 'applied', 'in-progress', 'rejected', 'accepted', 'dropped')),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_job_applications_user_id ON job_applications(user_id);
CREATE INDEX idx_job_applications_state ON job_applications(state);
CREATE INDEX idx_job_applications_company_name ON job_applications(company_name);
CREATE INDEX idx_job_applications_job_title ON job_applications(job_title);

-- +migrate Down
DROP TABLE IF EXISTS job_applications;
