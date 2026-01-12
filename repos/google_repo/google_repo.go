package google_repo

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"

	"github.com/golang-jwt/jwt/v5"

	"worknote-api/config"
)

// GoogleClaims represents the claims in the JWT.
type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	jwt.RegisteredClaims
}

// GooglePublicKeys represents the response from Google's public keys endpoint.
type GooglePublicKeys struct {
	Keys []struct {
		Kid string `json:"kid"`
		N   string `json:"n"`
		E   string `json:"e"`
	} `json:"keys"`
}

// FetchGooglePublicKeys fetches Google's public keys.
func FetchGooglePublicKeys() (*GooglePublicKeys, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v3/certs")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Google public keys: %v", err)
	}
	defer resp.Body.Close()

	var keys GooglePublicKeys
	if err := json.NewDecoder(resp.Body).Decode(&keys); err != nil {
		return nil, fmt.Errorf("failed to decode Google public keys: %v", err)
	}

	return &keys, nil
}

// ParseRSAPublicKey converts a modulus (N) and exponent (E) into an RSA public key.
func ParseRSAPublicKey(modulus, exponent string) (*rsa.PublicKey, error) {
	// Decode the base64 URL-encoded modulus and exponent.
	nBytes, err := base64.RawURLEncoding.DecodeString(modulus)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus: %v", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(exponent)
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent: %v", err)
	}

	// Convert the exponent bytes to an integer.
	var eInt int
	for _, b := range eBytes {
		eInt = eInt<<8 | int(b)
	}

	// Create the RSA public key.
	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: eInt,
	}, nil
}

// ValidateGoogleJWT validates the JWT token and returns the claims.
func ValidateGoogleJWT(tokenString string) (*GoogleClaims, error) {
	// Fetch Google's public keys.
	keys, err := FetchGooglePublicKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch public keys: %v", err)
	}

	// Parse the JWT without verifying the signature to get the kid (key ID).
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &GoogleClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	// Extract the kid from the token header.
	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("kid not found in token header")
	}

	// Find the corresponding public key.
	var publicKey *rsa.PublicKey
	for _, key := range keys.Keys {
		if key.Kid == kid {
			publicKey, err = ParseRSAPublicKey(key.N, key.E)
			if err != nil {
				return nil, fmt.Errorf("failed to parse public key: %v", err)
			}
			break
		}
	}
	if publicKey == nil {
		return nil, errors.New("public key not found for the given kid")
	}

	// Parse and validate the JWT.
	parsedToken, err := jwt.ParseWithClaims(tokenString, &GoogleClaims{}, func(token *jwt.Token) (any, error) {
		// Ensure the signing method is RSA.
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	// Verify the claims.
	if claims, ok := parsedToken.Claims.(*GoogleClaims); ok && parsedToken.Valid {
		// Check the issuer
		if claims.Issuer != "https://accounts.google.com" && claims.Issuer != "accounts.google.com" {
			return nil, errors.New("invalid issuer")
		}

		// Verify audience matches our client ID
		cfg := config.Get()
		if cfg.GoogleClientID != "" {
			validAudience := false
			for _, aud := range claims.Audience {
				if aud == cfg.GoogleClientID {
					validAudience = true
					break
				}
			}
			if !validAudience {
				return nil, errors.New("invalid audience")
			}
		}

		return claims, nil
	}

	return nil, errors.New("invalid token")
}
