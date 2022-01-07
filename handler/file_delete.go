package handler

import (
	"fmt"
	"signdoc_api/common"
	"signdoc_api/config"
	"signdoc_api/db"
	"signdoc_api/gcp"
	"signdoc_api/model"
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofiber/fiber/v2"
)

type fileDeleteRequest struct {
	FileName string `json:"files" form:"files"`
}

// Validate -> fileDeleteRequest
func (a fileDeleteRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.FileName, validation.Required),
	)
}

func FileDelete(c *fiber.Ctx) error {
	user := common.GetUser(c)

	workflowID, _ := strconv.Atoi(c.Params("workflow_id"))

	inputData := new(fileDeleteRequest)
	_ = c.BodyParser(inputData)
	err := inputData.Validate()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.JSONResponse{
			Code:    fiber.StatusBadRequest,
			Success: false,
			Message: "FileDelete BadRequest",
			Data:    err,
		})
	}

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

	files := strings.Split(inputData.FileName, ",")

	path := fmt.Sprintf("%susers/%s/%d/", config.C.DevPrefix, user, workflowID)

	for _, fileName := range files {

		// find fileMeta by fileName
		fileMeta := findFileMetaByFileName(fileName, path)
		if fileMeta == "" {
			return c.Status(fiber.StatusBadRequest).JSON(model.JSONResponse{
				Code:    fiber.StatusBadRequest,
				Success: false,
				Message: fmt.Sprintf("file name `%s` not found in workflow_id %d", fileName, workflowID),
				Data:    nil,
			})
		}

		// delete file on GCP
		err = gcp.Client.DeleteFiles(config.C.GCP.BucketName, path, fileMeta)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.JSONResponse{
				Code:    fiber.StatusInternalServerError,
				Success: false,
				Message: fmt.Sprintf("FileDelete: DeleteFiles: %v", err),
				Data:    nil,
			})
		}

		// delete fils on DB
		err = deleteFileMeta(workflowID, fileName, fileMeta, path)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.JSONResponse{
				Code:    fiber.StatusInternalServerError,
				Success: false,
				Message: fmt.Sprintf("FileDelete: deleteFileMeta:%v", err),
				Data:    nil,
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(model.JSONResponse{
		Code:    fiber.StatusOK,
		Success: true,
		Message: "files deleted.",
		Data:    nil,
	})
}

func findFileMetaByFileName(fileName, filePath string) (fileMeta string) {
	var listFileMeta model.FileMeta
	db.DB.Table(model.TableNameFileMeta).
		Select("save_file_name").
		Where("original_file_name = ? AND file_path = ? AND deleted = 0", fileName, filePath).
		Find(&listFileMeta)
	return listFileMeta.SaveFileName
}

func deleteFileMeta(workflowID int, fileName, fileMeta, path string) error {
	// delete file.Filename/meta
	now := time.Now()
	err := db.DB.Model(&model.FileMeta{}).
		Where("workflow_id = ? AND original_file_name = ? AND save_file_name = ? AND file_path = ? AND deleted = 0", workflowID, fileName, fileMeta, path).
		Updates(&model.FileMeta{
			Deleted:     true,
			UpdatedDate: now,
		}).
		Error
	if err != nil {
		db.DB.Rollback()
	}
	return err
}
