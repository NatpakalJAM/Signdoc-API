package model

// WorkflowState for TableName
const TableNameWorkflowState = "workflow_state"

// WorkflowState is model for workflow_state
type WorkflowState struct {
	ID         int    `gorm:"column:id;primary_key:true;AUTO_INCREMENT" json:"id"`
	WorkflowID int    `gorm:"column:workflow_id" json:"workflow_id"` //fk
	Name       string `gorm:"column:name;size:255" json:"name"`
	Order      int    `gorm:"column:state_order" json:"state_order"`
	Status     int    `gorm:"column:status" json:"-"`
	StatusStr  string `gorm:"-" json:"status"`
	AssignedTo string `gorm:"column:assigned_to;size:100" json:"assigned_to"` // fk
	// FK
	WorkflowHistory         []WorkflowHistory       `gorm:"foreignKey:WorkflowStateID;references:ID" json:"workflow_history"`
	WorkflowStateStatusType WorkflowStateStatusType `gorm:"foreignKey:Status;references:ID" json:"-"`
}

// TableName for model WorkflowState
func (WorkflowState) TableName() string {
	return "workflow_state"
}
