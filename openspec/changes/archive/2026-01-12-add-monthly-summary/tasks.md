# Tasks: Add Monthly Activity Summary

## 1. Configuration

- [x] 1.1 Add `OPENROUTER_API_KEY` and `OPENROUTER_MODEL` to `.env.sample`
- [x] 1.2 Add OpenRouter config fields to `config/config.go`

## 2. Database

- [x] 2.1 Create migration `005_create_work_log_summaries.sql` with `work_log_summaries` table

## 3. Models and Contracts

- [x] 3.1 Add `WorkLogSummary` struct to `model/model.go`
- [x] 3.2 Add request/response types to `contract/contract.go`
  - `GenerateSummaryRequest` (month string)
  - `WorkLogSummaryResponse`

## 4. Repository Layer

- [x] 4.1 Create `repos/work_log_summary_repo/work_log_summary_repo.go`
  - `Initialize()` - prepare named statements
  - `Upsert(userID, month, summary)` - insert or update summary
  - `GetByMonth(userID, month)` - retrieve existing summary

## 5. Service Layer

- [x] 5.1 Create `services/work_log_summary_service/work_log_summary_service.go`
  - `GenerateSummary(userID, month)` - orchestrates the flow:
    1. Fetch work logs for the month from `work_log_repo`
    2. Call OpenRouter API with work log contents
    3. Upsert summary to database
    4. Return summary
  - `GetSummary(userID, month)` - retrieve existing summary
- [x] 5.2 Create internal `callOpenRouter(content)` function for API call

## 6. Handler Layer

- [x] 6.1 Create `handlers/work_log_summary_handler/work_log_summary_handler.go`
  - `GenerateSummary` - POST /work-logs/summary
  - `GetSummary` - GET /work-logs/summary/:month

## 7. Integration

- [x] 7.1 Initialize `work_log_summary_repo` in `main.go`
- [x] 7.2 Register new routes in `main.go`

## 8. Validation

- [x] 8.1 Build verification passed (`go build ./...`)
- [x] 8.2 Repository initialization verified in startup logs
- [ ] 8.3 End-to-end testing (requires database migration and API key)
