# Change: Add Frontend Agent API Documentation

## Why

A front-end AI agent needs structured documentation to programmatically interact with the Worknote API for job hunting activities and work log management. Current specs describe behavior but lack agent-consumable API reference with endpoints, request/response schemas, and example payloads.

## What Changes

- Add a new `frontend-agent-docs` capability containing machine/agent-friendly API documentation
- Document all endpoints for:
  - **Authentication**: Google OAuth flow
  - **Job Applications**: CRUD operations with search/filter
  - **Job Application Logs**: Nested CRUD for hiring journey tracking
  - **Work Logs**: Upsert and retrieval of daily work notes
- Include request/response schemas with JSON examples
- Provide authentication header requirements

## Impact

- Affected specs: New `frontend-agent-docs` capability (references `auth`, `job-hunting-log`, `work-log`)
- Affected code: None (documentation only, no code changes)
