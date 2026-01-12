package auth_service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"

	"worknote-api/config"
	"worknote-api/contract"
	"worknote-api/model"
	"worknote-api/repos/google_repo"
	"worknote-api/repos/user_repo"
)

// TokenExpiry is the duration for which a token is valid (1 month)
const TokenExpiry = 30 * 24 * time.Hour

// GoogleTokenResponse represents the response from Google's token endpoint
type GoogleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// GenerateJWEToken generates an encrypted JWE token for the user
func GenerateJWEToken(user *model.User) (string, error) {
	publicKey := config.Get().JWEPublicKey

	// Create encrypter
	encrypter, err := jose.NewEncrypter(
		jose.A256GCM,
		jose.Recipient{
			Algorithm: jose.RSA_OAEP,
			Key:       publicKey,
		},
		(&jose.EncrypterOptions{}).WithType("JWT").WithContentType("JWT"),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create encrypter: %w", err)
	}

	// Create claims
	now := time.Now()
	claims := contract.TokenClaims{
		Claims: jwt.Claims{
			Issuer:    "worknote-api",
			Subject:   fmt.Sprintf("%d", user.ID),
			IssuedAt:  jwt.NewNumericDate(now),
			Expiry:    jwt.NewNumericDate(now.Add(TokenExpiry)),
			NotBefore: jwt.NewNumericDate(now),
		},
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
	}

	// Encrypt token
	token, err := jwt.Encrypted(encrypter).Claims(claims).CompactSerialize()
	if err != nil {
		return "", fmt.Errorf("failed to encrypt token: %w", err)
	}

	return token, nil
}

// DecryptJWEToken decrypts and validates a JWE token
func DecryptJWEToken(tokenString string) (*contract.TokenClaims, error) {
	privateKey := config.Get().JWEPrivateKey

	// Parse the encrypted token
	token, err := jwt.ParseEncrypted(tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Decrypt and get claims
	claims := &contract.TokenClaims{}
	if err := token.Claims(privateKey, claims); err != nil {
		return nil, fmt.Errorf("failed to decrypt token: %w", err)
	}

	// Validate expiry
	if claims.Expiry != nil && time.Now().After(claims.Expiry.Time()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

// AuthenticateWithGoogle authenticates a user with Google OAuth
func AuthenticateWithGoogle(ctx context.Context, idToken string) (*contract.AuthResponse, error) {
	// Validate ID token using JWT validation
	googleClaims, err := google_repo.ValidateGoogleJWT(idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to validate id token: %w", err)
	}

	// Check if user exists by email
	user, err := user_repo.GetByEmail(googleClaims.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		// Create new user
		username := generateUsername(googleClaims.Email)
		user = &model.User{
			Email:      googleClaims.Email,
			GoogleID:   googleClaims.RegisteredClaims.Subject,
			Username:   username,
			Name:       googleClaims.Name,
			PictureURL: googleClaims.Picture,
			Role:       "user",
		}
		if err := user_repo.Create(user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Generate JWE token
	accessToken, err := GenerateJWEToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &contract.AuthResponse{
		AccessToken: accessToken,
		User:        user,
	}, nil
}

// generateUsername creates a username from email
func generateUsername(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		return parts[0]
	}
	return email
}
