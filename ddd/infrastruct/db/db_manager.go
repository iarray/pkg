package db

import (
	"time"

	"gorm.io/gorm"
)

var db *gorm.DB

func InitDb(dbType string, conString string) error {
	a, err := gorm.Open(dbType, conString)
	if err != nil {
		return err
	}
	db = a
	db.Callback().Create().Before("gorm:create").Register("update_created_at", updateCreated)
	db.Callback().Update().Before("gorm:update").Register("update_updated_at", updateUpdated)
	return nil
}

func I() *gorm.DB {
	return db
}

func updateCreated(scope *gorm.Scope) {
	if scope.HasColumn("CreatedAt") {
		scope.SetColumn("CreatedAt", time.Now())
	}
}

func updateUpdated(scope *gorm.Scope) {
	if scope.HasColumn("UpdatedAt") {
		scope.SetColumn("UpdatedAt", time.Now())
	}
}
