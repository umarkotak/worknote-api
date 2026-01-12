# Tasks: Add Google Authentication

## 1. Database Setup

- [x] 1.1 Create `users` table migration

## 2. Repository Layer

- [x] 2.1 Create `repos/user_repo/` with `initialize()` function
- [x] 2.2 Implement `GetByEmail(email string) (*User, error)`
- [x] 2.3 Implement `GetByGoogleID(googleID string) (*User, error)`
- [x] 2.4 Implement `Create(user *User) error`

## 3. Service Layer

- [x] 3.1 Create `services/auth_service/`
- [x] 3.2 Implement Google ID token validation (call Google API)
- [x] 3.3 Implement JWE token generation using `go-jose`
- [x] 3.4 Implement JWE token decryption/validation
- [x] 3.5 Implement `AuthenticateWithGoogle(idToken string) (*TokenResponse, error)`
  - Validate Google token
  - Extract email and user info
  - Check if user exists by email
  - Create user if not exists
  - Generate and return JWE access token

## 4. Handler Layer

- [x] 4.1 Create `handlers/auth_handler/`
- [x] 4.2 Implement `POST /auth/google` endpoint
  - Accept: `{ "id_token": "..." }`
  - Return: `{ "access_token": "...", "user": {...} }`

## 5. Middleware

- [x] 5.1 Create `middleware/auth.go`
- [x] 5.2 Implement `AuthMiddleware` function
  - Extract Bearer token from Authorization header
  - Decrypt and validate JWE token
  - Inject user info into context
  - Return 401 on failure

## 6. Routing

- [x] 6.1 Update `main.go` to register `/auth/google` route
- [x] 6.2 Export `AuthMiddleware` for use on protected routes
- [x] 6.3 Call repository `initialize()` functions

## 7. Configuration

- [x] 7.1 Add RSA key pair generation/loading for JWE
- [x] 7.2 Add Google OAuth client configuration
