package handler

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"signdoc_api/common"
	"signdoc_api/config"
	"signdoc_api/db"
	"signdoc_api/gcp"
	"signdoc_api/model"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Sizer interface {
	Size() int64
}

const limitFileByteSize int64 = 1048576 // 10MB

func FileUpload(c *fiber.Ctx) error {
	user := common.GetUser(c)

	workflowID, _ := strconv.Atoi(c.Params("workflow_id"))

	// check workflowID belong to user & workflow on draft
	errCode, err := checkWorkflowIDPermission(workflowID, user)
	if err != nil {
		return c.Status(errCode).JSON(model.JSONResponse{
			Code:    errCode,
			Success: false,
			Message: fmt.Sprintf("checkWorkflowIDPermission: %v", err),
			Data:    nil,
		})
	}

	path := fmt.Sprintf("%susers/%s/%d/", config.C.DevPrefix, user, workflowID)

	// file upload
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.JSONResponse{
			Code:    fiber.StatusBadRequest,
			Success: false,
			Message: fmt.Sprintf("MultipartForm documents: %v", err),
			Data:    nil,
		})
	}
	type listFileName struct {
		FileName string `json:"file_name"`
		FileMeta string `json:"file_meta"`
	}
	var res model.JSONResponse
	files := form.File["documents"]
	if len(files) < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(model.JSONResponse{
			Code:    fiber.StatusBadRequest,
			Success: false,
			Message: "Attach at least one file to upload.",
			Data:    nil,
		})
	}
	lists := make([]listFileName, len(files))
	for i, file := range files {
		lists[i].FileName, lists[i].FileMeta, res = uploadFile(file, workflowID, user, path)
	}
	if res.Code == fiber.StatusOK {
		res.Data = lists
	}

	return c.Status(res.Code).JSON(res)
}

func uploadFile(file *multipart.FileHeader, workflowID int, user, path string) (fileName, fileMeta string, res model.JSONResponse) {
	fileName = file.Filename

	// Get first file from form field "documents":
	blobFile, err := file.Open()
	if err != nil {
		return fileName, fileMeta, model.JSONResponse{
			Code:    fiber.StatusInternalServerError,
			Success: false,
			Message: fmt.Sprintf("file.Open documents: %v", err),
			Data:    nil,
		}
	}

	// check duplicate fileName
	fileExist := checkDuplicateFileName(fileName, path)
	if fileExist {
		return fileName, fileMeta, model.JSONResponse{
			Code:    fiber.StatusBadRequest,
			Success: false,
			Message: "duplicate fileName in this location",
			Data:    nil,
		}
	}

	// check file mime type PDF
	buff := make([]byte, 512)
	_, err = blobFile.Read(buff)
	if err != nil {
		return fileName, fileMeta, model.JSONResponse{
			Code:    fiber.StatusInternalServerError,
			Success: false,
			Message: fmt.Sprintf("blobFile.Read: %v", err),
			Data:    nil,
		}
	}
	filetype := http.DetectContentType(buff)
	if filetype != "application/pdf" {
		return fileName, fileMeta, model.JSONResponse{
			Code:    fiber.StatusBadRequest,
			Success: false,
			Message: fmt.Sprintf("Query `documents` request File type or MIME type accept `application/pdf`, got `%s`", filetype),
			Data:    nil,
		}
	}
	if blobFile.(Sizer).Size() > limitFileByteSize {
		return fileName, fileMeta, model.JSONResponse{
			Code:    fiber.StatusBadRequest,
			Success: false,
			Message: "Query `documents` limit file size 10MB",
			Data:    nil,
		}
	}

	now := time.Now()
	fileMeta = common.GenerateFileMeta(now, fileName)
	fileMeta = checkDuplicateFileMeta(now, fileName, fileMeta, path)

	err = gcp.Client.UploadFile(file, path, fileMeta)
	if err != nil {
		return fileName, fileMeta, model.JSONResponse{
			Code:    fiber.StatusInternalServerError,
			Success: false,
			Message: fmt.Sprintf("GCP UploadFile: %v", err),
			Data:    nil,
		}
	}

	err = storeFileMeta(now, workflowID, fileName, fileMeta, user, path)
	if err != nil {
		return fileName, fileMeta, model.JSONResponse{
			Code:    fiber.StatusInternalServerError,
			Success: false,
			Message: fmt.Sprintf("storeFileMeta: %v", err),
			Data:    nil,
		}
	}

	return fileName, fileMeta, model.JSONResponse{
		Code:    fiber.StatusOK,
		Success: true,
		Message: "file uploaded",
		Data:    nil,
	}
}

func checkWorkflowIDPermission(workflowID int, user string) (int, error) {
	if workflowID <= 0 {
		return fiber.StatusBadRequest, fmt.Errorf("workflow_id should be number and over 0")
	}
	var wf model.Workflow
	ck := db.DB.Table(model.TableNameWorkflow).Where("id = ? AND created_by = ?", workflowID, user).Find(&wf).RowsAffected
	if ck <= 0 {
		return fiber.StatusForbidden, fmt.Errorf("user `%s` does not have access to workflow_id %d", user, workflowID)
	}
	if wf.Status != 1 {
		return fiber.StatusBadRequest, fmt.Errorf("add file denied; workflow_id %d is not in draft state", workflowID)
	}
	return 0, nil
}

func storeFileMeta(now time.Time, workflowID int, fileName, fileMeta, user, path string) (err error) {
	// store file.Filename/meta
	data := model.FileMeta{
		WorkflowID:       workflowID,
		OriginalFileName: fileName,
		SaveFileName:     fileMeta,
		File_uploader:    user,
		FilePath:         path,
		Deleted:          false,
		CreatedDate:      now,
		UpdatedDate:      now,
	}
	err = db.DB.Table(model.TableNameFileMeta).Create(&data).Error
	if err != nil {
		db.DB.Rollback()
	}
	return err
}

func checkDuplicateFileName(fileName, filePath string) (fileExist bool) {
	// var listFileMeta []model.FileMeta
	row := db.DB.Table(model.TableNameFileMeta).
		Select("id").
		Where("original_file_name = ? AND file_path = ?", fileName, filePath).
		Find(&model.FileMeta{}).
		RowsAffected
	if row > 0 {
		fileExist = true
	}
	return fileExist
}

func checkDuplicateFileMeta(now time.Time, fileName, fileMeta, filePath string) string {
	// var listFileMeta []model.FileMeta
	row := db.DB.Table(model.TableNameFileMeta).
		Select("id").
		Where("save_file_name = ? AND file_path = ?", fileMeta, filePath).
		Find(&model.FileMeta{}).
		RowsAffected
	if row > 0 {
		now = now.Add(time.Second * 1)
		fileMeta = common.GenerateFileMeta(now, fileName)
		checkDuplicateFileMeta(now, fileName, fileMeta, filePath)
	}
	return fileMeta
}
