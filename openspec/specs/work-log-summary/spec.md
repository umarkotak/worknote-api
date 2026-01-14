# work-log-summary Specification

## Purpose
TBD - created by archiving change add-monthly-summary. Update Purpose after archive.
## Requirements
### Requirement: Monthly Summary Generation

The system SHALL allow users to generate AI-powered summaries of their monthly work log activities.

#### Scenario: User generates summary for a month with work logs

- **Given** an authenticated user with work logs in the specified month
- **When** the user requests a summary for that month (format: YYYY-MM)
- **Then** the system fetches all work logs for that month
- **And** sends the content to OpenRouter for summarization
- **And** stores the summary in the database
- **And** returns the generated summary

#### Scenario: User generates summary for a month with no work logs

- **Given** an authenticated user with no work logs in the specified month
- **When** the user requests a summary for that month
- **Then** the system returns an error indicating no work logs exist for that month

#### Scenario: User regenerates summary for a month

- **Given** an authenticated user with an existing summary for a month
- **When** the user requests a new summary for that month
- **Then** the system regenerates the summary with current work log data
- **And** overwrites the existing summary

---

### Requirement: Monthly Summary Retrieval

The system SHALL allow users to retrieve previously generated monthly summaries.

#### Scenario: User retrieves existing summary

- **Given** an authenticated user with a summary for a specific month
- **When** the user requests to view the summary for that month
- **Then** the stored summary is returned

#### Scenario: User retrieves summary that does not exist

- **Given** an authenticated user without a summary for a specific month
- **When** the user requests to view the summary for that month
- **Then** a not found error is returned

---

### Requirement: Work Log Summary Data Fields

Each work log summary entry MUST contain the following fields:

| Field   | Type             | Required | Description                       |
| ------- | ---------------- | -------- | --------------------------------- |
| month   | string (YYYY-MM) | Yes      | The month of the summary          |
| summary | text             | Yes      | AI-generated summary of work logs |

#### Scenario: Valid summary entry

- **Given** a user generating a summary
- **When** the AI returns a valid summary
- **Then** the summary is stored with the user_id, month, and content

---

### Requirement: OpenRouter Integration

The system MUST integrate with OpenRouter API for AI summarization.

#### Scenario: Successful OpenRouter call

- **Given** valid OpenRouter credentials in environment
- **When** work log content is sent for summarization
- **Then** the AI-generated summary is returned

#### Scenario: OpenRouter API failure

- **Given** an authenticated user requesting a summary
- **When** the OpenRouter API call fails
- **Then** an internal error is returned to the user

---

### Requirement: Work Log Summary Storage

The system MUST persist work log summaries with referential integrity to users.

#### Scenario: Unique constraint on user and month

- **Given** the work_log_summaries table
- **When** two summaries with the same user_id and month are attempted
- **Then** only one entry exists (upsert behavior)

#### Scenario: User isolation

- **Given** multiple users with summaries
- **When** a user queries their summaries
- **Then** only their own entries are returned

