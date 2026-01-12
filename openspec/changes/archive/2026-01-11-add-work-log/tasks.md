# Tasks: Add Work-Log Feature

## Database Layer

- [x] Create migration `004_create_work_logs.sql` with `work_logs` table
  - Fields: `id`, `user_id`, `date`, `content`, `created_at`, `updated_at`
  - Unique constraint on `(user_id, date)`

## Model Layer

- [x] Add `WorkLog` struct to `model/model.go`

## Contract Layer

- [x] Add `UpsertWorkLogRequest` to `contract/contract.go`
- [x] Add `WorkLogResponse` and `WorkLogListResponse` to `contract/contract.go`

## Repository Layer

- [x] Create `repos/work_log_repo/` package
  - [x] `work_log_repo.go` with upsert, GetByDate, ListByUserID

## Service Layer

- [x] Create `services/work_log_service/` package
  - [x] `work_log_service.go` with business logic

## Handler Layer

- [x] Create `handlers/work_log_handler/` package
  - [x] `work_log_handler.go` with HTTP handlers

## Routing

- [x] Register routes in `main.go`:
  - `PUT /work-logs` (upsert)
  - `GET /work-logs/:date` (get by date)
  - `GET /work-logs` (list)

## Verification

- [x] Build verification passed
