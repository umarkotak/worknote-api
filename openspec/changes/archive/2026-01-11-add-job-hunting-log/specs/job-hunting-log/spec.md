## ADDED Requirements

### Requirement: Job Application Management

The system SHALL allow authenticated users to create, read, update, and delete job applications to track their job hunting activities.

#### Scenario: Create job application

- **WHEN** an authenticated user submits a new job application with company name and job title
- **THEN** the system creates a job application record linked to the user
- **AND** the application defaults to "todo" state if not specified

#### Scenario: View job applications

- **WHEN** an authenticated user requests their job applications
- **THEN** the system returns only the applications belonging to that user

#### Scenario: Update job application

- **WHEN** an authenticated user updates their job application fields or state
- **THEN** the system updates the record and returns the updated application

#### Scenario: Delete job application

- **WHEN** an authenticated user deletes their job application
- **THEN** the system removes the application and all associated logs

#### Scenario: Unauthorized access

- **WHEN** a user attempts to access another user's job application
- **THEN** the system returns 404 Not Found

---

### Requirement: Job Application Data Fields

The system SHALL store job application data including company information, application details, and user notes.

#### Scenario: Required and optional fields

- **WHEN** creating a job application
- **THEN** the system requires: company_name, job_title
- **AND** accepts optional fields: job_url, salary_range, email, notes, state

#### Scenario: State values

- **WHEN** setting the job application state
- **THEN** the system accepts one of: todo, applied, in-progress, rejected, accepted, dropped

---

### Requirement: Job Application Search and Filter

The system SHALL allow users to search and filter their job applications.

#### Scenario: Search by text

- **WHEN** a user provides a search query
- **THEN** the system returns applications matching company_name or job_title

#### Scenario: Filter by state

- **WHEN** a user provides a state filter
- **THEN** the system returns only applications in that state

#### Scenario: Combined search and filter

- **WHEN** a user provides both search query and state filter
- **THEN** the system returns applications matching both criteria

---

### Requirement: Job Application Log Management

The system SHALL allow authenticated users to create, read, update, and delete logs within their job applications to track the hiring journey.

#### Scenario: Create application log

- **WHEN** an authenticated user adds a log to their job application
- **THEN** the system creates a log entry with process_name and optional note/audio_url

#### Scenario: View application logs

- **WHEN** an authenticated user requests logs for their job application
- **THEN** the system returns all logs for that application ordered by creation time

#### Scenario: Update application log

- **WHEN** an authenticated user updates a log entry
- **THEN** the system updates the log record

#### Scenario: Delete application log

- **WHEN** an authenticated user deletes a log entry
- **THEN** the system removes the log record

---

### Requirement: Job Application Log Data Fields

The system SHALL store job application log data including process details and optional audio recordings.

#### Scenario: Log fields

- **WHEN** creating a job application log
- **THEN** the system requires: process_name
- **AND** accepts optional fields: note, audio_url

---

### Requirement: Job Application Storage

The system SHALL persist job applications and logs in the database with proper relationships.

#### Scenario: Job application data persistence

- **WHEN** a job application is created
- **THEN** the system stores id, user_id, company_name, job_title, job_url, salary_range, email, notes, state, created_at, updated_at

#### Scenario: Job application log data persistence

- **WHEN** a job application log is created
- **THEN** the system stores id, job_application_id, process_name, note, audio_url, created_at, updated_at

#### Scenario: Cascade delete

- **WHEN** a job application is deleted
- **THEN** all associated logs are also deleted
