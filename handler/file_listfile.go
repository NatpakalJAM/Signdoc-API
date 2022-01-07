package handler

import (
	"fmt"
	"signdoc_api/common"
	"signdoc_api/db"
	"signdoc_api/model"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func FileList(c *fiber.Ctx) error {
	user := common.GetUser(c)

	workflowID, _ := strconv.Atoi(c.Params("workflow_id"))

	// check user workflow permission
	err := checkUserWorkflowPermission(workflowID, user)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(model.JSONResponse{
			Code:    fiber.StatusForbidden,
			Success: false,
			Message: fmt.Sprintf("FileList: %v", err),
			Data:    err,
		})
	}

	// get list files by workflow_id
	listFiles, err := getListFileByWorkflowID(workflowID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.JSONResponse{
			Code:    fiber.StatusInternalServerError,
			Success: false,
			Message: fmt.Sprintf("FileList: getListFileByWorkflowID: %v", err),
			Data:    err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.JSONResponse{
		Code:    fiber.StatusOK,
		Success: true,
		Message: fmt.Sprintf("list files workflow_id: %d", workflowID),
		Data:    listFiles,
	})
}

func checkUserWorkflowPermission(workflowID int, user string) error {
	rowWF := db.DB.Table(model.TableNameWorkflow).
		Where("id = ? AND created_by = ?", workflowID, user).
		Find(&model.Workflow{}).
		RowsAffected
	rowWFS := db.DB.Table(model.TableNameWorkflowState).
		Where("workflow_id = ? AND assigned_to = ?", workflowID, user).
		Find(&model.WorkflowState{}).
		RowsAffected
	if rowWF <= 0 && rowWFS <= 0 {
		return fmt.Errorf("user `%s` does not have access to workflow_id %d", user, workflowID)
	}
	return nil
}

func getListFileByWorkflowID(workflowID int) (listFiles []model.FileMeta, err error) {
	err = db.DB.Table(model.TableNameFileMeta).
		Where("workflow_id = ? AND deleted = 0", workflowID).
		Find(&listFiles).
		Error
	return listFiles, err
}
