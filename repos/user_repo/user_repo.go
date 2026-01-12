package user_repo

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"worknote-api/datastore"
	"worknote-api/model"
)

var (
	stmtGetByEmail    *sqlx.NamedStmt
	stmtGetByGoogleID *sqlx.NamedStmt
	stmtCreate        *sqlx.NamedStmt
)

// Initialize prepares all named statements for user repository
func Initialize() {
	var err error

	stmtGetByEmail, err = datastore.DB.PrepareNamed(`
		SELECT id, email, google_id, username, name, picture_url, role, created_at, updated_at
		FROM users
		WHERE email = :email
	`)
	if err != nil {
		log.Fatalf("failed to prepare stmtGetByEmail: %v", err)
	}

	stmtGetByGoogleID, err = datastore.DB.PrepareNamed(`
		SELECT id, email, google_id, username, name, picture_url, role, created_at, updated_at
		FROM users
		WHERE google_id = :google_id
	`)
	if err != nil {
		log.Fatalf("failed to prepare stmtGetByGoogleID: %v", err)
	}

	stmtCreate, err = datastore.DB.PrepareNamed(`
		INSERT INTO users (email, google_id, username, name, picture_url, role)
		VALUES (:email, :google_id, :username, :name, :picture_url, :role)
		RETURNING id, created_at, updated_at
	`)
	if err != nil {
		log.Fatalf("failed to prepare stmtCreate: %v", err)
	}

	log.Info("user_repo initialized")
}

// GetByEmail retrieves a user by email
func GetByEmail(email string) (*model.User, error) {
	user := &model.User{}
	err := stmtGetByEmail.Get(user, map[string]interface{}{"email": email})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetByGoogleID retrieves a user by Google ID
func GetByGoogleID(googleID string) (*model.User, error) {
	user := &model.User{}
	err := stmtGetByGoogleID.Get(user, map[string]interface{}{"google_id": googleID})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Create inserts a new user into the database
func Create(user *model.User) error {
	return stmtCreate.QueryRow(user).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}
