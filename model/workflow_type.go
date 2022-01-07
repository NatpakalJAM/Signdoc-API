package model

// WorkflowType for TableName
const TableNameWorkflowType = "workflow_type"

// WorkflowType is model for workflow_type
type WorkflowType struct {
	ID   int    `gorm:"column:id;primary_key:true;AUTO_INCREMENT" json:"id"`
	Name string `gorm:"column:name;size:255" json:"name"`
}

// TableName for model WorkflowType
func (WorkflowType) TableName() string {
	return "workflow_type"
}

/*-----------------------------------------------------------------*/
// init db

var InitWorkflowType = []WorkflowType{
	{
		ID:   1,
		Name: "Sign Document",
	},
}
