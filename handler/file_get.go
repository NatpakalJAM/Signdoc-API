package handler

import (
	"fmt"
	"signdoc_api/common"
	"signdoc_api/config"
	"signdoc_api/db"
	"signdoc_api/gcp"
	"signdoc_api/model"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func FileGet(c *fiber.Ctx) error {
	user := common.GetUser(c)

	workflowID, _ := strconv.Atoi(c.Params("workflow_id"))
	saveFileName := c.Params("save_file_name")

	// check user access workflow
	if !checkUserAccessWorkflow(workflowID, user) {
		return c.Status(fiber.StatusForbidden).JSON(model.JSONResponse{
			Code:    fiber.StatusForbidden,
			Success: false,
			Message: fmt.Sprintf("user `%s` does not have access to workflow_id %d", user, workflowID),
			Data:    nil,
		})
	}

	// check file exist in workflow
	if !checkFileExistInWorkflow(workflowID, saveFileName) {
		return c.Status(fiber.StatusBadRequest).JSON(model.JSONResponse{
			Code:    fiber.StatusBadRequest,
			Success: false,
			Message: fmt.Sprintf("save_file_name `%s` does not exist in workflow_id %d", saveFileName, workflowID),
			Data:    nil,
		})
	}

	pathFile := fmt.Sprintf("%susers/%s/%d/%s", config.C.DevPrefix, user, workflowID, saveFileName)

	u, err := gcp.GenerateV4GetObjectSignedURL(config.C.GCP.BucketName, pathFile)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.JSONResponse{
			Code:    fiber.StatusInternalServerError,
			Success: false,
			Message: fmt.Sprintf("FileGet: GenerateV4GetObjectSignedURL: %v", err),
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.JSONResponse{
		Code:    fiber.StatusOK,
		Success: true,
		Message: "generate SignedURL successed.",
		Data:    u,
	})
}

// check file exist in workflow
func checkFileExistInWorkflow(workflowID int, saveFileName string) bool {
	row := db.DB.Table(model.TableNameFileMeta).
		Where("workflow_id = ? AND save_file_name = ?", workflowID, saveFileName).
		Find(&model.FileMeta{}).
		RowsAffected
	return row > 0
}
