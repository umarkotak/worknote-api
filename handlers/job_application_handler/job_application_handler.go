package job_application_handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"worknote-api/contract"
	"worknote-api/middleware"
	"worknote-api/model"
	"worknote-api/services/job_application_service"
	"worknote-api/utils/render"
)

// toJobApplicationResponse converts a model to response
func toJobApplicationResponse(app *model.JobApplication) contract.JobApplicationResponse {
	return contract.JobApplicationResponse{
		ID:          app.ID,
		CompanyName: app.CompanyName,
		JobTitle:    app.JobTitle,
		JobURL:      app.JobURL,
		SalaryRange: app.SalaryRange,
		Email:       app.Email,
		Notes:       app.Notes,
		State:       app.State,
		CreatedAt:   app.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   app.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// toJobApplicationLogResponse converts a model to response
func toJobApplicationLogResponse(log *model.JobApplicationLog) contract.JobApplicationLogResponse {
	return contract.JobApplicationLogResponse{
		ID:               log.ID,
		JobApplicationID: log.JobApplicationID,
		ProcessName:      log.ProcessName,
		Note:             log.Note,
		AudioURL:         log.AudioURL,
		CreatedAt:        log.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        log.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// CreateJobApplication handles POST /job-applications
func CreateJobApplication(c *fiber.Ctx) error {
	userInfo := middleware.GetUserFromContext(c)
	if userInfo == nil {
		return render.Unauthorized(c, "unauthorized")
	}

	var req contract.CreateJobApplicationRequest
	if err := c.BodyParser(&req); err != nil {
		return render.BadRequest(c, "invalid request body")
	}

	app, err := job_application_service.CreateJobApplication(userInfo.UserID, &req)
	if err != nil {
		return render.BadRequest(c, err.Error())
	}

	return render.JSON(c, fiber.StatusCreated, toJobApplicationResponse(app))
}

// GetJobApplication handles GET /job-applications/:id
func GetJobApplication(c *fiber.Ctx) error {
	userInfo := middleware.GetUserFromContext(c)
	if userInfo == nil {
		return render.Unauthorized(c, "unauthorized")
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return render.BadRequest(c, "invalid id")
	}

	app, err := job_application_service.GetJobApplication(id, userInfo.UserID)
	if err != nil {
		return render.Error(c, fiber.StatusInternalServerError, "internal error")
	}
	if app == nil {
		return render.Error(c, fiber.StatusNotFound, "not found")
	}

	return render.JSON(c, fiber.StatusOK, toJobApplicationResponse(app))
}

// ListJobApplications handles GET /job-applications
func ListJobApplications(c *fiber.Ctx) error {
	userInfo := middleware.GetUserFromContext(c)
	if userInfo == nil {
		return render.Unauthorized(c, "unauthorized")
	}

	search := c.Query("search")
	stateFilter := c.Query("state")
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	apps, total, err := job_application_service.ListJobApplications(userInfo.UserID, search, stateFilter, limit, offset)
	if err != nil {
		return render.BadRequest(c, err.Error())
	}

	responses := make([]contract.JobApplicationResponse, len(apps))
	for i, app := range apps {
		responses[i] = toJobApplicationResponse(&app)
	}

	return render.JSON(c, fiber.StatusOK, contract.JobApplicationListResponse{
		Data:  responses,
		Total: total,
	})
}

// UpdateJobApplication handles PUT /job-applications/:id
func UpdateJobApplication(c *fiber.Ctx) error {
	userInfo := middleware.GetUserFromContext(c)
	if userInfo == nil {
		return render.Unauthorized(c, "unauthorized")
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return render.BadRequest(c, "invalid id")
	}

	var req contract.UpdateJobApplicationRequest
	if err := c.BodyParser(&req); err != nil {
		return render.BadRequest(c, "invalid request body")
	}

	app, err := job_application_service.UpdateJobApplication(id, userInfo.UserID, &req)
	if err != nil {
		return render.BadRequest(c, err.Error())
	}
	if app == nil {
		return render.Error(c, fiber.StatusNotFound, "not found")
	}

	return render.JSON(c, fiber.StatusOK, toJobApplicationResponse(app))
}

// DeleteJobApplication handles DELETE /job-applications/:id
func DeleteJobApplication(c *fiber.Ctx) error {
	userInfo := middleware.GetUserFromContext(c)
	if userInfo == nil {
		return render.Unauthorized(c, "unauthorized")
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return render.BadRequest(c, "invalid id")
	}

	if err := job_application_service.DeleteJobApplication(id, userInfo.UserID); err != nil {
		return render.Error(c, fiber.StatusInternalServerError, "internal error")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// CreateJobApplicationLog handles POST /job-applications/:id/logs
func CreateJobApplicationLog(c *fiber.Ctx) error {
	userInfo := middleware.GetUserFromContext(c)
	if userInfo == nil {
		return render.Unauthorized(c, "unauthorized")
	}

	jobAppID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return render.BadRequest(c, "invalid id")
	}

	var req contract.CreateJobApplicationLogRequest
	if err := c.BodyParser(&req); err != nil {
		return render.BadRequest(c, "invalid request body")
	}

	log, err := job_application_service.CreateJobApplicationLog(jobAppID, userInfo.UserID, &req)
	if err != nil {
		return render.BadRequest(c, err.Error())
	}
	if log == nil {
		return render.Error(c, fiber.StatusNotFound, "job application not found")
	}

	return render.JSON(c, fiber.StatusCreated, toJobApplicationLogResponse(log))
}

// GetJobApplicationLog handles GET /job-applications/:id/logs/:log_id
func GetJobApplicationLog(c *fiber.Ctx) error {
	userInfo := middleware.GetUserFromContext(c)
	if userInfo == nil {
		return render.Unauthorized(c, "unauthorized")
	}

	jobAppID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return render.BadRequest(c, "invalid id")
	}

	logID, err := strconv.ParseInt(c.Params("log_id"), 10, 64)
	if err != nil {
		return render.BadRequest(c, "invalid log_id")
	}

	appLog, err := job_application_service.GetJobApplicationLog(logID, jobAppID, userInfo.UserID)
	if err != nil {
		return render.Error(c, fiber.StatusInternalServerError, "internal error")
	}
	if appLog == nil {
		return render.Error(c, fiber.StatusNotFound, "not found")
	}

	return render.JSON(c, fiber.StatusOK, toJobApplicationLogResponse(appLog))
}

// ListJobApplicationLogs handles GET /job-applications/:id/logs
func ListJobApplicationLogs(c *fiber.Ctx) error {
	userInfo := middleware.GetUserFromContext(c)
	if userInfo == nil {
		return render.Unauthorized(c, "unauthorized")
	}

	jobAppID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return render.BadRequest(c, "invalid id")
	}

	logs, err := job_application_service.ListJobApplicationLogs(jobAppID, userInfo.UserID)
	if err != nil {
		return render.Error(c, fiber.StatusInternalServerError, "internal error")
	}
	if logs == nil {
		return render.Error(c, fiber.StatusNotFound, "job application not found")
	}

	responses := make([]contract.JobApplicationLogResponse, len(logs))
	for i, log := range logs {
		responses[i] = toJobApplicationLogResponse(&log)
	}

	return render.JSON(c, fiber.StatusOK, contract.JobApplicationLogListResponse{
		Data: responses,
	})
}

// UpdateJobApplicationLog handles PUT /job-applications/:id/logs/:log_id
func UpdateJobApplicationLog(c *fiber.Ctx) error {
	userInfo := middleware.GetUserFromContext(c)
	if userInfo == nil {
		return render.Unauthorized(c, "unauthorized")
	}

	jobAppID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return render.BadRequest(c, "invalid id")
	}

	logID, err := strconv.ParseInt(c.Params("log_id"), 10, 64)
	if err != nil {
		return render.BadRequest(c, "invalid log_id")
	}

	var req contract.UpdateJobApplicationLogRequest
	if err := c.BodyParser(&req); err != nil {
		return render.BadRequest(c, "invalid request body")
	}

	appLog, err := job_application_service.UpdateJobApplicationLog(logID, jobAppID, userInfo.UserID, &req)
	if err != nil {
		return render.BadRequest(c, err.Error())
	}
	if appLog == nil {
		return render.Error(c, fiber.StatusNotFound, "not found")
	}

	return render.JSON(c, fiber.StatusOK, toJobApplicationLogResponse(appLog))
}

// DeleteJobApplicationLog handles DELETE /job-applications/:id/logs/:log_id
func DeleteJobApplicationLog(c *fiber.Ctx) error {
	userInfo := middleware.GetUserFromContext(c)
	if userInfo == nil {
		return render.Unauthorized(c, "unauthorized")
	}

	jobAppID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return render.BadRequest(c, "invalid id")
	}

	logID, err := strconv.ParseInt(c.Params("log_id"), 10, 64)
	if err != nil {
		return render.BadRequest(c, "invalid log_id")
	}

	if err := job_application_service.DeleteJobApplicationLog(logID, jobAppID, userInfo.UserID); err != nil {
		return render.Error(c, fiber.StatusInternalServerError, "internal error")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
