package handler

import (
	"fmt"
	"signdoc_api/common"
	"signdoc_api/db"
	"signdoc_api/model"

	"github.com/gofiber/fiber/v2"
)

func WorkflowList(c *fiber.Ctx) error {
	user := common.GetUser(c)

	workflowList, err := getWorkflowListByUser(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.JSONResponse{
			Code:    fiber.StatusInternalServerError,
			Success: false,
			Message: fmt.Sprintf("WorkflowList: getWorkflowListByUser: %v", err),
			Data:    err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.JSONResponse{
		Code:    fiber.StatusOK,
		Success: true,
		Message: fmt.Sprintf("list workflow; user `%s`", user),
		Data:    workflowList,
	})
}

func getWorkflowListByUser(user string) (workflowList []model.Workflow, err error) {
	err = db.DB.Table(fmt.Sprintf("%s w", model.TableNameWorkflow)).
		Select("DISTINCT(w.id), w.workflow_type, w.name, w.description, w.status, w.created_by, w.created_date, w.updated_date").
		Joins(fmt.Sprintf("INNER JOIN %s wk ON w.id = wk.workflow_id", model.TableNameWorkflowState)).
		Where("w.created_by = ? OR wk.assigned_to = ?", user, user).
		Preload("WorkflowStatus").
		Preload("WorkflowTypeSelect").
		Find(&workflowList).
		Error
	if err != nil {
		return workflowList, err
	}
	// str field
	for i := range workflowList {
		workflowList[i].WorkflowTypeStr = workflowList[i].WorkflowTypeSelect.Name
		workflowList[i].StatusStr = workflowList[i].WorkflowStatus.Name
	}
	return workflowList, nil
}
