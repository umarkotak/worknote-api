package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	log "github.com/sirupsen/logrus"

	"worknote-api/config"
	"worknote-api/datastore"
	"worknote-api/handlers/auth_handler"
	"worknote-api/handlers/job_application_handler"
	"worknote-api/handlers/work_log_handler"
	"worknote-api/handlers/work_log_summary_handler"
	"worknote-api/middleware"
	"worknote-api/repos/job_application_log_repo"
	"worknote-api/repos/job_application_repo"
	"worknote-api/repos/user_repo"
	"worknote-api/repos/work_log_repo"
	"worknote-api/repos/work_log_summary_repo"
	"worknote-api/utils/render"
)

func main() {
	// Initialize config first
	config.Initialize()

	// Initialize data stores (PostgreSQL, Redis)
	datastore.Initialize()
	defer datastore.Close()

	// Initialize repositories
	user_repo.Initialize()
	job_application_repo.Initialize()
	job_application_log_repo.Initialize()
	work_log_repo.Initialize()
	work_log_summary_repo.Initialize()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "worknote-api",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	// Public routes
	app.Post("/auth/google", auth_handler.GoogleAuth)

	// Protected routes
	app.Get("/me", middleware.AuthMiddleware, meHandler)

	// Job Application routes (protected)
	jobApps := app.Group("/job-applications", middleware.AuthMiddleware)
	jobApps.Post("/", job_application_handler.CreateJobApplication)
	jobApps.Get("/", job_application_handler.ListJobApplications)
	jobApps.Get("/:id", job_application_handler.GetJobApplication)
	jobApps.Put("/:id", job_application_handler.UpdateJobApplication)
	jobApps.Delete("/:id", job_application_handler.DeleteJobApplication)

	// Job Application Log routes (nested, protected)
	jobApps.Post("/:id/logs", job_application_handler.CreateJobApplicationLog)
	jobApps.Get("/:id/logs", job_application_handler.ListJobApplicationLogs)
	jobApps.Get("/:id/logs/:log_id", job_application_handler.GetJobApplicationLog)
	jobApps.Put("/:id/logs/:log_id", job_application_handler.UpdateJobApplicationLog)
	jobApps.Delete("/:id/logs/:log_id", job_application_handler.DeleteJobApplicationLog)

	// Work Log routes (protected)
	workLogs := app.Group("/work-logs", middleware.AuthMiddleware)
	workLogs.Put("/", work_log_handler.UpsertWorkLog)
	workLogs.Get("/", work_log_handler.ListWorkLogs)
	workLogs.Get("/download", work_log_handler.DownloadWorkLogs)
	workLogs.Post("/import", work_log_handler.ImportWorkLogs)
	workLogs.Post("/summary", work_log_summary_handler.GenerateSummary)
	workLogs.Get("/summary/:month", work_log_summary_handler.GetSummary)
	workLogs.Get("/:date", work_log_handler.GetWorkLogByDate)
	workLogs.Delete("/:date", work_log_handler.DeleteWorkLogByDate)

	// Start server
	cfg := config.Get()
	log.Infof("Server starting on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

// meHandler is an example protected endpoint that returns the current user
func meHandler(c *fiber.Ctx) error {
	userInfo := middleware.GetUserFromContext(c)
	if userInfo == nil {
		return render.Unauthorized(c, "unauthorized")
	}

	return render.JSON(c, fiber.StatusOK, userInfo)
}
