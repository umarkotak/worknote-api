# Change: Add Monthly Activity Summary

## Why

Users need the ability to review and understand their work activities at a monthly level. By leveraging AI via OpenRouter, the system can generate meaningful summaries of daily work logs, providing insights into monthly productivity and accomplishments.

## What Changes

- Add new `work_log_summaries` table to store AI-generated monthly summaries
- Add new environment variables for OpenRouter integration (`OPENROUTER_API_KEY`, `OPENROUTER_MODEL`)
- Create new API endpoint `POST /work-logs/summary` that:
  - Accepts a `month` parameter (format: `YYYY-MM`)
  - Fetches all work logs for the authenticated user in that month
  - Sends work log contents to OpenRouter for summarization
  - Stores the summary in the database
  - Returns the generated summary
- Create new API endpoint `GET /work-logs/summary/:month` to retrieve existing summaries

## Impact

- Affected specs: Creates new `work-log-summary` capability
- Affected code:
  - `config/config.go` - Add OpenRouter config fields
  - `.env.sample` - Add new environment variables
  - `model/model.go` - Add WorkLogSummary struct
  - `contract/contract.go` - Add request/response types
  - `db/migrations/` - New migration for work_log_summaries table
  - `repos/work_log_summary_repo/` - New repository
  - `services/work_log_summary_service/` - New service with OpenRouter integration
  - `handlers/work_log_summary_handler/` - New handler
  - `main.go` - Register new routes
