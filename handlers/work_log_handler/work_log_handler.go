package work_log_handler

import (
	"github.com/gofiber/fiber/v2"

	"worknote-api/contract"
	"worknote-api/middleware"
	"worknote-api/model"
	"worknote-api/services/work_log_download_service"
	"worknote-api/services/work_log_import_service"
	"worknote-api/services/work_log_service"
	"worknote-api/utils/render"
)

// toWorkLogResponse converts a model to response
func toWorkLogResponse(log *model.WorkLog) contract.WorkLogResponse {
	return contract.WorkLogResponse{
		ID:        log.ID,
		Date:      log.Date,
		Content:   log.Content,
		CreatedAt: log.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: log.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// UpsertWorkLog handles PUT /work-logs
func UpsertWorkLog(c *fiber.Ctx) error {
	userInfo := middleware.GetUserFromContext(c)
	if userInfo == nil {
		return render.Unauthorized(c, "unauthorized")
	}

	var req contract.UpsertWorkLogRequest
	if err := c.BodyParser(&req); err != nil {
		return render.BadRequest(c, "invalid request body")
	}

	workLog, err := work_log_service.UpsertWorkLog(userInfo.UserID, &req)
	if err != nil {
		return render.BadRequest(c, err.Error())
	}

	return render.JSON(c, fiber.StatusOK, toWorkLogResponse(workLog))
}

// GetWorkLogByDate handles GET /work-logs/:date
func GetWorkLogByDate(c *fiber.Ctx) error {
	userInfo := middleware.GetUserFromContext(c)
	if userInfo == nil {
		return render.Unauthorized(c, "unauthorized")
	}

	date := c.Params("date")
	if date == "" {
		return render.BadRequest(c, "date is required")
	}

	workLog, err := work_log_service.GetWorkLogByDate(userInfo.UserID, date)
	if err != nil {
		return render.Error(c, fiber.StatusInternalServerError, "internal error")
	}
	if workLog == nil {
		return render.Error(c, fiber.StatusNotFound, "not found")
	}

	return render.JSON(c, fiber.StatusOK, toWorkLogResponse(workLog))
}

// ListWorkLogs handles GET /work-logs
func ListWorkLogs(c *fiber.Ctx) error {
	userInfo := middleware.GetUserFromContext(c)
	if userInfo == nil {
		return render.Unauthorized(c, "unauthorized")
	}

	logs, err := work_log_service.ListWorkLogs(userInfo.UserID)
	if err != nil {
		return render.Error(c, fiber.StatusInternalServerError, "internal error")
	}

	responses := make([]contract.WorkLogResponse, len(logs))
	for i, log := range logs {
		responses[i] = toWorkLogResponse(&log)
	}

	return render.JSON(c, fiber.StatusOK, contract.WorkLogListResponse{
		Data: responses,
	})
}

// DeleteWorkLogByDate handles DELETE /work-logs/:date
func DeleteWorkLogByDate(c *fiber.Ctx) error {
	userInfo := middleware.GetUserFromContext(c)
	if userInfo == nil {
		return render.Unauthorized(c, "unauthorized")
	}

	date := c.Params("date")
	if date == "" {
		return render.BadRequest(c, "date is required")
	}

	err := work_log_service.DeleteWorkLogByDate(userInfo.UserID, date)
	if err != nil {
		return render.Error(c, fiber.StatusInternalServerError, "internal error")
	}

	return render.JSON(c, fiber.StatusOK, map[string]string{"message": "deleted"})
}

// DownloadWorkLogs handles GET /work-logs/download
func DownloadWorkLogs(c *fiber.Ctx) error {
	userInfo := middleware.GetUserFromContext(c)
	if userInfo == nil {
		return render.Unauthorized(c, "unauthorized")
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	req := &work_log_download_service.DownloadRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

	markdown, filename, err := work_log_download_service.DownloadWorkLogs(userInfo.UserID, req)
	if err != nil {
		return render.BadRequest(c, err.Error())
	}

	c.Set("Content-Type", "text/markdown; charset=utf-8")
	c.Set("Content-Disposition", "attachment; filename=\""+filename+"\"")

	return c.SendString(markdown)
}

// ImportWorkLogs handles POST /work-logs/import
func ImportWorkLogs(c *fiber.Ctx) error {
	userInfo := middleware.GetUserFromContext(c)
	if userInfo == nil {
		return render.Unauthorized(c, "unauthorized")
	}

	// Get file from form
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return render.BadRequest(c, "file is required")
	}

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return render.Error(c, fiber.StatusInternalServerError, "failed to open file")
	}
	defer file.Close()

	// Read file content
	content := make([]byte, fileHeader.Size)
	_, err = file.Read(content)
	if err != nil {
		return render.Error(c, fiber.StatusInternalServerError, "failed to read file")
	}

	// Import worklogs from markdown
	result, err := work_log_import_service.ImportFromMarkdown(userInfo.UserID, string(content))
	if err != nil {
		return render.BadRequest(c, err.Error())
	}

	return render.JSON(c, fiber.StatusOK, result)
}
