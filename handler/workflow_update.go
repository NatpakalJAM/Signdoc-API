package handler

import (
	"fmt"
	"signdoc_api/common"
	"signdoc_api/db"
	"signdoc_api/model"
	"strconv"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofiber/fiber/v2"
)

type workflowUpdateRequest struct {
	Action  int    `json:"action" form:"action"`
	Message string `json:"status_message" form:"status_message"`
}

// Validate -> workflowUpdateRequest
func (a workflowUpdateRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Action, validation.Required, validation.Min(1), validation.Max(2)),
	)
}

func WorkflowUpdate(c *fiber.Ctx) error {
	user := common.GetUser(c)

	workflowID, _ := strconv.Atoi(c.Params("workflow_id"))
	stateOrder, _ := strconv.Atoi(c.Params("state_order"))

	inputData := new(workflowUpdateRequest)
	_ = c.BodyParser(inputData)
	err := inputData.Validate()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.JSONResponse{
			Code:    fiber.StatusBadRequest,
			Success: false,
			Message: "WorkflowUpdate BadRequest",
			Data:    err,
		})
	}

	// check workflow status in updatable condition
	if !checkWorkflowStatus(workflowID) {
		return c.Status(fiber.StatusBadRequest).JSON(model.JSONResponse{
			Code:    fiber.StatusBadRequest,
			Success: false,
			Message: fmt.Sprintf("workflow_id %d not in updatable status", workflowID),
			Data:    nil,
		})
	}

	// check workflow current state
	currensState, ck := checkWorkflowCurrentState(workflowID, stateOrder)
	if !ck {
		errMsg := fmt.Sprintf("workflow_id %d now waiting on state %d", workflowID, currensState.Order)
		if currensState.Order == 0 {
			errMsg = "workflow complete or cancel"
		}
		return c.Status(fiber.StatusBadRequest).JSON(model.JSONResponse{
			Code:    fiber.StatusBadRequest,
			Success: false,
			Message: errMsg,
			Data:    nil,
		})
	}

	// check user has access to the state & state status in waiting
	if !checkPermissionState(workflowID, stateOrder, user) {
		return c.Status(fiber.StatusBadRequest).JSON(model.JSONResponse{
			Code:    fiber.StatusBadRequest,
			Success: false,
			Message: fmt.Sprintf("user `%s` does not have access to workflow_id %d", user, workflowID),
			Data:    nil,
		})
	}

	// update workflow state
	err = updateWorkflowState(workflowID, stateOrder, inputData.Action, currensState, user, inputData.Message)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.JSONResponse{
			Code:    fiber.StatusInternalServerError,
			Success: false,
			Message: "WorkflowUpdate: updateWorkflowState: error",
			Data:    err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.JSONResponse{
		Code:    fiber.StatusOK,
		Success: true,
		Message: fmt.Sprintf("update workflow_id %d / state_order %d successed.", workflowID, stateOrder),
		Data:    nil,
	})
}

// check workflow status in updatable condition
// 2 = publish
// 3 = processing
func checkWorkflowStatus(workflowID int) bool {
	row := db.DB.Table(model.TableNameWorkflow).
		Where("id = ? AND status IN (2, 3)", workflowID).
		Find(&model.Workflow{}).
		RowsAffected
	return row > 0
}

// check workflow current state
func checkWorkflowCurrentState(workflowID, stateOrder int) (workflowState model.WorkflowState, isCurrent bool) {
	db.DB.Table(model.TableNameWorkflowState).
		Where("workflow_id = ? AND status = 1", workflowID).
		Order("state_order ASC").
		Limit(1).
		Find(&workflowState)
	return workflowState, workflowState.Order == stateOrder
}

// check user has access to the state & state status in waiting
func checkPermissionState(workflowID, stateOrder int, user string) bool {
	row := db.DB.Table(model.TableNameWorkflowState).
		Where("workflow_id = ? AND state_order = ? AND assigned_to = ? ANd status = 1", workflowID, stateOrder, user).
		Find(&[]model.WorkflowState{}).
		RowsAffected
	return row > 0
}

func updateWorkflowState(workflowID, stateOrder, status int, currensState model.WorkflowState, user, statusMessage string) error {
	now := time.Now()
	_ = now
	var stateStatus int
	var historyType int
	switch status {
	case 1:
		stateStatus = 2 // succeeded
		historyType = 1 // verified
	case 2:
		stateStatus = 3 // rejected
		historyType = 2 // declined

		// update message onreject
		db.DB.Model(&model.Workflow{}).
			Where("id = ?", workflowID).
			Updates(&model.Workflow{
				StatusMessage: statusMessage,
			})
	}

	// update state
	db.DB.Model(&model.WorkflowState{}).
		Where("workflow_id = ? AND state_order = ? AND assigned_to = ? AND status = 1", workflowID, stateOrder, user).
		Updates(&model.WorkflowState{
			Status: stateStatus,
		})

	// add history
	db.DB.Table(model.TableNameWorkflowHistory).Create(&model.WorkflowHistory{
		WorkflowStateID: currensState.ID,
		HistoryType:     historyType,
		Date:            now,
	})

	// check total len state
	lenState := checkTotalLenState(workflowID)

	/* update workflow status
	check condition
	- processing
	- complete
	*/
	if stateOrder < lenState { // processing
		updateWorkflowStatus(workflowID, 3, now)
	} else { // complete
		updateWorkflowStatus(workflowID, 4, now)
	}

	return nil
}

// check total len state
func checkTotalLenState(workflowID int) (row int) {
	row = int(db.DB.Table(model.TableNameWorkflowState).
		Where("workflow_id = ?", workflowID).
		Find(&[]model.WorkflowState{}).
		RowsAffected)
	return row
}

// update workflow status
func updateWorkflowStatus(workflowID, wfStatus int, now time.Time) {
	db.DB.Model(&model.Workflow{}).
		Where("id = ?", workflowID).
		Updates(&model.Workflow{
			Status:      wfStatus,
			UpdatedDate: now,
		})
}
