package job_application_service

import (
	"database/sql"
	"errors"

	"worknote-api/contract"
	"worknote-api/model"
	"worknote-api/repos/job_application_log_repo"
	"worknote-api/repos/job_application_repo"
)

// Valid job application states
var validStates = map[string]bool{
	"todo":        true,
	"applied":     true,
	"in-progress": true,
	"rejected":    true,
	"accepted":    true,
	"dropped":     true,
}

// CreateJobApplication creates a new job application for a user
func CreateJobApplication(userID int64, req *contract.CreateJobApplicationRequest) (*model.JobApplication, error) {
	if req.CompanyName == "" {
		return nil, errors.New("company_name is required")
	}
	if req.JobTitle == "" {
		return nil, errors.New("job_title is required")
	}

	state := req.State
	if state == "" {
		state = "todo"
	}
	if !validStates[state] {
		return nil, errors.New("invalid state value")
	}

	app := &model.JobApplication{
		UserID:      userID,
		CompanyName: req.CompanyName,
		JobTitle:    req.JobTitle,
		JobURL:      req.JobURL,
		SalaryRange: req.SalaryRange,
		Email:       req.Email,
		Notes:       req.Notes,
		State:       state,
	}

	if err := job_application_repo.Create(app); err != nil {
		return nil, err
	}

	return app, nil
}

// GetJobApplication retrieves a job application by ID for a user
func GetJobApplication(id, userID int64) (*model.JobApplication, error) {
	return job_application_repo.GetByID(id, userID)
}

// ListJobApplications retrieves all job applications for a user with optional search/filter
func ListJobApplications(userID int64, search, stateFilter string, limit, offset int) ([]model.JobApplication, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	if stateFilter != "" && !validStates[stateFilter] {
		return nil, 0, errors.New("invalid state filter")
	}

	return job_application_repo.GetByUserID(userID, search, stateFilter, limit, offset)
}

// UpdateJobApplication updates a job application for a user
func UpdateJobApplication(id, userID int64, req *contract.UpdateJobApplicationRequest) (*model.JobApplication, error) {
	// Get existing application
	app, err := job_application_repo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, nil // Not found
	}

	// Update fields if provided
	if req.CompanyName != "" {
		app.CompanyName = req.CompanyName
	}
	if req.JobTitle != "" {
		app.JobTitle = req.JobTitle
	}
	if req.JobURL != "" {
		app.JobURL = req.JobURL
	}
	if req.SalaryRange != "" {
		app.SalaryRange = req.SalaryRange
	}
	if req.Email != "" {
		app.Email = req.Email
	}
	if req.Notes != "" {
		app.Notes = req.Notes
	}
	if req.State != "" {
		if !validStates[req.State] {
			return nil, errors.New("invalid state value")
		}
		app.State = req.State
	}

	if err := job_application_repo.Update(app); err != nil {
		return nil, err
	}

	return app, nil
}

// DeleteJobApplication deletes a job application for a user
func DeleteJobApplication(id, userID int64) error {
	err := job_application_repo.Delete(id, userID)
	if err == sql.ErrNoRows {
		return nil // Treat as success if not found
	}
	return err
}

// CreateJobApplicationLog creates a new log entry for a job application
func CreateJobApplicationLog(jobApplicationID, userID int64, req *contract.CreateJobApplicationLogRequest) (*model.JobApplicationLog, error) {
	// Verify job application belongs to user
	app, err := job_application_repo.GetByID(jobApplicationID, userID)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, nil // Not found
	}

	if req.ProcessName == "" {
		return nil, errors.New("process_name is required")
	}

	appLog := &model.JobApplicationLog{
		JobApplicationID: jobApplicationID,
		ProcessName:      req.ProcessName,
		Note:             req.Note,
		AudioURL:         req.AudioURL,
	}

	if err := job_application_log_repo.Create(appLog); err != nil {
		return nil, err
	}

	return appLog, nil
}

// GetJobApplicationLog retrieves a log entry by ID
func GetJobApplicationLog(logID, jobApplicationID, userID int64) (*model.JobApplicationLog, error) {
	// Verify job application belongs to user
	app, err := job_application_repo.GetByID(jobApplicationID, userID)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, nil // Not found
	}

	return job_application_log_repo.GetByID(logID, jobApplicationID)
}

// ListJobApplicationLogs retrieves all logs for a job application
func ListJobApplicationLogs(jobApplicationID, userID int64) ([]model.JobApplicationLog, error) {
	// Verify job application belongs to user
	app, err := job_application_repo.GetByID(jobApplicationID, userID)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, nil // Not found
	}

	return job_application_log_repo.GetByJobApplicationID(jobApplicationID)
}

// UpdateJobApplicationLog updates a log entry
func UpdateJobApplicationLog(logID, jobApplicationID, userID int64, req *contract.UpdateJobApplicationLogRequest) (*model.JobApplicationLog, error) {
	// Verify job application belongs to user
	app, err := job_application_repo.GetByID(jobApplicationID, userID)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, nil // Not found
	}

	// Get existing log
	appLog, err := job_application_log_repo.GetByID(logID, jobApplicationID)
	if err != nil {
		return nil, err
	}
	if appLog == nil {
		return nil, nil // Not found
	}

	// Update fields if provided
	if req.ProcessName != "" {
		appLog.ProcessName = req.ProcessName
	}
	if req.Note != "" {
		appLog.Note = req.Note
	}
	if req.AudioURL != "" {
		appLog.AudioURL = req.AudioURL
	}

	if err := job_application_log_repo.Update(appLog); err != nil {
		return nil, err
	}

	return appLog, nil
}

// DeleteJobApplicationLog deletes a log entry
func DeleteJobApplicationLog(logID, jobApplicationID, userID int64) error {
	// Verify job application belongs to user
	app, err := job_application_repo.GetByID(jobApplicationID, userID)
	if err != nil {
		return err
	}
	if app == nil {
		return nil // Treat as success if parent not found
	}

	err = job_application_log_repo.Delete(logID, jobApplicationID)
	if err == sql.ErrNoRows {
		return nil // Treat as success if not found
	}
	return err
}
