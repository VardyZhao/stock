package migration

import (
	"gorm.io/gorm"
	"stock/model"
)

var migrateMap = map[string][]interface{}{
	"default": {
		&model.Department{},
		&model.Stock{},
	},
}

func Run(name string, db *gorm.DB) {
	if modelList, exists := migrateMap[name]; exists {
		_ = db.AutoMigrate(modelList...)
	}
}
