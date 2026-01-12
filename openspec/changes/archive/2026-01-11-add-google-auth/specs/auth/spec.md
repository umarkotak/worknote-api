## ADDED Requirements

### Requirement: Google OAuth Authentication

The system SHALL authenticate users via Google OAuth by accepting a Google ID token, validating it with Google's API, and returning a JWE access token.

#### Scenario: New user registration via Google

- **WHEN** a valid Google ID token is provided for an unregistered email
- **THEN** the system creates a new user record with Google profile info
- **AND** returns a JWE access token and user details

#### Scenario: Existing user login via Google

- **WHEN** a valid Google ID token is provided for a registered email
- **THEN** the system returns a JWE access token and existing user details

#### Scenario: Invalid Google token

- **WHEN** an invalid or expired Google ID token is provided
- **THEN** the system returns a 401 Unauthorized error

---

### Requirement: JWE Access Token Generation

The system SHALL generate encrypted JWE access tokens containing user identity that can be used to authenticate subsequent requests.

#### Scenario: Token generation on successful auth

- **WHEN** a user successfully authenticates via Google
- **THEN** the system generates a JWE token with RSA-OAEP encryption
- **AND** the token contains user_id, email, role, and expiration time (1 month)

#### Scenario: Token expiration

- **WHEN** a JWE token has passed its 1-month expiration time
- **THEN** the token is considered invalid for authentication

---

### Requirement: Authentication Middleware

The system SHALL provide middleware that validates JWE tokens from the Authorization header and injects user context for protected routes.

#### Scenario: Valid token in Authorization header

- **WHEN** a request includes `Authorization: Bearer <valid_jwe_token>`
- **THEN** the middleware decrypts the token
- **AND** injects user info (user_id, email, role) into request context
- **AND** allows the request to proceed

#### Scenario: Missing Authorization header

- **WHEN** a request to a protected route lacks an Authorization header
- **THEN** the middleware returns 401 Unauthorized

#### Scenario: Invalid or expired token

- **WHEN** a request includes an invalid or expired JWE token
- **THEN** the middleware returns 401 Unauthorized

---

### Requirement: User Storage

The system SHALL store user information including Google profile data for authenticated users.

#### Scenario: User data persistence

- **WHEN** a new user registers via Google OAuth
- **THEN** the system stores email, google_id, username (unique), name, picture_url, and role
- **AND** records created_at and updated_at timestamps
