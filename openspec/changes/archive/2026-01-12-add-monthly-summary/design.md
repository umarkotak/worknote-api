# Design: Monthly Activity Summary

## Context

The worknote-api allows users to log daily work activities. Users need a way to review summarized insights of their monthly work. This feature uses OpenRouter as the AI gateway to generate summaries, leveraging environment variables for API key and model configuration.

## Goals / Non-Goals

### Goals

- Provide an API to generate AI-powered monthly summaries of work logs
- Persist summaries for later retrieval (avoid regenerating each time)
- Use OpenRouter for AI model access with configurable model selection
- Follow existing codebase patterns (handler → service → repo)

### Non-Goals

- Real-time summary updates when work logs change
- Multiple summary versions per month
- Summary deletion/editing (summaries are regenerated on request)

## Decisions

### Decision: Use OpenRouter for AI Integration

- **What**: Call OpenRouter API with the OpenAI-compatible chat completions endpoint
- **Why**: OpenRouter provides access to multiple models via a single API, allowing flexibility without code changes

### Decision: Regenerate Strategy

- **What**: When a user requests a summary for a month that already has one, regenerate and overwrite
- **Why**: Keeps implementation simple; users can regenerate to include new logs added since last summary
- **Alternatives considered**:
  - Only generate once, never update → May miss new logs
  - Store multiple versions → Added complexity not needed initially

### Decision: Summary Table Design

- **What**: `work_log_summaries` table with `user_id`, `month` (YYYY-MM string), `summary` text, timestamps
- **Why**: Simple schema matching existing patterns; month as string for easy querying

## Risks / Trade-offs

- **API Rate Limits** → Mitigated by storing summaries and not auto-regenerating
- **AI Response Quality** → Depends on prompt engineering; can be refined iteratively
- **OpenRouter Downtime** → Returns error to user; summaries still retrievable from DB

## Open Questions

- None blocking; prompt engineering can be refined after initial implementation
