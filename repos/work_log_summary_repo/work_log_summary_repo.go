package work_log_summary_repo

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"worknote-api/datastore"
	"worknote-api/model"
)

var (
	stmtUpsert     *sqlx.NamedStmt
	stmtGetByMonth *sqlx.NamedStmt
)

// Initialize prepares all named statements for work log summary repository
func Initialize() {
	var err error

	stmtUpsert, err = datastore.DB.PrepareNamed(`
		INSERT INTO work_log_summaries (user_id, month, summary)
		VALUES (:user_id, :month, :summary)
		ON CONFLICT (user_id, month)
		DO UPDATE SET summary = EXCLUDED.summary, updated_at = NOW()
		RETURNING id, created_at, updated_at
	`)
	if err != nil {
		log.Fatalf("failed to prepare work_log_summary stmtUpsert: %v", err)
	}

	stmtGetByMonth, err = datastore.DB.PrepareNamed(`
		SELECT id, user_id, month, summary, created_at, updated_at
		FROM work_log_summaries
		WHERE user_id = :user_id AND month = :month
	`)
	if err != nil {
		log.Fatalf("failed to prepare work_log_summary stmtGetByMonth: %v", err)
	}

	log.Info("work_log_summary_repo initialized")
}

// Upsert creates or updates a work log summary entry
func Upsert(userID int64, month, summary string) (*model.WorkLogSummary, error) {
	workLogSummary := &model.WorkLogSummary{
		UserID:  userID,
		Month:   month,
		Summary: summary,
	}
	err := stmtUpsert.QueryRow(workLogSummary).Scan(&workLogSummary.ID, &workLogSummary.CreatedAt, &workLogSummary.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return workLogSummary, nil
}

// GetByMonth retrieves a work log summary by user ID and month
func GetByMonth(userID int64, month string) (*model.WorkLogSummary, error) {
	workLogSummary := &model.WorkLogSummary{}
	err := stmtGetByMonth.Get(workLogSummary, map[string]interface{}{"user_id": userID, "month": month})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return workLogSummary, nil
}
