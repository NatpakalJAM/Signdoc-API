package model

// WorkflowStatus for TableName
const TableNameWorkflowStatus = "workflow_status"

// WorkflowStatus is model for workflow_status
type WorkflowStatus struct {
	ID   int    `gorm:"column:id;primary_key:true;AUTO_INCREMENT" json:"id"`
	Name string `gorm:"column:name;size:255" json:"name"`
}

// TableName for model WorkflowStatus
func (WorkflowStatus) TableName() string {
	return "workflow_status"
}

/*-----------------------------------------------------------------*/
// init db

var InitWorkflowStatus = []WorkflowStatus{
	{
		ID:   1,
		Name: "draft",
	},
	{
		ID:   2,
		Name: "publish",
	},
	{
		ID:   3,
		Name: "processing",
	},
	{
		ID:   4,
		Name: "complete",
	},
	{
		ID:   5,
		Name: "cancel",
	},
}
