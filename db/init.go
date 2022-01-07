package db

import (
	"fmt"
	"signdoc_api/config"
	"signdoc_api/model"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// const dbType string = "mysql"
// const dbType string = "sqlite"

// DB is Database Instance
var DB *gorm.DB

// Init connect database
func Init() {
	var err error
	// database connect
	dbType := config.C.DBtype
	switch dbType {
	case "mysql":
		dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=true&loc=Local",
			config.C.DB.Username,
			config.C.DB.Password,
			config.C.DB.Host,
			config.C.DB.Port,
			config.C.DB.Database)
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "sqlite":
		DB, err = gorm.Open(sqlite.Open("./db/dump/dump.db"), &gorm.Config{})
	}
	if err != nil {
		panic(fmt.Errorf("Error connect db %v %v %v", dbType, config.C.DB.Database, err))
	}

	/* Migration */
	DB.AutoMigrate(&model.FileData{})
	DB.AutoMigrate(&model.FileMeta{})
	DB.AutoMigrate(&model.Workflow{})
	DB.AutoMigrate(&model.WorkflowHistory{})
	DB.AutoMigrate(&model.WorkflowHistoryType{})
	DB.AutoMigrate(&model.WorkflowState{})
	DB.AutoMigrate(&model.WorkflowStateStatusType{})
	DB.AutoMigrate(&model.WorkflowStatus{})
	DB.AutoMigrate(&model.WorkflowType{})

	/* insert init data */
	DB.Table(model.TableNameWorkflowHistoryType).Create(model.InitWorkflowHistoryType)
	DB.Table(model.TableNameWorkflowStateStatusType).Create(model.InitWorkflowStateStatusType)
	DB.Table(model.TableNameWorkflowStatus).Create(model.InitWorkflowStatus)
	DB.Table(model.TableNameWorkflowType).Create(model.InitWorkflowType)

}
