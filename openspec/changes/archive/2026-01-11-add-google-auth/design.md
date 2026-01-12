# Design: Google Authentication with JWE Token

## Context

This is a new Go project that needs user authentication. The architecture follows handler → service → repository layers. We use PostgreSQL for persistence and Redis for caching.

## Goals

- Provide secure Google OAuth-based authentication
- Use JWE tokens for encrypted access tokens
- Create reusable auth middleware for protected routes
- Keep implementation minimal and straightforward

## Non-Goals

- Refresh token handling (can be added later)
- Role-based access control (can be added later)
- Session management with Redis (can be added later)

## Decisions

### 1. Google OAuth Flow

**Decision**: Use server-side token validation approach

- Client obtains Google ID token (via Google Sign-In SDK)
- Client sends ID token to our `/auth/google` endpoint
- Server validates token with Google and extracts user info

**Alternatives considered**:

- Full OAuth2 authorization code flow → More complex, overkill for mobile/SPA clients

### 2. JWE Token Format

**Decision**: Use `go-jose` library for JWE with:

- Algorithm: RSA-OAEP for key encryption
- Content encryption: A256GCM
- Token payload: `{ user_id, email, role, exp }`

**Alternatives considered**:

- JWT (JWS) → Not encrypted, tokens visible if intercepted
- Opaque tokens with Redis lookup → More Redis dependency

### 3. User Storage

**Decision**: Store minimal user data in `users` table:

```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  google_id TEXT UNIQUE NOT NULL,
  username TEXT UNIQUE NOT NULL,
  name TEXT,
  picture_url TEXT,
  role TEXT NOT NULL DEFAULT 'user',
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

### 4. Middleware Design

**Decision**: Create `AuthMiddleware` function that:

- Extracts `Authorization: Bearer <token>` header
- Decrypts and validates JWE token
- Injects user info into request context
- Returns 401 Unauthorized on failure

## Risks / Trade-offs

| Risk                            | Mitigation                                        |
| ------------------------------- | ------------------------------------------------- |
| RSA key management              | Generate key pair during setup, store in env vars |
| Token expiry handling           | Set 1-month expiry, client must re-auth after     |
| Google token validation latency | Cache Google's public keys                        |

## Open Questions

1. ~~What should be the token expiry duration?~~ → **Resolved: 1 month**
2. ~~Additional user fields?~~ → **Resolved: username (unique), role**
3. Should we store the access token in Redis for revocation? (Assuming no for MVP)
