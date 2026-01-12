package job_application_repo

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"worknote-api/datastore"
	"worknote-api/model"
)

var (
	stmtCreate        *sqlx.NamedStmt
	stmtGetByID       *sqlx.NamedStmt
	stmtUpdate        *sqlx.NamedStmt
	stmtDelete        *sqlx.NamedStmt
	stmtGetByUserID   *sqlx.Stmt
	stmtCountByUserID *sqlx.Stmt
)

// Initialize prepares all named statements for job application repository
func Initialize() {
	var err error

	stmtCreate, err = datastore.DB.PrepareNamed(`
		INSERT INTO job_applications (user_id, company_name, job_title, job_url, salary_range, email, notes, state)
		VALUES (:user_id, :company_name, :job_title, :job_url, :salary_range, :email, :notes, :state)
		RETURNING id, created_at, updated_at
	`)
	if err != nil {
		log.Fatalf("failed to prepare job_application stmtCreate: %v", err)
	}

	stmtGetByID, err = datastore.DB.PrepareNamed(`
		SELECT id, user_id, company_name, job_title, job_url, salary_range, email, notes, state, created_at, updated_at
		FROM job_applications
		WHERE id = :id AND user_id = :user_id
	`)
	if err != nil {
		log.Fatalf("failed to prepare job_application stmtGetByID: %v", err)
	}

	stmtUpdate, err = datastore.DB.PrepareNamed(`
		UPDATE job_applications
		SET company_name = :company_name, job_title = :job_title, job_url = :job_url,
		    salary_range = :salary_range, email = :email, notes = :notes, state = :state,
		    updated_at = NOW()
		WHERE id = :id AND user_id = :user_id
		RETURNING updated_at
	`)
	if err != nil {
		log.Fatalf("failed to prepare job_application stmtUpdate: %v", err)
	}

	stmtDelete, err = datastore.DB.PrepareNamed(`
		DELETE FROM job_applications
		WHERE id = :id AND user_id = :user_id
	`)
	if err != nil {
		log.Fatalf("failed to prepare job_application stmtDelete: %v", err)
	}

	stmtGetByUserID, err = datastore.DB.Preparex(`
		SELECT id, user_id, company_name, job_title, job_url, salary_range, email, notes, state, created_at, updated_at
		FROM job_applications
		WHERE user_id = $1
		ORDER BY created_at DESC
	`)
	if err != nil {
		log.Fatalf("failed to prepare job_application stmtGetByUserID: %v", err)
	}

	stmtCountByUserID, err = datastore.DB.Preparex(`
		SELECT COUNT(*) FROM job_applications WHERE user_id = $1
	`)
	if err != nil {
		log.Fatalf("failed to prepare job_application stmtCountByUserID: %v", err)
	}

	log.Info("job_application_repo initialized")
}

// Create inserts a new job application into the database
func Create(app *model.JobApplication) error {
	if app.State == "" {
		app.State = "todo"
	}
	return stmtCreate.QueryRow(app).Scan(&app.ID, &app.CreatedAt, &app.UpdatedAt)
}

// GetByID retrieves a job application by ID and user ID
func GetByID(id, userID int64) (*model.JobApplication, error) {
	app := &model.JobApplication{}
	err := stmtGetByID.Get(app, map[string]interface{}{"id": id, "user_id": userID})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return app, nil
}

// GetByUserID retrieves all job applications for a user with optional search and filter
func GetByUserID(userID int64, search, stateFilter string, limit, offset int) ([]model.JobApplication, int, error) {
	var apps []model.JobApplication
	var total int

	// Build dynamic query for search and filter
	baseQuery := `
		SELECT id, user_id, company_name, job_title, job_url, salary_range, email, notes, state, created_at, updated_at
		FROM job_applications
		WHERE user_id = $1
	`
	countQuery := `SELECT COUNT(*) FROM job_applications WHERE user_id = $1`

	args := []interface{}{userID}
	argIndex := 2

	var conditions []string

	if search != "" {
		conditions = append(conditions, fmt.Sprintf("(company_name ILIKE $%d OR job_title ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+search+"%")
		argIndex++
	}

	if stateFilter != "" {
		conditions = append(conditions, fmt.Sprintf("state = $%d", argIndex))
		args = append(args, stateFilter)
		argIndex++
	}

	if len(conditions) > 0 {
		condStr := " AND " + strings.Join(conditions, " AND ")
		baseQuery += condStr
		countQuery += condStr
	}

	// Get total count
	err := datastore.DB.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// Add ordering and pagination
	baseQuery += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	err = datastore.DB.Select(&apps, baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return apps, total, nil
}

// Update updates a job application in the database
func Update(app *model.JobApplication) error {
	return stmtUpdate.QueryRow(app).Scan(&app.UpdatedAt)
}

// Delete removes a job application from the database
func Delete(id, userID int64) error {
	result, err := stmtDelete.Exec(map[string]interface{}{"id": id, "user_id": userID})
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
