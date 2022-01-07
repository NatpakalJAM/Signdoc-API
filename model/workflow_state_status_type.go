package model

// WorkflowStateStatusType for TableName
const TableNameWorkflowStateStatusType = "workflow_state_status_type"

// WorkflowStateStatusType is model for workflow_state_status_type
type WorkflowStateStatusType struct {
	ID   int    `gorm:"column:id;primary_key:true;AUTO_INCREMENT" json:"id"`
	Name string `gorm:"column:name;size:255" json:"name"`
}

// TableName for model WorkflowStateStatusType
func (WorkflowStateStatusType) TableName() string {
	return "workflow_state_status_type"
}

/*-----------------------------------------------------------------*/
// init db

var InitWorkflowStateStatusType = []WorkflowType{
	{
		ID:   1,
		Name: "waiting",
	},
	{
		ID:   2,
		Name: "succeeded",
	},
	{
		ID:   3,
		Name: "rejected",
	},
}
