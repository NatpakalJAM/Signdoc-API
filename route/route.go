package route

import (
	"signdoc_api/handler"
	middleware "signdoc_api/middleware"
	"signdoc_api/model"

	"github.com/gofiber/fiber/v2"
)

// Init -> init route
func Init(app *fiber.App) {

	// health check
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(model.JSONResponse{
			Code:    fiber.StatusOK,
			Success: true,
			Message: "Hello, World!",
			Data:    nil,
		})
	})

	// workflow
	app.Post("/workflow", middleware.CheckAccessToken, handler.WorkflowCreate)
	app.Get("/workflow", middleware.CheckAccessToken, handler.WorkflowList)
	app.Get("/workflow/:workflow_id", middleware.CheckAccessToken, handler.WorkflowGetByID)
	app.Put("/workflow/:workflow_id", middleware.CheckAccessToken, handler.WorkflowEdit)
	app.Put("/workflow/:workflow_id/state/:state_order", middleware.CheckAccessToken, handler.WorkflowUpdate)

	// file
	app.Post("/workflow/:workflow_id/file", middleware.CheckAccessToken, handler.FileUpload)
	app.Get("/workflow/:workflow_id/file", middleware.CheckAccessToken, handler.FileList)
	app.Get("/workflow/:workflow_id/file/:save_file_name", middleware.CheckAccessToken, handler.FileGet)
	app.Put("/workflow/:workflow_id/file/rename", middleware.CheckAccessToken, handler.FileRename)
	app.Delete("/workflow/:workflow_id/file", middleware.CheckAccessToken, handler.FileDelete)
}
