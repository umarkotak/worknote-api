package config

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

// Config holds all configuration values
type Config struct {
	// Google OAuth
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURI  string

	// Database
	DatabaseURL string

	// Redis
	RedisURL      string
	RedisPassword string

	// Server
	Port string

	// JWE Keys
	JWEPrivateKey *rsa.PrivateKey
	JWEPublicKey  *rsa.PublicKey

	// OpenRouter AI
	OpenRouterAPIKey string
	OpenRouterModel  string

	// Z.AI
	ZaiAPIKey string
	ZaiModel  string
}

// GoogleOAuthJSON represents the structure of Google OAuth credentials JSON
type GoogleOAuthJSON struct {
	Web       *GoogleOAuthCredentials `json:"web,omitempty"`
	Installed *GoogleOAuthCredentials `json:"installed,omitempty"`
}

// GoogleOAuthCredentials contains the actual OAuth credentials
type GoogleOAuthCredentials struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURIs []string `json:"redirect_uris"`
}

var cfg *Config

// Initialize loads all configuration from environment variables
func Initialize() {
	// Load .env file if it exists (optional, won't fail if missing)
	if err := godotenv.Load(); err != nil {
		log.Warn("No .env file found, using environment variables directly")
	}

	cfg = &Config{
		DatabaseURL:      getEnvOrDefault("DATABASE_URL", "postgres://localhost:5432/worknote?sslmode=disable"),
		RedisURL:         getEnvOrDefault("REDIS_URL", "localhost:6379"),
		RedisPassword:    os.Getenv("REDIS_PASSWORD"),
		Port:             getEnvOrDefault("PORT", "8080"),
		OpenRouterAPIKey: os.Getenv("OPENROUTER_API_KEY"),
		OpenRouterModel:  getEnvOrDefault("OPENROUTER_MODEL", "openai/gpt-4o-mini"),
		ZaiAPIKey:        os.Getenv("ZAI_API_KEY"),
		ZaiModel:         getEnvOrDefault("ZAI_MODEL", "glm-4.7"),
	}

	// Parse Google OAuth JSON
	parseGoogleOAuthJSON()

	// Initialize JWE keys
	initJWEKeys()

	log.Info("Config initialized")
}

// parseGoogleOAuthJSON parses the GOOGLE_OAUTH_JSON environment variable
func parseGoogleOAuthJSON() {
	oauthJSON := os.Getenv("GOOGLE_OAUTH_JSON")
	if oauthJSON == "" {
		log.Warn("GOOGLE_OAUTH_JSON not set, Google OAuth will not work")
		return
	}

	var googleOAuth GoogleOAuthJSON
	if err := json.Unmarshal([]byte(oauthJSON), &googleOAuth); err != nil {
		log.Fatalf("failed to parse GOOGLE_OAUTH_JSON: %v", err)
	}

	// Try to get credentials from 'web' or 'installed' key
	var creds *GoogleOAuthCredentials
	if googleOAuth.Web != nil {
		creds = googleOAuth.Web
	} else if googleOAuth.Installed != nil {
		creds = googleOAuth.Installed
	} else {
		log.Fatal("GOOGLE_OAUTH_JSON must contain 'web' or 'installed' credentials")
	}

	cfg.GoogleClientID = creds.ClientID
	cfg.GoogleClientSecret = creds.ClientSecret
	if len(creds.RedirectURIs) > 0 {
		cfg.GoogleRedirectURI = creds.RedirectURIs[0]
	}

	// Allow override of redirect URI from env
	if redirectURI := os.Getenv("GOOGLE_REDIRECT_URI"); redirectURI != "" {
		cfg.GoogleRedirectURI = redirectURI
	}
}

// Get returns the global configuration
func Get() *Config {
	if cfg == nil {
		log.Fatal("config not initialized, call config.Initialize() first")
	}
	return cfg
}

func initJWEKeys() {
	// In production, load from environment or file
	// For now, generate a key pair (should be persisted in real usage)
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("failed to generate RSA key: %v", err)
	}
	cfg.JWEPrivateKey = key
	cfg.JWEPublicKey = &key.PublicKey
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
