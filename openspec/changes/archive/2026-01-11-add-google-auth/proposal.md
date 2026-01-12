# Change: Add Google Authentication with JWE Token

## Why

The application needs user authentication to secure endpoints. Google OAuth provides a trusted, passwordless login experience. JWE (JSON Web Encryption) tokens ensure secure transmission of user credentials between client and server.

## What Changes

- Add Google OAuth integration for user registration/login
- Create `users` table to store user information
- Implement email-based user lookup (register if new, login if exists)
- Generate JWE access tokens upon successful authentication
- Add authentication middleware for protected routes

## Impact

- Affected specs: `auth` (new capability)
- Affected code:
  - `main.go` - routing and middleware setup
  - `handlers/auth_handler/` - HTTP handlers for auth endpoints
  - `services/auth_service/` - business logic for authentication
  - `repos/user_repo/` - user database operations
  - `middleware/auth.go` - JWE token validation middleware
