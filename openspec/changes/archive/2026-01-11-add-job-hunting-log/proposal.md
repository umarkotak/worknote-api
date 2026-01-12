# Change: Add Job Hunting Log Feature

## Why

The worknote app aims to help workers track their job hunting activities. Users need a way to log job applications, track the hiring journey through multiple stages, and maintain notes/recordings of their interview experiences for later evaluation.

## What Changes

- Add `job_applications` table to track job opportunities with states (todo, applied, in-progress, rejected, accepted, dropped)
- Add `job_application_logs` table to track the hiring journey with process name, notes, and optional audio recordings
- Implement CRUD APIs for job applications
- Implement nested CRUD APIs for job application logs
- Add search and filter capabilities for job applications

## Impact

- Affected specs: `job-hunting-log` (new capability)
- Affected code:
  - `db/migrations/` - new migration files
  - `model/` - new Job Application and Log models
  - `contract/` - new request/response types
  - `repos/` - new job_application_repo and job_application_log_repo
  - `services/` - new job_application_service
  - `handlers/` - new job_application_handler
  - `main.go` - new route registrations
