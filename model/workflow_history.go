package model

import "time"

// WorkflowHistory for TableName
const TableNameWorkflowHistory = "workflow_history"

// WorkflowHistory is model for workflow_history
type WorkflowHistory struct {
	ID              int       `gorm:"column:id;primary_key:true;AUTO_INCREMENT" json:"id"`
	WorkflowStateID int       `gorm:"column:workflow_state_id" json:"workflow_state_id"` //fk
	HistoryType     int       `gorm:"column:type;size:255" json:"-"`                     //fk
	HistoryTypeStr  string    `gorm:"-" json:"type"`
	Date            time.Time `gorm:"column:date;type:datetime" json:"date"`
	// FK
	WorkflowHistoryType WorkflowHistoryType `gorm:"foreignKey:HistoryType;references:ID" json:"-"`
}

// TableName for model WorkflowHistory
func (WorkflowHistory) TableName() string {
	return "workflow_history"
}
