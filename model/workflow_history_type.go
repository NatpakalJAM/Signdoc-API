package model

// WorkflowHistoryType for TableName
const TableNameWorkflowHistoryType = "workflow_history_type"

// WorkflowHistoryType is model for workflow_history_type
type WorkflowHistoryType struct {
	ID   int    `gorm:"column:id;primary_key:true;AUTO_INCREMENT" json:"id"`
	Name string `gorm:"column:name;size:255" json:"name"`
}

// TableName for model WorkflowHistoryType
func (WorkflowHistoryType) TableName() string {
	return "workflow_history_type"
}

/*-----------------------------------------------------------------*/
// init db

var InitWorkflowHistoryType = []WorkflowHistoryType{
	{
		ID:   1,
		Name: "verified",
	},
	{
		ID:   2,
		Name: "declined",
	},
}
