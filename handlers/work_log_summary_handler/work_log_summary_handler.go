package work_log_summary_handler

import (
	"github.com/gofiber/fiber/v2"

	"worknote-api/contract"
	"worknote-api/middleware"
	"worknote-api/model"
	"worknote-api/services/work_log_summary_service"
	"worknote-api/utils/render"
)

// toSummaryResponse converts a model to response
func toSummaryResponse(summary *model.WorkLogSummary) contract.WorkLogSummaryResponse {
	return contract.WorkLogSummaryResponse{
		ID:        summary.ID,
		Month:     summary.Month,
		Summary:   summary.Summary,
		CreatedAt: summary.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: summary.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// GenerateSummary handles POST /work-logs/summary
func GenerateSummary(c *fiber.Ctx) error {
	userInfo := middleware.GetUserFromContext(c)
	if userInfo == nil {
		return render.Unauthorized(c, "unauthorized")
	}

	var req contract.GenerateSummaryRequest
	if err := c.BodyParser(&req); err != nil {
		return render.BadRequest(c, "invalid request body")
	}

	if req.Month == "" {
		return render.BadRequest(c, "month is required")
	}

	summary, err := work_log_summary_service.GenerateSummary(userInfo.UserID, req.Month)
	if err != nil {
		return render.BadRequest(c, err.Error())
	}

	return render.JSON(c, fiber.StatusOK, toSummaryResponse(summary))
}

// GetSummary handles GET /work-logs/summary/:month
func GetSummary(c *fiber.Ctx) error {
	userInfo := middleware.GetUserFromContext(c)
	if userInfo == nil {
		return render.Unauthorized(c, "unauthorized")
	}

	month := c.Params("month")
	if month == "" {
		return render.BadRequest(c, "month is required")
	}

	summary, err := work_log_summary_service.GetSummary(userInfo.UserID, month)
	if err != nil {
		return render.BadRequest(c, err.Error())
	}
	if summary == nil {
		return render.Error(c, fiber.StatusNotFound, "summary not found for the specified month")
	}

	return render.JSON(c, fiber.StatusOK, toSummaryResponse(summary))
}
