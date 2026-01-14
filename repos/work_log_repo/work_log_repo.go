package work_log_repo

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"worknote-api/datastore"
	"worknote-api/model"
)

var (
	stmtUpsert       *sqlx.NamedStmt
	stmtGetByDate    *sqlx.NamedStmt
	stmtListByUser   *sqlx.Stmt
	stmtDeleteByDate *sqlx.NamedStmt
)

// Initialize prepares all named statements for work log repository
func Initialize() {
	var err error

	stmtUpsert, err = datastore.DB.PrepareNamed(`
		INSERT INTO work_logs (user_id, date, content)
		VALUES (:user_id, :date, :content)
		ON CONFLICT (user_id, date)
		DO UPDATE SET content = EXCLUDED.content, updated_at = NOW()
		RETURNING id, created_at, updated_at
	`)
	if err != nil {
		log.Fatalf("failed to prepare work_log stmtUpsert: %v", err)
	}

	stmtGetByDate, err = datastore.DB.PrepareNamed(`
		SELECT id, user_id, date, content, created_at, updated_at
		FROM work_logs
		WHERE user_id = :user_id AND date = :date
	`)
	if err != nil {
		log.Fatalf("failed to prepare work_log stmtGetByDate: %v", err)
	}

	stmtListByUser, err = datastore.DB.Preparex(`
		SELECT id, user_id, date, content, created_at, updated_at
		FROM work_logs
		WHERE user_id = $1
		ORDER BY date DESC
	`)
	if err != nil {
		log.Fatalf("failed to prepare work_log stmtListByUser: %v", err)
	}

	stmtDeleteByDate, err = datastore.DB.PrepareNamed(`
		DELETE FROM work_logs
		WHERE user_id = :user_id AND date = :date
	`)
	if err != nil {
		log.Fatalf("failed to prepare work_log stmtDeleteByDate: %v", err)
	}

	log.Info("work_log_repo initialized")
}

// Upsert creates or updates a work log entry
func Upsert(userID int64, date, content string) (*model.WorkLog, error) {
	workLog := &model.WorkLog{
		UserID:  userID,
		Date:    date,
		Content: content,
	}
	err := stmtUpsert.QueryRow(workLog).Scan(&workLog.ID, &workLog.CreatedAt, &workLog.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return workLog, nil
}

// GetByDate retrieves a work log by user ID and date
func GetByDate(userID int64, date string) (*model.WorkLog, error) {
	workLog := &model.WorkLog{}
	err := stmtGetByDate.Get(workLog, map[string]interface{}{"user_id": userID, "date": date})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return workLog, nil
}

// ListByUserID retrieves all work logs for a user
func ListByUserID(userID int64) ([]model.WorkLog, error) {
	var logs []model.WorkLog
	err := stmtListByUser.Select(&logs, userID)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// DeleteByDate deletes a work log by user ID and date
func DeleteByDate(userID int64, date string) error {
	_, err := stmtDeleteByDate.Exec(map[string]interface{}{"user_id": userID, "date": date})
	return err
}
