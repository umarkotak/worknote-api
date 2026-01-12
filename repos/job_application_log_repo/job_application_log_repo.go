package job_application_log_repo

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"worknote-api/datastore"
	"worknote-api/model"
)

var (
	stmtCreate                *sqlx.NamedStmt
	stmtGetByID               *sqlx.NamedStmt
	stmtGetByJobApplicationID *sqlx.NamedStmt
	stmtUpdate                *sqlx.NamedStmt
	stmtDelete                *sqlx.NamedStmt
)

// Initialize prepares all named statements for job application log repository
func Initialize() {
	var err error

	stmtCreate, err = datastore.DB.PrepareNamed(`
		INSERT INTO job_application_logs (job_application_id, process_name, note, audio_url)
		VALUES (:job_application_id, :process_name, :note, :audio_url)
		RETURNING id, created_at, updated_at
	`)
	if err != nil {
		log.Fatalf("failed to prepare job_application_log stmtCreate: %v", err)
	}

	stmtGetByID, err = datastore.DB.PrepareNamed(`
		SELECT id, job_application_id, process_name, note, audio_url, created_at, updated_at
		FROM job_application_logs
		WHERE id = :id AND job_application_id = :job_application_id
	`)
	if err != nil {
		log.Fatalf("failed to prepare job_application_log stmtGetByID: %v", err)
	}

	stmtGetByJobApplicationID, err = datastore.DB.PrepareNamed(`
		SELECT id, job_application_id, process_name, note, audio_url, created_at, updated_at
		FROM job_application_logs
		WHERE job_application_id = :job_application_id
		ORDER BY created_at ASC
	`)
	if err != nil {
		log.Fatalf("failed to prepare job_application_log stmtGetByJobApplicationID: %v", err)
	}

	stmtUpdate, err = datastore.DB.PrepareNamed(`
		UPDATE job_application_logs
		SET process_name = :process_name, note = :note, audio_url = :audio_url, updated_at = NOW()
		WHERE id = :id AND job_application_id = :job_application_id
		RETURNING updated_at
	`)
	if err != nil {
		log.Fatalf("failed to prepare job_application_log stmtUpdate: %v", err)
	}

	stmtDelete, err = datastore.DB.PrepareNamed(`
		DELETE FROM job_application_logs
		WHERE id = :id AND job_application_id = :job_application_id
	`)
	if err != nil {
		log.Fatalf("failed to prepare job_application_log stmtDelete: %v", err)
	}

	log.Info("job_application_log_repo initialized")
}

// Create inserts a new job application log into the database
func Create(appLog *model.JobApplicationLog) error {
	return stmtCreate.QueryRow(appLog).Scan(&appLog.ID, &appLog.CreatedAt, &appLog.UpdatedAt)
}

// GetByID retrieves a job application log by ID and job application ID
func GetByID(id, jobApplicationID int64) (*model.JobApplicationLog, error) {
	appLog := &model.JobApplicationLog{}
	err := stmtGetByID.Get(appLog, map[string]interface{}{"id": id, "job_application_id": jobApplicationID})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return appLog, nil
}

// GetByJobApplicationID retrieves all logs for a job application
func GetByJobApplicationID(jobApplicationID int64) ([]model.JobApplicationLog, error) {
	var logs []model.JobApplicationLog
	rows, err := stmtGetByJobApplicationID.Queryx(map[string]interface{}{"job_application_id": jobApplicationID})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var appLog model.JobApplicationLog
		if err := rows.StructScan(&appLog); err != nil {
			return nil, err
		}
		logs = append(logs, appLog)
	}

	return logs, nil
}

// Update updates a job application log in the database
func Update(appLog *model.JobApplicationLog) error {
	return stmtUpdate.QueryRow(appLog).Scan(&appLog.UpdatedAt)
}

// Delete removes a job application log from the database
func Delete(id, jobApplicationID int64) error {
	result, err := stmtDelete.Exec(map[string]interface{}{"id": id, "job_application_id": jobApplicationID})
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
