# Add Work-Log Feature

## Summary

Add a simple work-log feature that allows users to record daily work notes. Each entry is identified by user ID and date (unique constraint), and uses an upsert API to create or update entries.

## Motivation

From `project.md`:

> Help user log their daily activities during work.

This feature fulfills the core need of tracking daily work activities in a simple, freetext format.

## Scope

- New `work_logs` database table with unique constraint on `(user_id, date)`
- Upsert API endpoint: `PUT /work-logs` (creates or updates based on user_id + date)
- Get API endpoint: `GET /work-logs/:date` (get a specific day's log)
- List API endpoint: `GET /work-logs` (list all logs for the user)

## Out of Scope

- Search/filter functionality (can be added later)
- Rich text or markdown support (plain freetext only)
- Attachments or media
