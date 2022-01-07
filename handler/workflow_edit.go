package handler

import (
	"encoding/json"
	"fmt"
	"signdoc_api/common"
	"signdoc_api/db"
	"signdoc_api/model"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func WorkflowEdit(c *fiber.Ctx) error {
	user := common.GetUser(c)

	workflowID, _ := strconv.Atoi(c.Params("workflow_id"))

	// check that the user is the workflow creator.
	if !checkUserWorkflowCreator(workflowID, user) {
		return c.Status(fiber.StatusBadRequest).JSON(model.JSONResponse{
			Code:    fiber.StatusBadRequest,
			Success: false,
			Message: fmt.Sprintf("user `%s` does not workflow_id %d creator", user, workflowID),
			Data:    nil,
		})
	}

	// check the workflow is in editable conditions. (1 = draft)
	/* if !checkWorkflowIsEditable(workflowID) {
		return c.Status(fiber.StatusBadRequest).JSON(model.JSONResponse{
			Code:    fiber.StatusBadRequest,
			Success: false,
			Message: fmt.Sprintf("workflow_id %d is in an uneditable state.", workflowID),
			Data:    nil,
		})
	} */

	inputData := new(workflowRequest)
	_ = c.BodyParser(inputData)

	// Unmarshal wfState
	wfState := new([]workflowStateResponse)
	if inputData.State != "" {
		err := json.Unmarshal([]byte(inputData.State), wfState)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.JSONResponse{
				Code:    fiber.StatusBadRequest,
				Success: false,
				Message: fmt.Sprintf("workflow state: %v", err),
				Data:    nil,
			})
		}
	}

	// check workflow_type in database
	if inputData.WorkflowType > 0 {
		var wfType model.WorkflowType
		lenWFType := db.DB.Table(model.TableNameWorkflowType).Where("id = ?", inputData.WorkflowType).Find(&wfType).RowsAffected
		if lenWFType <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(model.JSONResponse{
				Code:    fiber.StatusBadRequest,
				Success: false,
				Message: "workflow_type invalid",
				Data:    nil,
			})
		}
	}

	// check duplicate name
	if inputData.Name != "" {
		var names []model.Workflow
		lenName := db.DB.Table(model.TableNameWorkflow).Where("name = ?", inputData.Name).Find(&names).RowsAffected
		if lenName > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(model.JSONResponse{
				Code:    fiber.StatusBadRequest,
				Success: false,
				Message: "workflow name already taken",
				Data:    nil,
			})
		}
	}

	// check state order
	if len(*wfState) > 0 {
		err := checkStateOrder(*wfState)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.JSONResponse{
				Code:    fiber.StatusBadRequest,
				Success: false,
				Message: fmt.Sprintf("checkStateOrder: %v", err),
				Data:    nil,
			})
		}
	}

	// edit workflow
	err := editWorkflow(workflowID, user, inputData, *wfState)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.JSONResponse{
			Code:    fiber.StatusInternalServerError,
			Success: false,
			Message: fmt.Sprintf("WorkflowEdit: editWorkflow: %v", err),
			Data:    err,
		})
	}

	// return workflow data
	return c.Status(fiber.StatusOK).JSON(model.JSONResponse{
		Code:    fiber.StatusOK,
		Success: true,
		Message: "edit workflow successed.",
		Data:    nil,
	})
}

// check that the user is the workflow creator.
func checkUserWorkflowCreator(workflowID int, user string) bool {
	row := db.DB.Table(model.TableNameWorkflow).
		Where("id = ? AND created_by = ?", workflowID, user).
		Find(&[]model.Workflow{}).
		RowsAffected
	return row > 0
}

// check the workflow is in editable conditions. (1 = draft)
func checkWorkflowIsEditable(workflowID int) bool {
	row := db.DB.Table(model.TableNameWorkflow).
		Where("id = ? AND status = 1", workflowID). // status = draft
		Find(&[]model.Workflow{}).
		RowsAffected
	return row > 0
}

// edit workflow
func editWorkflow(workflowID int, user string, inputData *workflowRequest, state []workflowStateResponse) (err error) {
	now := time.Now()

	// update workflow
	workflowData := model.Workflow{
		WorkflowType: inputData.WorkflowType,
		Name:         inputData.Name,
		Description:  inputData.Description,
		Status:       inputData.Status,
		UpdatedDate:  now,
	}
	err = db.DB.Table(model.TableNameWorkflow).Where("id = ?", workflowID).Updates(&workflowData).Error
	if err != nil {
		db.DB.Rollback()
		return err
	}

	// update state
	if len(state) > 0 {
		// delete state
		db.DB.Table(model.TableNameWorkflowState).
			Where("workflow_id = ?", workflowID).
			Delete(&model.WorkflowState{})

		// insert state
		insertStateStr := []string{}
		insertStateArgs := []interface{}{}
		for _, v := range state {
			insertStateStr = append(insertStateStr, "(?, ?, ?, ?, ?)")
			insertStateArgs = append(insertStateArgs, workflowID)
			insertStateArgs = append(insertStateArgs, v.Name)
			insertStateArgs = append(insertStateArgs, v.Order)
			insertStateArgs = append(insertStateArgs, 1)
			insertStateArgs = append(insertStateArgs, v.AssignedTo)
		}
		smt := fmt.Sprintf("INSERT INTO `%s` (`workflow_id`,`name`,`state_order`,`status`,`assigned_to`) VALUES %s", model.TableNameWorkflowState, "%s")
		smt = fmt.Sprintf(smt, strings.Join(insertStateStr, ","))
		if err := db.DB.Exec(smt, insertStateArgs...).Error; err != nil {
			db.DB.Rollback()
			return err
		}
	}

	return nil
}
