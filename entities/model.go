package entities

import (
	"log"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type Todo struct {
	Id          int32  `gorm:"primaryKey;autoIncrement"`
	Title       string `gorm:"size:255"`
	Description string `gorm:"size:1000"`
	IsCompleted bool
}

func InitDB() *gorm.DB {
	dsn := "sqlserver://sa:m@s@192.168.5.130:1433?database=TodoList"

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	err = db.AutoMigrate(&Todo{})
	if err != nil {
		log.Fatal("failed to migrate schema:", err)
	}
	return db
}
