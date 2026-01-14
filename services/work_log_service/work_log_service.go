package work_log_service

import (
	"errors"

	"worknote-api/contract"
	"worknote-api/model"
	"worknote-api/repos/work_log_repo"
)

// UpsertWorkLog creates or updates a work log entry for a user
func UpsertWorkLog(userID int64, req *contract.UpsertWorkLogRequest) (*model.WorkLog, error) {
	if req.Date == "" {
		return nil, errors.New("date is required")
	}
	if req.Content == "" {
		return nil, errors.New("content is required")
	}

	content := req.Content

	// If append mode, fetch existing content and append new content
	if req.Append {
		existing, err := work_log_repo.GetByDate(userID, req.Date)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.Content != "" {
			content = existing.Content + "\n\n" + req.Content
		}
	}

	return work_log_repo.Upsert(userID, req.Date, content)
}

// GetWorkLogByDate retrieves a work log by user ID and date
func GetWorkLogByDate(userID int64, date string) (*model.WorkLog, error) {
	if date == "" {
		return nil, errors.New("date is required")
	}
	return work_log_repo.GetByDate(userID, date)
}

// ListWorkLogs retrieves all work logs for a user
func ListWorkLogs(userID int64) ([]model.WorkLog, error) {
	return work_log_repo.ListByUserID(userID)
}

// DeleteWorkLogByDate deletes a work log by user ID and date
func DeleteWorkLogByDate(userID int64, date string) error {
	if date == "" {
		return errors.New("date is required")
	}
	return work_log_repo.DeleteByDate(userID, date)
}
