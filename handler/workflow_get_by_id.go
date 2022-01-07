package handler

import (
	"fmt"
	"signdoc_api/common"
	"signdoc_api/db"
	"signdoc_api/model"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func WorkflowGetByID(c *fiber.Ctx) error {
	user := common.GetUser(c)

	workflowID, _ := strconv.Atoi(c.Params("workflow_id"))

	// check user access workflow
	if !checkUserAccessWorkflow(workflowID, user) {
		return c.Status(fiber.StatusForbidden).JSON(model.JSONResponse{
			Code:    fiber.StatusForbidden,
			Success: false,
			Message: fmt.Sprintf("user `%s` does not have access to workflow_id %d", user, workflowID),
			Data:    nil,
		})
	}

	// get workflow
	workflowList, err := getWorkflowByID(workflowID, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.JSONResponse{
			Code:    fiber.StatusInternalServerError,
			Success: false,
			Message: fmt.Sprintf("WorkflowList: getWorkflowByID: %v", err),
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

func checkUserAccessWorkflow(workflowID int, user string) bool {
	row := db.DB.Table(fmt.Sprintf("%s w", model.TableNameWorkflow)).
		Select("DISTINCT(w.id)").
		Joins(fmt.Sprintf("INNER JOIN %s wk ON w.id = wk.workflow_id", model.TableNameWorkflowState)).
		Where("w.id = ? AND w.created_by = ? OR wk.assigned_to = ?", workflowID, user, user).
		Preload("WorkflowStatus").
		Preload("WorkflowTypeSelect").
		Find(&[]model.Workflow{}).
		RowsAffected
	return row > 0
}

func getWorkflowByID(workflowID int, user string) (workflowList model.Workflow, err error) {
	err = db.DB.Table(model.TableNameWorkflow).
		Where("id = ?", workflowID).
		Preload("WorkflowState").
		Preload("WorkflowState.WorkflowHistory").
		Preload("WorkflowState.WorkflowHistory.WorkflowHistoryType").
		Preload("WorkflowState.WorkflowStateStatusType").
		Preload("WorkflowStatus").
		Preload("WorkflowTypeSelect").
		Find(&workflowList).
		Error
	if err != nil {
		return workflowList, err
	}
	// str field
	workflowList.WorkflowTypeStr = workflowList.WorkflowTypeSelect.Name
	workflowList.StatusStr = workflowList.WorkflowStatus.Name
	for i := range workflowList.WorkflowState {
		workflowList.WorkflowState[i].StatusStr = workflowList.WorkflowState[i].WorkflowStateStatusType.Name
		for j := range workflowList.WorkflowState[i].WorkflowHistory {
			workflowList.WorkflowState[i].WorkflowHistory[j].HistoryTypeStr =
				workflowList.WorkflowState[i].WorkflowHistory[j].WorkflowHistoryType.Name
		}
	}
	return workflowList, nil
}
