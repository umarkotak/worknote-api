## Context

This change introduces the first core feature of worknote: job hunting activity logging. The system needs to support tracking multiple job applications per user, each with its own hiring journey consisting of multiple log entries (interviews, assessments, etc.).

**Stakeholders**: Job seekers using worknote to manage their job search.

**Constraints**:

- Follow existing project conventions (no DI, no interfaces, layered architecture)
- Simple state transitions (direct update, no state machine enforcement)
- Audio files stored as URLs (external storage assumed)

## Goals / Non-Goals

**Goals**:

- Track job applications with essential fields (company, title, URL, salary, email, notes)
- Support 6 states: todo, applied, in-progress, rejected, accepted, dropped
- Track hiring journey with N log entries per application
- Support audio recording URLs for interview evaluation
- Provide CRUD + search/filter APIs

**Non-Goals**:

- Reminders/notifications (deferred)
- Analytics/statistics (deferred)
- Audio file upload/storage (external URL only for now)
- Strict state machine validation

## Decisions

### Database Schema

**Decision**: Two tables with 1:N relationship

```sql
job_applications (
  id, user_id, company_name, job_title, job_url,
  salary_range, email, notes, state,
  created_at, updated_at
)

job_application_logs (
  id, job_application_id, process_name, note,
  audio_url, created_at, updated_at
)
```

**Rationale**: Simple normalized schema. Logs are append-only journey entries. Audio stored as external URL to avoid blob storage complexity initially.

### State Values

**Decision**: Use text enum with 6 states: `todo`, `applied`, `in-progress`, `rejected`, `accepted`, `dropped`

**Rationale**:

- `todo` - saved but not yet applied
- `applied` - application submitted
- `in-progress` - active interview process
- `rejected` - not moving forward
- `accepted` - offer accepted
- `dropped` - user decided not to continue

### API Design

**Decision**: RESTful nested resources

- `POST/GET/PUT/DELETE /job-applications` - application CRUD
- `GET /job-applications?search=&state=` - search/filter
- `POST/GET/PUT/DELETE /job-applications/:id/logs` - log CRUD

**Rationale**: Follows REST conventions, logs are naturally nested under applications.

## Risks / Trade-offs

| Risk                                 | Mitigation                                                |
| ------------------------------------ | --------------------------------------------------------- |
| Audio URLs point to deleted files    | Accept for MVP; add validation later                      |
| No state transition validation       | Simple updates per user request; can add validation later |
| Search performance on large datasets | Add indexes; paginate results                             |

## Open Questions

- ~~Audio file storage mechanism~~ → External URL for now, revisit if upload needed
