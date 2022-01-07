package model

import (
	"time"
)

// FileMeta for TableName
const TableNameFileMeta = "file_meta"

// FileMeta is model for file_meta
type FileMeta struct {
	ID               int       `gorm:"column:id;primary_key:true;AUTO_INCREMENT" json:"id"`
	WorkflowID       int       `gorm:"column:workflow_id" json:"workflow_id"`
	OriginalFileName string    `gorm:"column:original_file_name;size:100" json:"original_file_name"`
	SaveFileName     string    `gorm:"column:save_file_name;size:100;index" json:"save_file_name"`
	File_uploader    string    `gorm:"column:file_uploader;size:50;index" json:"file_uploader"`
	FilePath         string    `gorm:"column:file_path;type:text" json:"file_path"`
	Deleted          bool      `gorm:"column:deleted" json:"deleted"`
	CreatedDate      time.Time `gorm:"column:created_date;type:datetime" json:"created_date"`
	UpdatedDate      time.Time `gorm:"column:updated_date;type:datetime" json:"updated_date"`
}

// TableName for model FileMeta
func (FileMeta) TableName() string {
	return "file_meta"
}
