package model

import "time"

// Workflow for TableName
const TableNameWorkflow = "workflow"

// Workflow is model for workflow
type Workflow struct {
	ID              int       `gorm:"column:id;primary_key:true;AUTO_INCREMENT" json:"id"`
	WorkflowType    int       `gorm:"column:workflow_type" json:"-"`
	WorkflowTypeStr string    `gorm:"-" json:"type"`
	Name            string    `gorm:"column:name;size:255" json:"name"`
	Description     string    `gorm:"column:description;type:text" json:"description"`
	Status          int       `gorm:"column:status" json:"-"`
	StatusStr       string    `gorm:"-" json:"status"`
	StatusMessage   string    `gorm:"column:status_message" json:"status_message"`
	CreatedBy       string    `gorm:"column:created_by" json:"created_by"`
	CreatedDate     time.Time `gorm:"column:created_date;type:datetime" json:"created_date"`
	UpdatedDate     time.Time `gorm:"column:updated_date;type:datetime" json:"updated_date"`
	// FK
	WorkflowState      []WorkflowState `gorm:"foreignKey:WorkflowID;references:ID" json:"workflow_state,omitempty"`
	WorkflowTypeSelect WorkflowType    `gorm:"foreignKey:WorkflowType;references:ID" json:"-"`
	WorkflowStatus     WorkflowStatus  `gorm:"foreignKey:Status;references:ID" json:"-"`
}

// TableName for model Workflow
func (Workflow) TableName() string {
	return "workflow"
}
