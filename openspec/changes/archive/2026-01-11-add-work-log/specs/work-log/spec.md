# Work-Log Capability

## ADDED Requirements

### Requirement: Work Log Management

The system SHALL allow users to record and manage their daily work notes.

#### Scenario: User creates a new work log entry

- **Given** an authenticated user
- **When** the user submits a work log for a date that doesn't exist
- **Then** a new work log entry is created for that date

#### Scenario: User updates an existing work log entry

- **Given** an authenticated user with an existing work log for a specific date
- **When** the user submits a work log for that same date
- **Then** the existing entry is updated with the new content

#### Scenario: User views a specific day's work log

- **Given** an authenticated user with a work log for a specific date
- **When** the user requests the work log for that date
- **Then** the work log content is returned

#### Scenario: User lists all their work logs

- **Given** an authenticated user with multiple work log entries
- **When** the user requests their work log list
- **Then** all work log entries are returned ordered by date descending

---

### Requirement: Work Log Data Fields

Each work log entry MUST contain the following fields:

| Field   | Type              | Required | Description                      |
| ------- | ----------------- | -------- | -------------------------------- |
| date    | date (YYYY-MM-DD) | Yes      | The date of the work log         |
| content | text              | Yes      | Freetext content of the work log |

#### Scenario: Valid work log entry

- **Given** a user submitting a work log
- **When** both date and content are provided
- **Then** the work log is created/updated successfully

#### Scenario: Missing required fields

- **Given** a user submitting a work log
- **When** date or content is missing
- **Then** a validation error is returned

---

### Requirement: Work Log Storage

The system MUST persist work logs with referential integrity to users.

#### Scenario: Unique constraint on user and date

- **Given** the work_logs table
- **When** two entries with the same user_id and date are attempted
- **Then** only one entry exists (upsert behavior)

#### Scenario: User isolation

- **Given** multiple users with work logs
- **When** a user queries their work logs
- **Then** only their own entries are returned
