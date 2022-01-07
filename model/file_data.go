package model

import (
	"time"
)

// FileData for TableName
const TableNameFileData = "file_data"

// FileData is model for file_data
type FileData struct {
	ID          int       `gorm:"column:id;primary_key:true;AUTO_INCREMENT" json:"id"`
	FileMetaID  int       `gorm:"column:filemeta_id" json:"filemeta_id"`
	FromUser    string    `gorm:"column:from_user;size:100" json:"from_user"`
	Data        string    `gorm:"column:data;type:text" json:"data"`
	CreatedDate time.Time `gorm:"column:created_date;type:datetime" json:"created_date"`
}

// TableName for model FileData
func (FileData) TableName() string {
	return "file_data"
}
