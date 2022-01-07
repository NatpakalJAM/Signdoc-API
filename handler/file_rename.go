package handler

import (
	"fmt"
	"signdoc_api/common"
	"signdoc_api/config"
	"signdoc_api/db"
	"signdoc_api/model"
	"strconv"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofiber/fiber/v2"
)

type fileRenameRequest struct {
	OldFileName string `json:"old_file_name" form:"old_file_name"`
	NewFileName string `json:"new_file_name" form:"new_file_name"`
}

// Validate -> fileRenameRequest
func (a fileRenameRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.OldFileName, validation.Required),
		validation.Field(&a.NewFileName, validation.Required),
	)
}

func FileRename(c *fiber.Ctx) error {
	user := common.GetUser(c)

	workflowID, _ := strconv.Atoi(c.Params("workflow_id"))

	inputData := new(fileRenameRequest)
	_ = c.BodyParser(inputData)
	err := inputData.Validate()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.JSONResponse{
			Code:    fiber.StatusBadRequest,
			Success: false,
			Message: "FileRename BadRequest",
			Data:    err,
		})
	}

	path := fmt.Sprintf("%susers/%s/%d/", config.C.DevPrefix, user, workflowID)

	// check workflow status = 1(draft)
	errCode, err := checkWorkflowIDPermission(workflowID, user)
	if err != nil {
		return c.Status(errCode).JSON(model.JSONResponse{
			Code:    errCode,
			Success: false,
			Message: fmt.Sprintf("checkWorkflowIDPermission: %v", err),
			Data:    nil,
		})
	}

	// check original_file_name exist
	err = checkFileNameExist(workflowID, inputData.OldFileName, path)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.JSONResponse{
			Code:    fiber.StatusBadRequest,
			Success: false,
			Message: fmt.Sprintf("FileRename: checkFileNameExist: %v", err),
			Data:    err,
		})
	}

	// rename original_file_name in DB file_meta
	err = renameFile(workflowID, inputData.OldFileName, inputData.NewFileName, path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.JSONResponse{
			Code:    fiber.StatusInternalServerError,
			Success: false,
			Message: fmt.Sprintf("FileRename: renameFile: %v", err),
			Data:    err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.JSONResponse{
		Code:    fiber.StatusOK,
		Success: true,
		Message: fmt.Sprintf("file rename succeed. workflow_id %d : %s -> %s", workflowID, inputData.OldFileName, inputData.NewFileName),
		Data:    nil,
	})
}

// check original_file_name exist
func checkFileNameExist(workflowID int, fileName, path string) error {
	row := db.DB.Table(model.TableNameFileMeta).
		Where("workflow_id = ? AND original_file_name = ? AND file_path = ? AND deleted = 0", workflowID, fileName, path).
		Find(&model.FileMeta{}).
		RowsAffected
	if row < 1 {
		return fmt.Errorf("file `%s` not found in workflow_id %d", fileName, workflowID)
	}
	return nil
}

// rename original_file_name in DB file_meta
func renameFile(workflowID int, oldName, newName, path string) error {
	now := time.Now()
	err := db.DB.Model(&model.FileMeta{}).
		Where("workflow_id = ? AND original_file_name = ? AND file_path = ? AND deleted = 0", workflowID, oldName, path).
		Updates(&model.FileMeta{
			OriginalFileName: newName,
			UpdatedDate:      now,
		}).
		Error
	if err != nil {
		db.DB.Rollback()
	}
	return err
}
