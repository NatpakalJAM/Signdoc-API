package handler

import (
	"encoding/json"
	"fmt"
	"signdoc_api/common"
	"signdoc_api/config"
	"signdoc_api/db"
	"signdoc_api/model"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofiber/fiber/v2"
)

type workflowRequest struct {
	WorkflowType int    `json:"type" form:"type"`
	Name         string `json:"name" form:"name"`
	Description  string `json:"description" form:"description"`
	Status       int    `json:"status" form:"status"`
	State        string `json:"state" form:"state"`
}

type workflowStateResponse struct {
	Order      int    `json:"state_order" form:"state_order"`
	Name       string `json:"name" form:"name"`
	AssignedTo string `json:"assigned_to" form:"assigned_to"`
}

type workflowResponse struct {
	WorkflowID   int                     `json:"workflow_id" form:"workflow_id"`
	WorkflowType int                     `json:"type" form:"type"`
	Name         string                  `json:"name" form:"name"`
	Description  string                  `json:"description" form:"description"`
	Status       int                     `json:"status" form:"status"`
	State        []workflowStateResponse `json:"state" form:"state"`
}

// Validate -> workflowRequest
func (a workflowRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.WorkflowType, validation.Required),
		validation.Field(&a.Name, validation.Required),
		validation.Field(&a.Status, validation.Required, validation.Min(1), validation.Max(2)),
		// validation.Field(&a.State, validation.Required),
	)
}

func WorkflowCreate(c *fiber.Ctx) error {
	user := common.GetUser(c)

	inputData := new(workflowRequest)
	_ = c.BodyParser(inputData)
	err := inputData.Validate()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.JSONResponse{
			Code:    fiber.StatusBadRequest,
			Success: false,
			Message: "WorkflowCreate BadRequest",
			Data:    err,
		})
	}

	// attach at least one file
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.JSONResponse{
			Code:    fiber.StatusBadRequest,
			Success: false,
			Message: fmt.Sprintf("MultipartForm documents: %v", err),
			Data:    nil,
		})
	}
	files := form.File["documents"]
	if len(files) < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(model.JSONResponse{
			Code:    fiber.StatusBadRequest,
			Success: false,
			Message: "Attach at least one file to create a workflow.",
			Data:    nil,
		})
	}

	// Unmarshal wfState
	wfState := new([]workflowStateResponse)
	if inputData.State != "" {
		err = json.Unmarshal([]byte(inputData.State), wfState)
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

	// check duplicate name
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

	// check state order
	if len(*wfState) > 0 {
		err = checkStateOrder(*wfState)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.JSONResponse{
				Code:    fiber.StatusBadRequest,
				Success: false,
				Message: fmt.Sprintf("checkStateOrder: %v", err),
				Data:    nil,
			})
		}
	}

	// store workflow & workflow_state
	workflowID, err := storeWorkflow(user, inputData, *wfState)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.JSONResponse{
			Code:    fiber.StatusInternalServerError,
			Success: false,
			Message: "WorkflowCreate: storeWorkflow: error",
			Data:    err,
		})
	}

	// file upload
	path := fmt.Sprintf("%susers/%s/%d/", config.C.DevPrefix, user, workflowID)
	var uploadRes model.JSONResponse
	for _, file := range files {
		_, _, uploadRes = uploadFile(file, workflowID, user, path)
	}
	if uploadRes.Code != fiber.StatusOK {
		return c.Status(uploadRes.Code).JSON(uploadRes)
	}

	res := workflowResponse{
		WorkflowID:   workflowID,
		WorkflowType: inputData.WorkflowType,
		Name:         inputData.Name,
		Description:  inputData.Description,
		State:        *wfState,
	}
	// return workflow data
	return c.Status(fiber.StatusOK).JSON(model.JSONResponse{
		Code:    fiber.StatusOK,
		Success: true,
		Message: fmt.Sprintf("create workflow `%s` by `%s` successed.", inputData.Name, user),
		Data:    res,
	})
}

func checkStateOrder(state []workflowStateResponse) error {
	var oldOrder int
	for _, v := range state {
		if v.Order == 0 {
			return fmt.Errorf("state order invalid, state order number start from 1")
		}
		if v.Order < oldOrder {
			return fmt.Errorf("state order invalid, state order number should not be duplicated")
		}
		oldOrder = v.Order
	}
	return nil
}

func storeWorkflow(user string, inputData *workflowRequest, state []workflowStateResponse) (workflowID int, err error) {
	now := time.Now()

	// store workflow
	workflowData := model.Workflow{
		WorkflowType: inputData.WorkflowType,
		Name:         inputData.Name,
		Description:  inputData.Description,
		Status:       inputData.Status,
		CreatedBy:    user,
		CreatedDate:  now,
		UpdatedDate:  now,
	}
	err = db.DB.Table(model.TableNameWorkflow).Create(&workflowData).Error
	if err != nil {
		db.DB.Rollback()
		return 0, err
	}

	// store workflow_state
	if len(state) > 0 {
		insertStateStr := []string{}
		insertStateArgs := []interface{}{}
		for _, v := range state {
			insertStateStr = append(insertStateStr, "(?, ?, ?, ?, ?)")
			insertStateArgs = append(insertStateArgs, workflowData.ID)
			insertStateArgs = append(insertStateArgs, v.Name)
			insertStateArgs = append(insertStateArgs, v.Order)
			insertStateArgs = append(insertStateArgs, 1)
			insertStateArgs = append(insertStateArgs, v.AssignedTo)
		}
		smt := fmt.Sprintf("INSERT INTO `%s` (`workflow_id`,`name`,`state_order`,`status`,`assigned_to`) VALUES %s", model.TableNameWorkflowState, "%s")
		smt = fmt.Sprintf(smt, strings.Join(insertStateStr, ","))
		if err := db.DB.Exec(smt, insertStateArgs...).Error; err != nil {
			db.DB.Rollback()
			return workflowData.ID, err
		}
	}

	return workflowData.ID, nil
}
